package main

import (
	tb "gopkg.in/tucnak/telebot.v2"
)

//BotState Состояние бота
type BotState int

const (
	//DefaultState Начальное состояние бота
	DefaultState BotState = iota
	//MainMenuState Главное меню
	MainMenuState
	WriteSetNameState
	WriteSetDescriptionState
	AccountCheckState
)

func (d BotState) String() string {
	return [...]string{"Default", "MainMenu", "WriteSetName", "AccountCheckState"}[d]
}
func (d BotState) Message() string {
	return [...]string{
		"Default",
		"MainMenu",
		"Введите название вашего инфика",
		`🗝 <b>Аккаунт</b>
		...`,
	}[d]
}

//AdminBot Telegram ID
const AdminBot int = 303629013

var (
	RBtnRead    = tb.ReplyButton{Text: "📕 Читать"}
	RBtnAccount = tb.ReplyButton{Text: "🗝 Аккаунт"}
	RBtnWrite   = tb.ReplyButton{Text: "✍️ Писать"}

	//ReplyMain Главное меню бота
	ReplyMain = &tb.ReplyMarkup{
		ResizeReplyKeyboard: true,
		ReplyKeyboard:       [][]tb.ReplyButton{{RBtnRead}, {RBtnAccount, RBtnWrite}},
	}
)
