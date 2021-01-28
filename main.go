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

			var sendable interface{}
			var keyboard *tb.ReplyMarkup
			var id = u.EditableInficID

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

	allListFuncCallback := func(c *tb.Callback, order string, title string) {
		b.Respond(c)
		u := GetUser(c.Sender.ID)

		message := fmt.Sprintf("📚 <b>" + title + "</b>")
		myInfics := u.GetList(order)

		for _, inf := range myInfics {

			message += fmt.Sprintf("\n<b>/i%d %s</b> - %s", inf.ID, inf.Name, inf.Author.Name)
		}

		b.Edit(c.Message, message, InlineRead)
	}

	//РЕПЛИКЕЙБОРДЫ
	b.Handle(&RBtnRead, func(m *tb.Message) {
		u := GetUser(m.Sender.ID)

		message := fmt.Sprintf("📚 <b>Твоя библиотека</b>")
		myInfics := u.GetList("name ASC")

		for _, inf := range myInfics {

			message += fmt.Sprintf("\n<b>/i%d %s</b> - %s", inf.ID, inf.Name, inf.Author.Name)
		}

		b.Send(m.Sender, message, InlineRead)
	})

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

	//ИНЛИНЕКЕЙБОРДЫ ДЛЯ СОСТОЯНИЙ
	eeefunc := func(c *tb.Callback, state BotState) {
		b.Respond(c)
		u := GetUser(c.Sender.ID)

		u.SetBotState(state)
		_, err := b.Send(c.Sender, state.Message())
		fmt.Println(err)
	}
	b.Handle(EditNameState.Endpoint(), func(c *tb.Callback) {
		eeefunc(c, EditNameState)
	})
	b.Handle(EditDescriptionState.Endpoint(), func(c *tb.Callback) {
		eeefunc(c, EditDescriptionState)
	})
	b.Handle(EditImageState.Endpoint(), func(c *tb.Callback) {
		eeefunc(c, EditImageState)
	})
	b.Handle(EditTextState.Endpoint(), func(c *tb.Callback) {
		eeefunc(c, EditTextState)
	})
	b.Handle(EditTitleState.Endpoint(), func(c *tb.Callback) {
		eeefunc(c, EditTitleState)
	})

	//ИНЛИНЕКЕЙБОРДЫ
	b.Handle(&IBtnCreate, func(c *tb.Callback) {
		b.Respond(c)
		u := GetUser(c.Sender.ID)

		inficID := CreateInfic(u.ID)
		u.SetEditableInfic(inficID)
		message, _, _ := SprintInfic(inficID, b)
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

	b.Handle(&IBtnAllListAZ, func(c *tb.Callback) {
		allListFuncCallback(c, "name ASC", "Все инфики по алфавиту")
	})

	b.Handle(&IBtnAllListID, func(c *tb.Callback) {
		allListFuncCallback(c, "id ASC", "Все инфики по ID")
	})

	b.Handle(&IBtnMyLibrary, func(c *tb.Callback) {
		allListFuncCallback(c, "id ASC", "Моя библиотека")
	})

	b.Handle(&IBtnEdit, func(c *tb.Callback) {
		b.Respond(c)
		u := GetUser(c.Sender.ID)

		infic, _ := GetInfic(u.EditableInficID)

		m, k := GetMessageMessage(u, infic, 0)
		_, err := b.Send(c.Sender, m, k)
		fmt.Println(err)
	})

	b.Handle(&IBtnNewMessage, func(c *tb.Callback) {
		b.Respond(c)
		u := GetUser(c.Sender.ID)
		infic, _ := GetInfic(u.EditableInficID)

		newID := infic.AddNewMessage(u.EditableMessageID)

		m, k := GetMessageMessage(u, infic, newID)

		b.Edit(c.Message, m, k)
	})

	b.Handle("\fmessage", func(c *tb.Callback) {
		b.Respond(c)
		u := GetUser(c.Sender.ID)
		infic, _ := GetInfic(u.EditableInficID)
		id, _ := strconv.Atoi(c.Data)
		m, k := GetMessageMessage(u, infic, id)

		b.Edit(c.Message, m, k)
	})

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

func GetMessageMessage(u User, infic Infic, mID int) (string, *tb.ReplyMarkup) {
	u.SetEditableMessage(mID)
	thisMess := infic.Story[mID]

	var linkRow []tb.InlineButton
	var keyboardRows [][]tb.InlineButton

	parentMess := infic.Story[thisMess.Parent]
	keyboardRows = append(keyboardRows, []tb.InlineButton{{Text: parentMess.Title, Unique: "message", Data: fmt.Sprint(parentMess.ID)}})
	keyboardRows = append(keyboardRows, []tb.InlineButton{IBtnEditMessageText, IBtnEditMessageTitle})

	i := 0
	for _, num := range thisMess.Childs {
		linkRow = append(linkRow, tb.InlineButton{Text: infic.Story[num].Title, Unique: "message", Data: fmt.Sprint(infic.Story[num].ID)})
		i++

		if i > 3 {
			i = 0
			keyboardRows = append(keyboardRows, linkRow)
			linkRow = []tb.InlineButton{}
		}

	}
	keyboardRows = append(keyboardRows, linkRow)
	keyboardRows = append(keyboardRows, []tb.InlineButton{IBtnNewMessage})

	message := fmt.Sprintf("<b>Сообщение ID %d</b> <i>\"%s\"</i>\n%s", thisMess.ID, thisMess.Title, thisMess.Text)

	keyboard := &tb.ReplyMarkup{
		InlineKeyboard: keyboardRows,
	}
	return message, keyboard
}
