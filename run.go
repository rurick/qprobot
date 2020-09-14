package main

import (
	"encoding/json"
	"log"

	"./ad"
	"./tguser"
	"github.com/tgbotapi"
	"go.mongodb.org/mongo-driver/mongo"
)

//Run - новый поток обработчика
func Run(update *tgbotapi.Update, bot *tgbotapi.BotAPI, dbClient *mongo.Client) {
	tguser.Init(dbClient)
	ad.Init(dbClient)

	var (
		chatID int64
		inMsg  string
		UserID int64
	)
	if update.Message != nil {
		chatID = update.Message.Chat.ID
		inMsg = update.Message.Text
		UserID = int64(update.Message.From.ID)
	}
	if update.CallbackQuery != nil {
		chatID = update.CallbackQuery.Message.Chat.ID
		inMsg = update.CallbackQuery.Data
		UserID = int64(update.CallbackQuery.From.ID)

	}
	var reply interface{}

	user, err := tguser.CheckUser(UserID)
	if err != nil {
		log.Printf("CRITICAL ERROR: %s", err)
	}

	log.Println(inMsg)

	if update.Message.Command() == "start" {
		reply = pageStart(update, chatID, &user, bot)
	}

	//Отправка сообщения
	if reply != nil {
		m, err := bot.Send(reply.(tgbotapi.Chattable))
		user.SetLastMsgID(int64(m.MessageID))
		if err != nil {
			log.Printf("CRITICAL ERROR: %s", err)
		}
	}

}

func pageStart(update *tgbotapi.Update, chatID int64, user *tguser.TgUser, bot *tgbotapi.BotAPI) interface{} {
	user.Username = update.Message.From.UserName
	user.Name = update.Message.From.FirstName
	user.SetState("root")
	text := "<b>Я КуПи-Робот!</b>\n\n" +
		"🤖Я помогу Тебе быстро и выгодно купить или продать ВСЁ что угодно!\n" +
		"Мы находимся в телеграме, где все данные хранятся и передаются в надежно зашифрованом виде😎!\n\n" +
		"<i>Если Ты вдруг потеряешься, набери команду /start и Ты вернёшся на эту страницу</i>\n"

	lastMsgID := user.GetLastMsgID()
	if lastMsgID >= 0 {
		reply := NewEditMessageMarkup(chatID, int(lastMsgID), keyBoardRoot)
		reply.ParseMode = "HTML"
		reply.Text = text
		v, err := reply.values()
		if err != nil {
			log.Printf("%s", err)
		}
		resp, err := bot.MakeRequest(reply.method(), v)
		if err != nil {
			log.Printf("%s", err)
		}
		var message tgbotapi.Message
		json.Unmarshal(resp.Result, &message)
		return nil
	}
	reply := tgbotapi.NewMessage(chatID, "")
	reply.ParseMode = "HTML"
	reply.ReplyMarkup = keyBoardRoot
	reply.Text = text
	return reply
}
