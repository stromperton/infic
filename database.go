package main

import (
	"fmt"
	"os"

	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
)

func createSchema(db *pg.DB) error {
	for _, model := range []interface{}{&User{}, &Infic{}} {
		err := db.CreateTable(model, &orm.CreateTableOptions{
			IfNotExists: true,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

//ConnectDataBase Подключение к базе данных
func ConnectDataBase() {
	opt, err := pg.ParseURL(os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}

	db = pg.Connect(opt)

	err = createSchema(db)
	if err != nil {
		panic(err)
	}
	fmt.Println("Успешное подключение к базе данных")
}

//GetUser Извлечение игрока из базы данных по ID
func GetUser(id int) User {
	u := &User{ID: id}
	err := db.Model(u).WherePK().Select()
	if err != nil {
		fmt.Println(err)
	}

	return *u
}

//GetInfic Извлечение инфика из базы данных по ID
func GetInfic(id int) (Infic, error) {
	inf := &Infic{ID: id}
	err := db.Model(inf).Relation("Author").WherePK().Select()
	if err != nil {
		fmt.Println(err)
	}

	return *inf, err
}

//UpdateModel Обновление игрока из базы данных по ID
func UpdateModel(m interface{}) {
	_, err := db.Model(m).WherePK().Update()
	if err != nil {
		fmt.Println(err)
	}
}

//NewDefaultUser Новый стандартный игрок
func NewDefaultUser(id int, ref int) (User, bool) {
	u := &User{}
	u.ID = id
	u.Name = "Непризнанный гений"
	u.Ref = ref
	u.BotState = DefaultState
	u.Library = map[int]InficMeta{}

	res, err := db.Model(u).OnConflict("DO NOTHING").Insert()
	if err != nil {
		panic(err)
	}

	if res.RowsAffected() > 0 {
		return *u, true
	}
	return *u, false
}

//CreateInfic Создать историю
func CreateInfic(author int) int {
	infic := &Infic{
		Name:        "Новый инфик",
		Description: "Если вы когда-нибудь и видели новые инфики, то могу гарантировать, что этот будет свежее всех свежы́х!",
		isPublic:    false,
		AuthorID:    author,
	}
	infic.Story = map[int]Message{}
	infic.Story[0] = Message{
		ID:    0,
		Title: "Стартовое сообщение",
		Text:  fmt.Sprintf("Перед тобой минимальная функциональная единица инфика — обычное сообщение.\n\nТы можешь создавать новые части истории и перемещаться по уже созданным блокам с помощью кнопок внизу.\n\nЗаголовки сообщений предназначены для наименования выбора во время развилок, и для того чтобы тебе было легче ориентироваться в собственной рукописи.\n\nЭто начало твоей новой истории, просто начни творить!"),
	}

	_, err := db.Model(infic).Insert()
	if err != nil {
		panic(err)
	}

	return infic.ID
}
