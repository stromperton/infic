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

	b.Handle(tb.OnText, func(m *tb.Message) {
		u := GetUser(m.Sender.ID)

		if m.Text[:2] == "/i" {
			id, _ := strconv.Atoi(m.Text[2:])
			message, err := SprintInfic(id, b)
			if err != nil {
				b.Send(m.Sender, "Инфик не существует...")
			} else {
				b.Send(m.Sender, message)
			}
		} else {

			u.Action(m.Text)
			b.Send(m.Sender, u.BotState.Message())
		}
	})

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

		message := fmt.Sprintf("Твои инфики:")
		myInfics := u.GetMyInfics()

		for _, inf := range myInfics {

			message += fmt.Sprintf("\n*/i%d %s*", inf.ID, inf.Name)
		}

		b.Send(m.Sender, message, InlineWhrite)
	})

	b.Handle(&RBtnAccount, func(m *tb.Message) {
		u := GetUser(m.Sender.ID)

		u.SetBotState(AccountCheckState)
		b.Send(m.Sender, AccountCheckState.Message())
	})

	//ИНЛИНЕКЕЙБОРДЫ
	b.Handle(&IBtnCreare, func(c *tb.Callback) {
		b.Respond(c)
		u := GetUser(c.Sender.ID)
		message, _ := SprintInfic(CreateInfic(u.ID), b)
		b.Send(c.Sender, message)
	})

	b.Start()
}

func SprintInfic(id int, b *tb.Bot) (tb.Sendable, error) {
	inf, err := GetInfic(id)

	var file tb.File
	if inf.Image == "" {
		file = tb.FromDisk("pustota.jpg")
	} else {
		file, _ = b.FileByID(inf.Image)
	}

	sendable := &tb.Photo{
		File: file,
		Caption: fmt.Sprintf(`*%s*
_%s_`, inf.Name, inf.Description),
	}
	return sendable, err
}
