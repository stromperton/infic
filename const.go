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
	EditNameState
	EditDescriptionState
	EditImageState
	AccountCheckState
	EditTextState
	EditTitleState
	EndEnumState
)

func (d BotState) String() string {
	return [...]string{
		"Default",
		"MainMenu",
		"EditNameState",
		"EditDescriptionState",
		"EditImageState",
		"AccountCheckState",
		"EditTextState",
		"EditTitleState",
		"EndEnumState",
	}[d]
}
func (d BotState) Message() string {
	return [...]string{
		"Default",
		"MainMenu",
		"Хорошо. Отправь мне новое <b>название</b> для этого инфика.",
		"Хорошо. Отправь мне новое <b>описание</b> для этого инфика.",
		"Хорошо. Отправь мне новую <b>обложку</b> для этого инфика.",
		`🗝 <b>Аккаунт</b>
...`,
		"Хорошо. Отправь новый <b>текст</b> этого сообщения",
		"Хорошо. Отправь новый <b>заголовок</b> для этого сообщения",
		"EndEnumState",
	}[d]
}

func (d BotState) Endpoint() string {
	return "\f" + [...]string{
		"",
		"",
		"editName",
		"editDesc",
		IBtnEditImage.Unique,
		"",
		"editMessageText",
		"editMessageTitle",
		"EndEnumState",
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

	IBtnCreate        = tb.InlineButton{Text: "✍️ Начать новый", Unique: "create"}
	IBtnRead          = tb.InlineButton{Text: "📕 Читать", Unique: "read"}
	IBtnAddLibrary    = tb.InlineButton{Text: "⭐️ Добавить в библиотеку", Unique: "addLibrary"}
	IBtnRemoveLibrary = tb.InlineButton{Text: "❌ Убрать из библиотеки", Unique: "addLibrary"}
	IBtnToList        = tb.InlineButton{Text: "⬅️ Назад к списку", Unique: "toList"}
	IBtnNext          = tb.InlineButton{Text: "📃 Далее", Unique: "next"}

	IBtnEdit      = tb.InlineButton{Text: "✍️ Редактировать", Unique: "edit"}
	IBtnPublic    = tb.InlineButton{Text: "📖 Опубликовать", Unique: "public"}
	IBtnEditName  = tb.InlineButton{Text: "Сменить название", Unique: "editName"}
	IBtnEditDesc  = tb.InlineButton{Text: "Сменить описание", Unique: "editDesc"}
	IBtnEditImage = tb.InlineButton{Text: "Сменить обложку", Unique: "editImage"}

	IBtnMyLibrary = tb.InlineButton{Text: "📚 Моя библиотека", Unique: "library"}
	IBtnAllListAZ = tb.InlineButton{Text: "Все инфики (A - Я)", Unique: "allListAZ"}
	IBtnAllListID = tb.InlineButton{Text: "Все инфики (ID)", Unique: "allListID"}
	IBtnRandom    = tb.InlineButton{Text: "🎲 Случайный", Unique: "random"}

	IBtnEditMessageText  = tb.InlineButton{Text: "✍️ Текст", Unique: "editMessageText"}
	IBtnEditMessageTitle = tb.InlineButton{Text: "✍️ Заголовок", Unique: "editMessageTitle"}
	IBtnNewMessage       = tb.InlineButton{Text: "+", Unique: "newMessage"}

	InlineWhrite = &tb.ReplyMarkup{
		InlineKeyboard: [][]tb.InlineButton{{IBtnCreate}},
	}

	InlineInficEdit = &tb.ReplyMarkup{
		InlineKeyboard: [][]tb.InlineButton{{IBtnRead, IBtnEdit}, {IBtnEditName, IBtnEditDesc}, {IBtnEditImage, IBtnPublic}},
	}

	InlineInfic = &tb.ReplyMarkup{
		InlineKeyboard: [][]tb.InlineButton{{IBtnRead}},
	}
	InlineInficWithRemove = &tb.ReplyMarkup{
		InlineKeyboard: [][]tb.InlineButton{{IBtnRead}, {IBtnRemoveLibrary}},
	}

	InlineRead = &tb.ReplyMarkup{
		InlineKeyboard: [][]tb.InlineButton{{IBtnAllListAZ, IBtnAllListID}, {IBtnMyLibrary, IBtnRandom}},
	}

	InlineInficRead = &tb.ReplyMarkup{
		InlineKeyboard: [][]tb.InlineButton{{IBtnNext}},
	}
)
