package main

import (
	"strconv"
)

//User Игрок
type User struct {
	ID       int
	Ref      int `pg:"ref,use_zero,notnull"`
	Lang     string
	BotState BotState `pg:"bot_state,use_zero,notnull"`
}

//Infic Интерактивный рассказ
type Infic struct {
	Name string
}

//GetState Получить состояние
func (u *User) GetState() string {
	return "Всё хорошо!" + strconv.Itoa(u.ID)
}

func (u *User) SetBotState(newState BotState) {
	u.BotState = newState
	UpdateUser(*u)
}
