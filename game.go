package main

import (
	"strconv"

	"github.com/go-pg/pg/v9"
)

var (
	//DefaultUser Стандартный игрок
	DefaultUser = &User{
		ID:       0,
		BotState: DefaultState,
	}
)

//User Игрок
type User struct {
	ID       int
	Ref      int
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
