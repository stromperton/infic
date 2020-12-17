package main

import (
	"strconv"

	"github.com/go-pg/pg/v9"
)

//User Игрок
type User struct {
	ID       int
	Ref      int `pg:"ref,use_zero,notnull"`
	Lang     string
	BotState BotState
}

//Infic Интерактивный рассказ
type Infic struct {
	Name string
}

//GetState Получить состояние
func (u *User) GetState() string {
	return "Всё хорошо!" + strconv.Itoa(u.ID)
}

func (u *User) SetBotState(db *pg.DB, newState BotState) {
	u.BotState = newState
	UpdateUser(db, *u)
}
