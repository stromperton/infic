package main

import (
	"fmt"
	"os"

	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
)

func createSchema(db *pg.DB) error {
	for _, model := range []interface{}{&User{}} {
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
	u := &User{}
	u.ID = id
	err := db.Select(u)
	if err != nil {
		fmt.Println(err)
	}

	return *u
}

//UpdateUser Обновление игрока из базы данных по ID
func UpdateUser(u User) {
	_, err := db.Model(u).WherePK().Update()
	if err != nil {
		fmt.Println(err)
	}
}

//NewDefaultUser Новый стандартный игрок
func NewDefaultUser(id int, ref int) (User, bool) {
	u := &User{}
	u.ID = id
	u.Ref = ref
	u.Lang = "rus"
	u.BotState = DefaultState

	res, err := db.Model(u).OnConflict("DO NOTHING").Insert()
	if err != nil {
		panic(err)
	}

	if res.RowsAffected() > 0 {
		return *u, true
	}
	return *u, false
}

//CreateStory Создать историю
func CreateStory(name string) {
	infic := Infic{
		Name: name,
	}

	_, err := db.Model(infic).Insert()
	if err != nil {
		panic(err)
	}
}
