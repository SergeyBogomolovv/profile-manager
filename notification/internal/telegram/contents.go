package telegram

import (
	tele "gopkg.in/telebot.v4"
)

const (
	AccountTypeTelegram = "Телеграм"
	AccountTypeEmail    = "Почта"
)

var clearMenu = &tele.ReplyMarkup{
	RemoveKeyboard: true,
}

var accountTypeMenu = &tele.ReplyMarkup{
	ReplyKeyboard: [][]tele.ReplyButton{
		{{Text: AccountTypeEmail}, {Text: AccountTypeTelegram}},
	},
}
