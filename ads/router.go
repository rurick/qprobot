package ads

import (
	"sort"
	"strconv"
	"strings"

	"../ad"
	"../botlib"
	"../tguser"

	"github.com/tgbotapi"
)

//RouterCalback выполняем маршрутизация запроса
func RouterCalback(update *tgbotapi.Update, chatID int64, user *tguser.TgUser, bot *tgbotapi.BotAPI) {
	switch update.CallbackQuery.Data {
	case "nextPage:":
		user.SetState("root/ads/1")
		Router(update, chatID, user, bot)
	case "prevPage:":
		user.SetState("root/ads/0")
		Router(update, chatID, user, bot)
	}
}

//Router выполняем маршрутизация запроса
func Router(update *tgbotapi.Update, chatID int64, user *tguser.TgUser, bot *tgbotapi.BotAPI) {
	path := strings.Split(user.State, "/")
	if path[1] == "ads" {
		botlib.DeleteIncomingMsg(update, chatID, bot)
		page, _ := strconv.Atoi(path[2])
		PageStart(update, chatID, user, bot, page)
	}
	if len(path) > 2 {
	}
}

//PageStart - первая страница объявлений с рубриками
func PageStart(update *tgbotapi.Update, chatID int64, user *tguser.TgUser, bot *tgbotapi.BotAPI, page int) {
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
	l := len(keys)
	for k := l / 2 * page; k < l/2*(page+1); k++ {
		row := tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				keys[k]+" (...)",
				"switchCategory:"+keys[k],
			),
		)
		inlineKB = append(inlineKB, row)
	}
	newxtbtn := tgbotapi.NewInlineKeyboardButtonData("след. стр. ⏩", "nextPage:")
	if page > 0 {
		newxtbtn = tgbotapi.NewInlineKeyboardButtonData("⏪ пред. стр.", "prevPage:")
	}
	row := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Избранное ❤️", "favorits:"),
		newxtbtn,
	)
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
		for k := l / 2 * page; k < l/2*(page+1); k++ {
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
		row := tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Избранное ❤️", "favorits:"),
			newxtbtn,
		)
		inlineKB = append(inlineKB, row)
		reply := tgbotapi.NewEditMessageReplyMarkup(chatID, msgID, tgbotapi.InlineKeyboardMarkup{InlineKeyboard: inlineKB})
		bot.Send(reply)
	}(m.MessageID, chatID, user, bot)
}
