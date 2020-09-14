package main

import (
	"log"

	"go.mongodb.org/mongo-driver/mongo"

	"./ad"
	"./tguser"
	"github.com/tgbotapi"
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

	//Инициализация пользователя
	user, err := tguser.CheckUser(UserID)
	if err != nil {
		log.Printf("CRITICAL ERROR: %s", err)
	}

	log.Println(inMsg)

	//Проверка команд
	if update.Message.Command() == "start" {
		pageStart(update, chatID, &user, bot)
	} else {
		Router(update, chatID, &user, bot)
	}

}
