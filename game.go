package main

import (
	"fmt"

	tb "gopkg.in/tucnak/telebot.v2"
)

//User Игрок
type User struct {
	ID                int
	Name              string
	Keys              int `pg:"keys,use_zero,notnull"`
	Ref               int `pg:"ref,use_zero,notnull"`
	EditableInficID   int
	EditableMessageID int      `pg:"editable_message_id,use_zero,notnull"`
	BotState          BotState `pg:"bot_state,use_zero,notnull"`
	Library           []InficMeta
}

//Infic Интерактивный рассказ
type Infic struct {
	ID          int
	Name        string
	Description string
	Image       string
	AuthorID    int
	Author      *User `pg:"rel:has-one"`
	Story       map[int]Message
	isPublic    bool
}

//InficMeta Данные о закладках
type InficMeta struct {
	InficID int
	Level   int
	Branch  int
}
type Message struct {
	ID     int
	Title  string
	Text   string
	Parent int
	Childs []int
}

func (u *User) Action(message *tb.Message) string {
	infic := &Infic{ID: u.EditableInficID}
	err := db.Model(infic).WherePK().Select()
	if err != nil {
		fmt.Println(err)
	}

	switch u.BotState {
	case EditNameState:
		infic.Name = message.Text
	case EditDescriptionState:
		infic.Description = message.Text
	case EditImageState:
		if message.Photo != nil {
			infic.Image = message.Photo.FileID
		} else {
			return "Для обложки подойдут только <b>сжатые изображения</b> формата .png/.jpg/.jpeg"
		}
	case EditTextState:
		mess := infic.Story[u.EditableMessageID]
		mess.Text = message.Text
		infic.Story[u.EditableMessageID] = mess
	case EditTitleState:
		mess := infic.Story[u.EditableMessageID]
		mess.Title = message.Text
		infic.Story[u.EditableMessageID] = mess
	}
	UpdateModel(infic)
	return ""
}

//isInLibrary В библиотеке?
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
func (u *User) SetEditableMessage(id int) {
	u.EditableMessageID = id
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

//GetList Список инфиков, из библиотеки
func (u *User) GetList(order string) []Infic {
	infArr := &[]Infic{}
	err := db.Model(infArr).Relation("Author").Order(order).Select()
	if err != nil {
		fmt.Println(err)
	}
	return *infArr
}

func (infic *Infic) AddNewMessage(editableMessageID int) int {
	id := len(infic.Story)
	message := infic.Story[editableMessageID]
	message.Childs = append(message.Childs, id)
	infic.Story[editableMessageID] = message

	infic.Story[id] = Message{
		ID:     id,
		Title:  "Новое сообщение",
		Text:   "И вновь, и вновь...",
		Parent: editableMessageID,
	}

	UpdateModel(infic)
	return id
}
