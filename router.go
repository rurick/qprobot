package main

import (
	"strings"

	"./ads"
	"./tguser"
	"github.com/tgbotapi"
)

//RouterCalback выполняем маршрутизация запроса
func RouterCalback(update *tgbotapi.Update, chatID int64, user *tguser.TgUser, bot *tgbotapi.BotAPI) {
	path := strings.Split(user.State, "/")
	if len(path) > 1 {
		switch path[1] {
		case "ads":
			ads.RouterCalback(update, chatID, user, bot)
		}
	}
}

//Router выполняем маршрутизация запроса
func Router(update *tgbotapi.Update, chatID int64, user *tguser.TgUser, bot *tgbotapi.BotAPI) {
	path := strings.Split(user.State, "/")
	if len(path) == 1 {
		if !actStart(update, chatID, user, bot) {
			pageStart(update, chatID, user, bot)
		}
	}
	if len(path) > 1 {
		switch path[1] {
		case "ads":
			ads.Router(update, chatID, user, bot)
		}
	}

}

/** надо сделать стартовое сообщение с инлайн кнопкой = append(/** надо сделать стартовое сообщение с инлайн кнопкой,
и после клика по кнопке открытие главной стницы которую уже можно праить
сообщение в котором клавиатура(большая) нельзя редактировать.*/

//pageStart  выводит первую страницу по команде start
func pageStart(update *tgbotapi.Update, chatID int64, user *tguser.TgUser, bot *tgbotapi.BotAPI) {
	user.Username = update.Message.From.UserName
	user.Name = update.Message.From.FirstName
	user.SetState("root")
	text := "<b>Я КуПи-Робот!</b>\n\n" +
		"🤖Я помогу Тебе быстро и выгодно купить или продать ВСЁ что угодно!\n" +
		"Мы находимся в телеграме, где все данные хранятся и передаются в надежно зашифрованом виде😎!\n\n" +
		"<i>Если Ты вдруг потеряешься, набери команду /start и Ты вернёшся на эту страницу</i>\n"

	reply := tgbotapi.NewMessage(chatID, "")
	reply.ParseMode = "HTML"
	reply.ReplyMarkup = keyBoardRoot
	reply.Text = text
	m, _ := bot.Send(reply)
	user.SetLastMsgID(m.MessageID)
}

func actStart(update *tgbotapi.Update, chatID int64, user *tguser.TgUser, bot *tgbotapi.BotAPI) bool {
	// Клик по кнопке обхъявления
	if update.Message.Text == ButtonAds {
		user.SetState("root/ads/0")
		ads.PageStart(update, chatID, user, bot, 0)
		return true
	}
	return false
}
