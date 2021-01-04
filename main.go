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
		ParseMode: tb.ModeHTML,
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	ConnectDataBase()
	defer db.Close()

	b.Handle(tb.OnPhoto, func(m *tb.Message) {
		u := GetUser(m.Sender.ID)

		if u.BotState == EditImageState {
			u.Action(m)
		} else {
			b.Send(m.Sender, "А зачем мне сейчас эта фотография?")
		}
	})

	b.Handle(tb.OnText, func(m *tb.Message) {
		u := GetUser(m.Sender.ID)

		var sendable interface{}
		var keyboard *tb.ReplyMarkup
		var id int

		if u.BotState == DefaultState && m.Text[:2] == "/i" {
			id, _ = strconv.Atoi(m.Text[2:])
			u.SetEditableInfic(id)
		} else {
			u.Action(m)
			id = u.EditableInficID
		}

		var aid int
		sendable, aid, err = SprintInfic(id, b)
		if err != nil {
			b.Send(m.Sender, "Инфик не существует...")
		} else {
			keyboard = InlineInfic
			if m.Sender.ID == aid {
				keyboard = InlineInficEdit
			} else if u.isInLibrary(id) {
				keyboard = InlineInficWithRemove
			}

		}
		b.Send(m.Sender, sendable, keyboard)
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

	libraryFunc := func(m *tb.Message) {
		u := GetUser(m.Sender.ID)

		message := fmt.Sprintf("📚 <b>Твоя библиотека</b>")
		myInfics := u.GetMyLibrary("name ASC")

		for _, inf := range myInfics {

			message += fmt.Sprintf("\n<b>/i%d %s</b> - %s", inf.ID, inf.Name, inf.Author.Name)
		}

		b.Send(m.Sender, message, InlineRead)
	}

	//РЕПЛИКЕЙБОРДЫ
	b.Handle(&RBtnRead, libraryFunc)

	b.Handle(&RBtnWrite, func(m *tb.Message) {
		u := GetUser(m.Sender.ID)

		message := fmt.Sprintf("✍️ <b>Твои рукописи</b>")
		myInfics := u.GetMyWorks()

		for _, inf := range myInfics {

			message += fmt.Sprintf("\n<b>/i%d %s</b> - %s", inf.ID, inf.Name, inf.Author.Name)
		}

		b.Send(m.Sender, message, InlineWhrite)
	})

	b.Handle(&RBtnAccount, func(m *tb.Message) {
		u := GetUser(m.Sender.ID)

		u.SetBotState(AccountCheckState)
		b.Send(m.Sender, AccountCheckState.Message())
	})

	//ИНЛИНЕКЕЙБОРДЫ
	b.Handle(&IBtnCreate, func(c *tb.Callback) {
		b.Respond(c)
		u := GetUser(c.Sender.ID)
		message, _, _ := SprintInfic(CreateInfic(u.ID), b)
		_, err := b.Send(c.Sender, message, InlineInficEdit)
		fmt.Println(err)
	})

	b.Handle(&IBtnRead, func(c *tb.Callback) {
		b.Respond(c)
		u := GetUser(c.Sender.ID)

		if u.isInLibrary(u.EditableInficID) {

		} else {
			meta := InficMeta{InficID: u.EditableInficID, Level: 0, Branch: 0}
			u.Library = append(u.Library, meta)
		}
		message := "Текст"
		_, err := b.Send(c.Sender, message, InlineInficRead)
		fmt.Println(err)
	})

	b.Handle(&IBtnAllListAZ, func(m *tb.Message) {
		u := GetUser(m.Sender.ID)

		message := fmt.Sprintf("<b>Все инфики по алфавиту</b>")
		myInfics := u.GetMyLibrary("name ASC")

		for _, inf := range myInfics {

			message += fmt.Sprintf("\n<b>/i%d %s</b> - %s", inf.ID, inf.Name, inf.Author.Name)
		}

		b.Send(m.Sender, message, InlineRead)
	})

	b.Handle(&IBtnAllListID, func(m *tb.Message) {
		u := GetUser(m.Sender.ID)

		message := fmt.Sprintf("<b>Все инфики по ID</b>")
		myInfics := u.GetMyLibrary("id ASC")

		for _, inf := range myInfics {

			message += fmt.Sprintf("\n<b>/i%d %s</b> - %s", inf.ID, inf.Name, inf.Author.Name)
		}

		b.Send(m.Sender, message, InlineRead)
	})

	b.Handle(&IBtnMyLibrary, libraryFunc)

	for i := BotState(0); i < EndEnumState; i++ {
		b.Handle(i.Endpoint(), func(c *tb.Callback) {
			b.Respond(c)
			u := GetUser(c.Sender.ID)

			u.SetBotState(i)
			_, err := b.Send(c.Sender, i.Message())
			fmt.Println(err)
		})
	}

	b.Start()
}

func SprintInfic(id int, b *tb.Bot) (*tb.Photo, int, error) {
	inf, err := GetInfic(id)

	var file tb.File
	if inf.Image == "" {
		file = tb.FromDisk("pustota.jpg")
	} else {
		file, _ = b.FileByID(inf.Image)
	}

	sendable := &tb.Photo{
		File: file,
		Caption: fmt.Sprintf(`<b>%s</b>
%s
<i>%s</i>`, inf.Name, inf.Description, inf.Author.Name),
	}
	return sendable, inf.AuthorID, err
}
