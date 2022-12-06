package handlers

import (
	"database/sql"
	"fmt"
	tele "gopkg.in/telebot.v3"
	"log"
	"strings"
	"studyAndRepeat/src/database"
	"time"
)

func AddCard(db *sql.DB) tele.HandlerFunc {
	return func(c tele.Context) error {
		frontText := c.Message().Payload
		id, err := database.FindCardByFront(db, c.Sender().ID, frontText)
		if err != nil {
			return err
		}
		if id == 0 {
			id, err = database.StoreCard(db, c.Sender().ID, frontText)
		}
		if err != nil {
			return err
		}

		return c.Send("*"+frontText+"* was added!\nNow write the definition", tele.ModeMarkdown)
	}
}

func DeleteCard(db *sql.DB) tele.HandlerFunc {
	return func(c tele.Context) error {
		frontText := c.Message().Payload
		id, err := database.FindCardByFront(db, c.Sender().ID, frontText)
		if err != nil {
			return err
		}
		if id == 0 {
			return c.Send("Hmm 🤨... Such cards wasn't found in your dictionary")
		}
		err = database.DeleteCard(db, c.Sender().ID, frontText)
		if err != nil {
			return err
		}

		return c.Send("*"+frontText+"* was deleted", tele.ModeMarkdown)
	}
}

func SetDefinition(db *sql.DB) tele.HandlerFunc {
	return func(c tele.Context) error {
		backText := c.Message().Text
		id, frontText, err := database.FindLatestUserCard(db, c.Sender().ID)
		log.Printf("FindLatestUserCard: %d, %s", id, frontText)
		if err != nil {
			return err
		}
		if id == 0 {
			return c.Send("Hmm 🤨... It seems you need to use /add <phrase> beforehand")
		}

		err = database.SetBackForCard(db, c.Sender().ID, backText, id)
		if err != nil {
			return err
		}

		return c.Send("✅Card: *"+frontText+" - "+backText+"* was completed", tele.ModeMarkdown)
	}
}

func GetDictionary(db *sql.DB) tele.HandlerFunc {
	return func(c tele.Context) error {
		cards, err := database.FindUserCards(db, c.Sender().ID)
		if err != nil {
			return err
		}

		results := []string{}
		for k, card := range cards {
			results = append(results, fmt.Sprintf("%d. *%s - %s*. Repeat %s",
				k+1, card.Front, card.Back.String, readableAfter(card.RepeatAfter)))
		}

		return c.Send(strings.Join(results, "\n"), tele.ModeMarkdown)
	}
}

func readableAfter(after sql.NullTime) string {
	if !after.Valid {
		return "available"
	}

	out := time.Time{}.Add(after.Time.Sub(time.Now()))
	return "in " + out.Format("15h 04m")
}
