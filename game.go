package main

import (
	"fmt"
	"strconv"

	tb "gopkg.in/tucnak/telebot.v2"
)

//User Игрок
type User struct {
	ID              int
	Name            string
	Keys            int `pg:"keys,use_zero,notnull"`
	Ref             int `pg:"ref,use_zero,notnull"`
	EditableInficID int
	BotState        BotState `pg:"bot_state,use_zero,notnull"`
	Library         []InficMeta
}

//Infic Интерактивный рассказ
type Infic struct {
	ID          int
	Name        string
	Description string
	Image       string
	AuthorID    int
	Author      *User `pg:"rel:has-one"`
	Story       [][]Message
	isPublic    bool
}

//InficMeta Данные о закладках
type InficMeta struct {
	InficID int
	Level   int
	Branch  int
}
type Message struct {
	Text   string
	Level  int
	Branch int
}

func (u *User) Action(message *tb.Message) {
	switch u.BotState {
	case EditNameState:
		UpdateModel(&Infic{
			ID:   u.EditableInficID,
			Name: message.Text,
		})
	case EditDescriptionState:
		UpdateModel(&Infic{
			ID:          u.EditableInficID,
			Description: message.Text,
		})
	case EditImageState:
		UpdateModel(&Infic{
			ID:    u.EditableInficID,
			Image: message.Photo.FileID,
		})
	}
	u.SetBotState(DefaultState)
}

//GetState Получить состояние
func (u *User) GetState() string {
	return "Всё хорошо!" + strconv.Itoa(u.ID)
}

//GetState Получить состояние
func (u *User) isInLibrary(inficID int) bool {
	inLibrary := false

	for _, inf := range u.Library {
		if inf.InficID == inficID {
			inLibrary = true
		}
	}

	return inLibrary
}

func (u *User) SetBotState(newState BotState) {
	u.BotState = newState
	UpdateModel(u)
}

func (u *User) SetEditableInfic(id int) {
	u.EditableInficID = id
	UpdateModel(u)
}

func (u *User) AddKeys(num int) {
	u.Keys += num
	UpdateModel(u)
}

//GetMyWorks Список моих инфиков
func (u *User) GetMyWorks() []Infic {
	infArr := &[]Infic{}
	err := db.Model(infArr).Relation("Author").Where("author.id = ?", u.ID).Select()
	if err != nil {
		fmt.Println(err)
	}
	return *infArr
}

//GetMyLibrary Список инфиков, из библиотеки
func (u *User) GetMyLibrary(order string) []Infic {
	infArr := &[]Infic{}
	err := db.Model(infArr).Order(order).Select()
	if err != nil {
		fmt.Println(err)
	}
	return *infArr
}
