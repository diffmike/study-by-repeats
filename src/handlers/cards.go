package handlers

import (
	"database/sql"
	"fmt"
	tele "gopkg.in/telebot.v3"
	"log"
	"strings"
	"studyAndRepeat/src/database"
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

		return c.Send("Nice!\n" + frontText + " was added\nNow write the definition")
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

		return c.Send(frontText + " was deleted")
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

		return c.Send("Thank you!\nCard: " + frontText + " - " + backText + " was completed")
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
			results = append(results, fmt.Sprintf("%d. %s - %s. Created at %s",
				k+1, card.Front, card.Back, card.CreatedAt.Format("15:04 01-02-2006")))
		}

		return c.Send(strings.Join(results, "\n"))
	}
}
