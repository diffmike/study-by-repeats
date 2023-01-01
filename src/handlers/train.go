package handlers

import (
	tele "gopkg.in/telebot.v3"
	"studyAndRepeat/src/database"
	"studyAndRepeat/src/services"
)

func Train(db *database.DB, showAnswer tele.Btn) tele.HandlerFunc {
	return func(c tele.Context) error {
		result, reply, err := services.Train(db, c.Sender().ID, showAnswer)
		if err != nil {
			return err
		}

		return c.Send(result, reply, tele.Silent)
	}
}

func ShowAnswer(db *database.DB, answers []tele.Btn) tele.HandlerFunc {
	return func(c tele.Context) error {
		result, reply, err := services.ShowAnswer(db, c.Sender().ID, c.Data(), answers)
		if err != nil {
			return err
		}

		return c.Send(result, reply, tele.Silent)
	}
}

func SaveAnswer(db *database.DB, repeatInHours int64, showAnswer tele.Btn) tele.HandlerFunc {
	return func(c tele.Context) error {
		result, reply, err := services.SaveAnswer(db, c.Sender().ID, c.Data(), repeatInHours, showAnswer)
		if err != nil {
			return err
		}

		return c.Send(result, reply, tele.Silent)
	}
}
