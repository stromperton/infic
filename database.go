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

//GetInfic Извлечение игрока из базы данных по ID
func GetInfic(id int) Infic {
	inf := &Infic{ID: id}
	err := db.Model(inf).WherePK().Select()
	if err != nil {
		fmt.Println(err)
	}

	return *inf
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

//CreateInfic Создать историю
func CreateInfic(name string) int {
	infic := Infic{
		Name:     name,
		isPublic: false,
	}

	_, err := db.Model(infic).Insert()
	if err != nil {
		panic(err)
	}

	return infic.ID
}
