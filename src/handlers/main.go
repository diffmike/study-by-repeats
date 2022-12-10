package handlers

import (
	tele "gopkg.in/telebot.v3"
	"studyAndRepeat/src/database"
	"studyAndRepeat/src/services"
)

func Start(db *database.DB) tele.HandlerFunc {
	return func(c tele.Context) error {
		result, err := services.StoreUser(db, c.Sender().ID, c.Sender().Username)
		if err != nil {
			return err
		}

		return c.Send(result)
	}
}
