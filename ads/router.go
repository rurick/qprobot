package ads

import (
	"log"
	"strconv"
	"strings"

	"../ad"
	"../tguser"

	"github.com/tgbotapi"
)

//Router выполняем маршрутизация запроса
func Router(update *tgbotapi.Update, chatID int64, user *tguser.TgUser, bot *tgbotapi.BotAPI) {
	path := strings.Split(user.State, "/")
	if len(path) == 2 {
		if !actStart(update, chatID, user, bot) {
			pageStart(update, chatID, user, bot)
		}
	}
	if len(path) > 2 {
	}
}

func actStart(update *tgbotapi.Update, chatID int64, user *tguser.TgUser, bot *tgbotapi.BotAPI) bool {
	return false
}

func pageStart(update *tgbotapi.Update, chatID int64, user *tguser.TgUser, bot *tgbotapi.BotAPI) {
	lastMessage := user.GetLastMsgID()
	text := "<b>Объявления</b>\n\n" +
		"Опубликовано " + strconv.Itoa(int(ad.Count("total"))) + "объявлений\n"
	var inlineKB [][]tgbotapi.InlineKeyboardButton
	for k := range adsCategories {
		inlineKB = append(inlineKB, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(k, "switchCategory:"+k),
		))
	}

	reply := tgbotapi.NewEditMessageText(chatID, lastMessage, text)
	m, err := bot.Send(reply)
	if err != nil {
		log.Println(err)
	}

	reply1 := tgbotapi.NewEditMessageReplyMarkup(chatID, m.MessageID, tgbotapi.InlineKeyboardMarkup{InlineKeyboard: inlineKB})
	m, _ = bot.Send(reply1)
	user.SetLastMsgID(m.MessageID)
}
