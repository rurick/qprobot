package main

//answerCallbackQuery
//git clone https://github.com/go-telegram-bot-api/telegram-bot-api.git tgbotapi
//git checkout 99b74b8efaa519636cf7f56afed97b65ecafb512
// GetFileDirectURL

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"time"

	"github.com/tgbotapi"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//Config - конфиг
type Config struct {
	TelegramBotToken string
	DebugLog         bool
}

var configuration Config

//инициализация приложения
func init() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	file, _ := os.Open(path.Join(dir, "config.json"))
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration = Config{}
	err = decoder.Decode(&configuration)
	if err != nil {
		log.Panic(err)
	}
}

func main() {
	log.Println("Starting...")

	//создаем контексты
	ctxDB, canselDB := context.WithTimeout(context.Background(), 3*time.Second)
	ctx, cansel := context.WithCancel(context.Background()) //контекст программы для обработки событий
	go handleSignals(cansel)

	//инициализирую БД
	log.Print("Opening DB connection... ")
	dbClient, err := mongo.NewClient(options.Client().ApplyURI("mongodb://qprobot:rcode286@rurick.ru:27017"))
	if err != nil {
		log.Fatal(err)
	}
	defer canselDB()
	err = dbClient.Connect(ctxDB)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("OK")
	defer func() {
		log.Println("Close DB connection")
		dbClient.Disconnect(ctxDB)
	}()

	// TELEGRAM BOT
	// используя токен создаем новый инстанс бота
	bot, err := tgbotapi.NewBotAPI(configuration.TelegramBotToken)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = false
	log.Println("Authorized Telegram BOT on account ", bot.Self.UserName)

	// u - структура с конфигом для получения апдейтов
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	// используя конфиг u создаем канал в который будут прилетать новые сообщения
	updates := bot.GetUpdatesChan(u)
	log.Println("OK")

	//основной цикл чтения каналов
	for {
		select {
		case <-ctx.Done(): //ждем завершения программы по ctrl+c
			log.Println("Stopping programm...")
			return

		case update := <-updates: //ждем сообщение из телеги
			go Run(&update, bot, dbClient)
		}
	}
}

//обработчик сигналов
func handleSignals(cansel context.CancelFunc) {
	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, os.Interrupt)
	for {
		sig := <-sigCh
		switch sig {
		case os.Interrupt:
			cansel()
			return
		}
	}
}
