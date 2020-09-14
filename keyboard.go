package main

import "github.com/tgbotapi"

const (
	//ButtonAds -
	ButtonAds = "🔎 Объявления"
	//ButtonProfile -
	ButtonProfile = "🔑 Личный кабинет"
	//ButtonFavorits -
	ButtonFavorits = "❤️ Избранное"
)

//keyBoard - Клавиатура
var keyBoardRoot = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(ButtonAds),
		tgbotapi.NewKeyboardButton(ButtonProfile),
	),
)
