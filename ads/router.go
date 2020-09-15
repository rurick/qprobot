package ads

import (
	"sort"
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
			PageStart(update, chatID, user, bot)
		}
	}
	if len(path) > 2 {
	}
}

func actStart(update *tgbotapi.Update, chatID int64, user *tguser.TgUser, bot *tgbotapi.BotAPI) bool {
	if update.Message.Text == "🔎 Объявления" {
		del := tgbotapi.NewDeleteMessage(chatID, update.Message.MessageID)
		bot.Send(del)
		return true
	}
	return false

}

//PageStart - первая страница объявлений с рубриками
func PageStart(update *tgbotapi.Update, chatID int64, user *tguser.TgUser, bot *tgbotapi.BotAPI) {
	text := "<b>Объявления</b>\n\n" +
		strconv.Itoa(int(ad.Count("total"))) + " объявлений\n"
	reply := tgbotapi.NewMessage(chatID, text)
	reply.ParseMode = "HTML"

	//keyboard
	var inlineKB [][]tgbotapi.InlineKeyboardButton
	keys := make([]string, 0, len(adsCategories))
	for k := range adsCategories {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for k := 0; k < len(keys); k++ {
		row := tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				keys[k]+" (...)",
				"switchCategory:"+keys[k],
			),
		)
		inlineKB = append(inlineKB, row)
	}
	row := tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(
		"Избранное ❤️",
		"favorits:",
	))
	inlineKB = append(inlineKB, row)
	reply.ReplyMarkup = tgbotapi.InlineKeyboardMarkup{InlineKeyboard: inlineKB}
	//--

	m, _ := bot.Send(reply)
	user.SetLastMsgID(m.MessageID)

	//запуск задачи на расчёт  количества объявлений
	go func(msgID int, chatID int64, user *tguser.TgUser, bot *tgbotapi.BotAPI) {
		//keyboard
		var inlineKB [][]tgbotapi.InlineKeyboardButton
		keys := make([]string, 0, len(adsCategories))
		for k := range adsCategories {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for k := 0; k < len(keys); k++ {
			c := int(ad.Count(keys[k]))
			if c == 0 { //не выводить пустые рубрики
				//continue
			}
			row := tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(
					keys[k]+" ("+strconv.Itoa(c)+")",
					"switchCategory:"+keys[k],
				),
			)
			inlineKB = append(inlineKB, row)
		}
		row := tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(
			"Избранное ❤️",
			"favorits:",
		))
		inlineKB = append(inlineKB, row)
		reply := tgbotapi.NewEditMessageReplyMarkup(chatID, msgID, tgbotapi.InlineKeyboardMarkup{InlineKeyboard: inlineKB})
		bot.Send(reply)
	}(m.MessageID, chatID, user, bot)
}
