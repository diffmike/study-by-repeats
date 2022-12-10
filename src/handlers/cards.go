package handlers

import (
	tele "gopkg.in/telebot.v3"
	"studyAndRepeat/src/database"
	"studyAndRepeat/src/services"
)

func AddCard(db *database.DB) tele.HandlerFunc {
	return func(c tele.Context) error {
		result, err := services.AddCard(db, c.Sender().ID, c.Message().Payload)
		if err != nil {
			return err
		}

		return c.Send(result, tele.ModeMarkdown)
	}
}

func DeleteCard(db *database.DB) tele.HandlerFunc {
	return func(c tele.Context) error {
		result, err := services.DeleteCard(db, c.Sender().ID, c.Message().Payload)
		if err != nil {
			return err
		}

		return c.Send(result, tele.ModeMarkdown)
	}
}

func SetDefinition(db *database.DB) tele.HandlerFunc {
	return func(c tele.Context) error {
		result, err := services.SetDefinition(db, c.Sender().ID, c.Message().Text)
		if err != nil {
			return err
		}

		return c.Send(result, tele.ModeMarkdown)
	}
}

func GetDictionary(db *database.DB) tele.HandlerFunc {
	return func(c tele.Context) error {
		result, err := services.GetDictionary(db, c.Sender().ID)
		if err != nil {
			return err
		}

		return c.Send(result, tele.ModeMarkdown)
	}
}
