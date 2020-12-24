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
	WriteSetImageState
	AccountCheckState
)

func (d BotState) String() string {
	return [...]string{
		"Default",
		"MainMenu",
		"WriteSetName",
		"WriteSetDescription",
		"WriteSetImage",
		"AccountCheckState",
	}[d]
}
func (d BotState) Message() string {
	return [...]string{
		"Default",
		"MainMenu",
		"Как назовем инфик?",
		"Дайте краткое описание вашему инфику",
		"Отправьте фотографию для обложки",
		`🗝 *Аккаунт*
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

	IBtnCreare   = tb.InlineButton{Text: "Начать новый", Unique: "create"}
	InlineWhrite = &tb.ReplyMarkup{
		InlineKeyboard: [][]tb.InlineButton{{IBtnCreare}},
	}
)
