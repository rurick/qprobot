package main

import (
	"encoding/json"
	"net/url"
	"strconv"

	"github.com/tgbotapi"
)

//EditMessageMarkup -
type EditMessageMarkup struct {
	tgbotapi.BaseEdit
	ReplyMarkup *tgbotapi.ReplyKeyboardMarkup
	ParseMode   string
	Text        string
}

// NewEditMessageMarkup allows you to edit the inline
// keyboard markup.
func NewEditMessageMarkup(chatID int64, messageID int, markup tgbotapi.ReplyKeyboardMarkup) EditMessageMarkup {
	return EditMessageMarkup{
		BaseEdit: tgbotapi.BaseEdit{
			ChatID:    chatID,
			MessageID: messageID,
		},
		ReplyMarkup: &markup,
	}
}

func (config EditMessageMarkup) values() (url.Values, error) {
	v := url.Values{}

	if config.ChannelUsername != "" {
		v.Add("chat_id", config.ChannelUsername)
	} else {
		v.Add("chat_id", strconv.FormatInt(config.ChatID, 10))
	}
	v.Add("message_id", strconv.Itoa(config.MessageID))

	if config.ReplyMarkup != nil {
		data, err := json.Marshal(config.ReplyMarkup)
		if err != nil {
			return v, err
		}
		v.Add("reply_markup", string(data))
	}

	return v, nil
}

func (config EditMessageMarkup) method() string {
	return "EditMessageMarkup"
}
