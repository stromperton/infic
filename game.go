package main

import (
	"strconv"
)

//User Игрок
type User struct {
	ID              int
	Ref             int `pg:"ref,use_zero,notnull"`
	Lang            string
	BotState        BotState `pg:"bot_state,use_zero,notnull"`
	Keys            int      `pg:"keys,use_zero,notnull"`
	EditableInficID int
}

//Infic Интерактивный рассказ
type Infic struct {
	ID          int
	Name        string
	Description string
	Image       string
	Story       [][]Message
	isPublic    bool
}

type Message struct {
	Text    string
	Level   int
	Version int
}

//GetState Получить состояние
func (u *User) GetState() string {
	return "Всё хорошо!" + strconv.Itoa(u.ID)
}

func (u *User) SetBotState(newState BotState) {
	u.BotState = newState
	UpdateModel(u)
}

func (u *User) AddKeys(num int) {
	u.Keys += num
	UpdateModel(u)
}

func (u *User) Action(text string) {
	switch u.BotState {
	case WriteSetNameState:
		u.EditableInficID = CreateInfic(text)
		u.SetBotState(WriteSetDescriptionState)
	case WriteSetDescriptionState:
		UpdateModel(&Infic{
			ID:          u.EditableInficID,
			Description: text,
		})
		u.SetBotState(DefaultState)
	}
}
