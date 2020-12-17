package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/go-pg/pg/v9"
	tb "gopkg.in/tucnak/telebot.v2"
)

func main() {
	var (
		port      = os.Getenv("PORT")
		publicURL = os.Getenv("PUBLIC_URL")
		token     = os.Getenv("TOKEN")
	)

	poller := &tb.Webhook{
		Listen:   ":" + port,
		Endpoint: &tb.WebhookEndpoint{PublicURL: publicURL},
	}

	b, err := tb.NewBot(tb.Settings{
		Token:     token,
		Poller:    poller,
		ParseMode: tb.ModeMarkdownV2,
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	db := &pg.DB{}
	ConnectDataBase(db)
	defer db.Close()

	//КОМАНДЫ
	b.Handle("/start", func(m *tb.Message) {

		referral, err := strconv.Atoi(m.Payload)
		if err != nil {
			referral = AdminBot
		}
		u, isNewUser := NewDefaultUser(db, m.Sender.ID, referral)

		if isNewUser {
			fmt.Printf("Новый игрок: @%s[%d]\n", m.Sender.Username, u.ID)
		}
		if err != nil {
			fmt.Println(err)
		}
		b.Send(m.Sender, GetTextFile("start"), ReplyMain)
	})

	b.Handle("/test", func(m *tb.Message) {
		b.Send(m.Sender, GetTextFile("test"))
	})

	//РЕПЛИКЕЙБОРДЫ
	b.Handle(&RBtnRead, func(m *tb.Message) {
		b.Send(m.Sender, "Список")
	})

	b.Handle(&RBtnWrite, func(m *tb.Message) {
		u := GetUser(db, m.Sender.ID)

		u.SetBotState(db, WriteSetNameState)
		b.Send(m.Sender, WriteSetNameState.Message())
	})

	b.Handle(tb.OnText, func(m *tb.Message) {
		u := GetUser(db, m.Sender.ID)

		u.BotState.Action()
	})

	b.Start()
}
