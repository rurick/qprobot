package tguser

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//SetLastMsgID -
func (u *TgUser) SetLastMsgID(msgID int) {
	col := dbClient.Database("qprobot").Collection("lastMsgIDs")
	opts := options.Update().SetUpsert(true)
	filter := bson.M{"UID": u.ID}
	update := bson.M{
		"$set": bson.M{
			"MsgID": msgID,
		},
	}
	_, err := col.UpdateOne(context.TODO(), filter, update, opts)
	if err != nil {
		_, err := col.InsertOne(context.TODO(), bson.M{"UID": u.ID, "MsgID": msgID})
		if err != nil {
		}
	}
}

//GetLastMsgID -
func (u *TgUser) GetLastMsgID() int {
	type mt struct {
		UID   int
		MsgID int
	}
	col := dbClient.Database("qprobot").Collection("lastMsgIDs")
	filter := bson.M{"UID": u.ID}
	d := mt{}
	err := col.FindOne(context.TODO(), filter).Decode(&d)
	if err != nil {
		return -1
	}
	return d.MsgID
}

//TgUser - данные пользователя
//db.users.createIndex({"Username" : 1}, {"unique" : true})
type TgUser struct {
	ID        int64
	Username  string //телеграм имя
	Name      string //имя
	Phone     string //телефон
	State     string //текущее состояние. Где и в каком месте находится пользователь
	StateBack string //пердыдущее состояние. Куда возвращаться
	Location  struct {
		Longitude float64
		Latitude  float64
	} //текущее состояние. Где и в каком месте находится пользователь
	LastVisit time.Time            //Время последнего визита
	Ads       []primitive.ObjectID //объявления пользователя
	AddingAd  primitive.ObjectID   //Добавляемое сейчас ad
	AdData    interface{}          //данные для объявлений
}

var dbClient *mongo.Client
var collection *mongo.Collection

//Save - сохранить
func (u *TgUser) Save() {
	opts := options.Update().SetUpsert(true)
	filter := bson.M{"ID": u.ID}
	u.LastVisit = time.Now()
	update := bson.M{
		"$set": bson.M{
			"Username":  u.Username,
			"Name":      u.Name,
			"Phone":     u.Phone,
			"State":     u.State,
			"StateBack": u.StateBack,
			"Location":  u.Location,
			"LastVisit": u.LastVisit,
			"Ads":       u.Ads,
			"AddingAd":  u.AddingAd,
			"AdData":    u.AdData,
		},
	}

	_, err := collection.UpdateOne(context.TODO(), filter, update, opts)
	if err != nil {
		log.Printf("CRITICAL ERROR: %s", err)
	}
}

//Init - инициализация модуля
func Init(db *mongo.Client) {
	dbClient = db
	collection = dbClient.Database("qprobot").Collection("users")
}

//CheckUser - проверяет есть ли такой бользователь и если нет - создаёт его
func CheckUser(ID int64) (TgUser, error) {
	var user TgUser
	filter := bson.M{"ID": ID}
	err := collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		log.Printf("User %d not found. Create", ID)
		_, err := collection.InsertOne(context.TODO(), bson.M{"ID": ID})
		if err != nil {
			return user, err
		}
		user.ID = ID
		return user, nil
	}
	return user, nil
}

//ResetAd -
func (u *TgUser) ResetAd() {
	u.AddingAd = primitive.ObjectID{0}
	u.Save()
}

//SetState - становить новое состояние state
func (u *TgUser) SetState(state string) {
	if u.State == state {
		//менять на саму себя не будем т.к. итак тут.
		return
	}
	u.StateBack = u.State
	u.State = state
	u.Save()
}

//GoBack - перейти на предыдущее состояние
func (u *TgUser) GoBack() {
	u.State = u.StateBack
	u.StateBack = "root" //предыдущее авыставить root
	u.Save()
}
