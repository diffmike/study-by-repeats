package handlers

import (
	"database/sql"
	tele "gopkg.in/telebot.v3"
	"studyAndRepeat/src/database"
)

func Start(db *sql.DB) tele.HandlerFunc {
	return func(c tele.Context) error {
		id, err := database.FindUserId(db, c.Sender().ID)
		if err != nil {
			return err
		}
		if id == 0 {
			id, err = database.StoreUser(db, c.Sender().ID, c.Sender().Username)
		}
		if err != nil {
			return err
		}

		return c.Send("Thank you ! Now lets add some cards to study. Use /add for this")
	}
}
