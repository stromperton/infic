package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/go-pg/pg/v9"
	tb "gopkg.in/tucnak/telebot.v2"
)

var db *pg.DB

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

	ConnectDataBase()
	defer db.Close()

	//КОМАНДЫ
	b.Handle("/start", func(m *tb.Message) {

		referral, err := strconv.Atoi(m.Payload)
		if err != nil {
			referral = AdminBot
		}
		u, isNewUser := NewDefaultUser(m.Sender.ID, referral)

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
		u := GetUser(m.Sender.ID)

		u.GetMyInfics()

		message := fmt.Sprintf("Твои инфики:")
		b.Send(m.Sender, message)
	})

	b.Handle(&RBtnAccount, func(m *tb.Message) {
		u := GetUser(m.Sender.ID)

		u.SetBotState(AccountCheckState)
		b.Send(m.Sender, AccountCheckState.Message())
	})

	b.Handle(tb.OnText, func(m *tb.Message) {
		u := GetUser(m.Sender.ID)

		if m.Text[:2] == "/i" {
			b.Send(m.Sender, "это инфик")
		} else {

			u.Action(m.Text)
			b.Send(m.Sender, u.BotState.Message())
		}
	})

	b.Start()
}
