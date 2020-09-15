package ad

import (
	"context"
	"log"
	"time"

	"../location"
	"../tguser"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	//StatusAdding - Добавляется
	StatusAdding = 0
	//StatusWait - Ожидает публикации
	StatusWait = 1
	//StatusPublished - Опубликовано
	StatusPublished = 2
	//StatusExpire - истёк срок
	StatusExpire = 3
	//StatusBlock - заблокировано
	StatusBlock = 4
	//StatusComplete - завершено/продано
	StatusComplete = 5
)

//Ad - Объявление
type Ad struct {
	ID          primitive.ObjectID `bson:"_id"`
	Title       string             //заголовок
	Description string
	Category    []string
	Address     location.Address
	Price       int64
	UserID      int64             //user ID в телеграме
	Name        string            //Имя в контактах
	Phone       string            //телефон в контактах
	Location    location.Location //
	Photos      []string
	Status      int16
	AddTime     time.Time //время подачи
	PublishTime time.Time //время публикации
}

//PhotoTmp - временная колекция для фото .возникают конкурентные вычисления при загрузке нескольких фото
type PhotoTmp struct {
	ID     primitive.ObjectID `bson:"_id"`
	FileID string
	AdID   primitive.ObjectID
}

//CntItem - Элемент AdsCount
type CntItem struct {
	Count int64
	Time  time.Time
}

//AdsCount Количество объявлений всего
var AdsCount = make(map[string]CntItem)

//Count - считает количество объявлений в рубрике(кеширует в модуле)
func Count(rubric string) int64 {
	calc := false //признак что наджо пересчитать
	if r, ok := AdsCount[rubric]; ok {
		dl := r.Time.Add(time.Minute * 10)
		if dl.Before(time.Now()) {
			//время жизни значения 10 минут. потом пересчитать
			calc = true
		}
	} else { //если такой еще нет
		AdsCount[rubric] = CntItem{
			Count: 0,
			Time:  time.Now(),
		}
		calc = true
	}

	if calc {
		var (
			cnt    int64
			err    error
			filter interface{}
		)
		if rubric == "total" {
			filter = bson.M{}

		} else {
			filter = bson.M{"rubric": bson.M{"$in": rubric}}
		}
		if cnt, err = collection.CountDocuments(context.TODO(), filter); err != nil {
			log.Printf("ERROR IN COUNT: %s", err)
		}
		var tmp CntItem
		tmp.Count = cnt
		tmp.Time = time.Now()
		AdsCount[rubric] = tmp
	}
	return AdsCount[rubric].Count
}

//ClearPhotoTmp - удалаяет все объекты PhotoTmp для объявления AdID
func ClearPhotoTmp(AdID primitive.ObjectID) {
	c := dbClient.Database("qprobot").Collection("phototmp")

	filter := bson.M{"AdID": AdID}
	_, err := c.DeleteMany(context.TODO(), filter)
	if err != nil {
		log.Printf("CRITICAL ERROR: %s", err)
	}
}

//PhotoTmpCnt -
func (a *Ad) PhotoTmpCnt() int64 {
	c := dbClient.Database("qprobot").Collection("phototmp")
	filter := bson.M{"AdID": a.ID}
	cnt, err := c.CountDocuments(context.TODO(), filter)
	if err != nil {
		log.Printf("ERROR IN COUNT: %s", err)
	}
	return cnt
}

//PhotoTmpToAd - Переносит временные фотки в объяву
func (a *Ad) PhotoTmpToAd() {
	c := dbClient.Database("qprobot").Collection("phototmp")
	filter := bson.M{"AdID": a.ID}
	cur, err := c.Find(context.TODO(), filter)
	if err != nil {
		log.Printf("CRITICAL ERROR: %s", err)
	}
	for cur.Next(context.TODO()) {
		var pt PhotoTmp
		err := cur.Decode(&pt)
		if err != nil {
			log.Printf("CRITICAL ERROR: %s", err)
		}
		a.Photos = append(a.Photos, pt.FileID)
	}
	a.Save()
}

//GetPhotoTmp -получить временные фотки
func (a *Ad) GetPhotoTmp() []string {
	c := dbClient.Database("qprobot").Collection("phototmp")
	filter := bson.M{"AdID": a.ID}
	cur, err := c.Find(context.TODO(), filter)
	if err != nil {
		log.Printf("CRITICAL ERROR: %s", err)
	}
	res := []string{}
	for cur.Next(context.TODO()) {
		var pt PhotoTmp
		err := cur.Decode(&pt)
		if err != nil {
			log.Printf("CRITICAL ERROR: %s", err)
		}
		res = append(res, pt.FileID)
	}
	return res
}

//Create -
func (pt *PhotoTmp) Create() error {
	c := dbClient.Database("qprobot").Collection("phototmp")
	_, err := c.InsertOne(context.TODO(), bson.M{
		"FileID": pt.FileID,
		"AdID":   pt.AdID,
	})
	if err != nil {
		return err
	}
	return nil
}

var dbClient *mongo.Client
var collection *mongo.Collection

//Init - инициализация модуля
func Init(db *mongo.Client) {
	if dbClient == nil {
		dbClient = db
		collection = dbClient.Database("qprobot").Collection("ads")
	}
}

//Create -
func Create(user *tguser.TgUser, category string) (Ad, error) {
	res, err := collection.InsertOne(context.TODO(), bson.M{
		"UserID":   user.ID,
		"Status":   StatusAdding,
		"Category": category,
		"AddTime":  time.Now(),
	})
	if err != nil {
		return Ad{}, err
	}
	user.AddingAd = res.InsertedID.(primitive.ObjectID)
	user.Save()
	return Ad{ID: res.InsertedID.(primitive.ObjectID)}, nil
}

//Save - сохранить данные
func (a *Ad) Save() {
	opts := options.Update().SetUpsert(true)
	filter := bson.M{"_id": a.ID}
	update := bson.M{"$set": bson.M{
		"Title":       a.Title,
		"Description": a.Description,
		"Category":    a.Category,
		"Price":       a.Price,
		"UserID":      a.UserID,
		"Name":        a.Name,
		"Phone":       a.Phone,
		"Photos":      a.Photos,
		"Address":     a.Address,
		"Location":    a.Location,
		"Status":      a.Status,
		"AddTime":     a.AddTime,
		"PublishTime": a.PublishTime,
	}}

	_, err := collection.UpdateOne(context.TODO(), filter, update, opts)
	if err != nil {
		log.Printf("CRITICAL ERROR: %s", err)
	}
}

//Get - найти
func Get(ID primitive.ObjectID) *Ad {
	var ad Ad
	filter := bson.M{"_id": ID}
	err := collection.FindOne(context.TODO(), filter).Decode(&ad)
	if err != nil {
		return nil
	}
	return &ad
}

//Delete -
func (a *Ad) Delete() {
	ClearPhotoTmp(a.ID)
	delete(a.ID)
}

//Delete -
func delete(ID primitive.ObjectID) {
	filter := bson.M{"_id": ID}
	_, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Printf("CRITICAL ERROR: %s", err)
	}
}

//List -
func List(offset int64, limit int64) []Ad {
	filter := bson.M{"Status": StatusPublished}
	options := options.Find()
	options.SetSkip(offset)
	options.SetLimit(limit)
	//options.SetSort(bson.D{{"_id", -1}})
	res := []Ad{}
	cur, err := collection.Find(context.TODO(), filter, options)
	if err != nil {
		log.Printf("CRITICAL ERROR: %s", err)
	}
	for cur.Next(context.TODO()) {
		var a Ad
		err := cur.Decode(&a)
		if err != nil {
			log.Printf("CRITICAL ERROR: %s", err)
		}
		res = append(res, a)
	}
	return res
}
