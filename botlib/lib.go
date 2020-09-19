package botlib

import (
	"strconv"
	"strings"

	"github.com/tgbotapi"
)

//DeleteIncomingMsg  удалит входящее сообщение
func DeleteIncomingMsg(update *tgbotapi.Update, chatID int64, bot *tgbotapi.BotAPI) {
	if update.Message != nil {
		del := tgbotapi.NewDeleteMessage(chatID, update.Message.MessageID)
		bot.Send(del)
	}
}

/*Cut * Обрезка строки
 * возвращает новую строку*/
func Cut(text string, limit int) string {
	runes := []rune(text)
	if len(runes) >= limit {
		return string(runes[:limit])
	}
	return text
}

//PriceFormatted -
func PriceFormatted(price int64) string {
	res := ""
	for o := price % 1000; price > 1000; {
		res = " " + strconv.Itoa(int(o)) + res
		price = price / 1000
	}
	return strings.Trim(strconv.Itoa(int(price))+res, " ")
}
