package handlers

import (
	"database/sql"
	tele "gopkg.in/telebot.v3"
	"log"
	"strconv"
	"strings"
	"studyAndRepeat/src/database"
)

func Train(db *sql.DB, showAnswer tele.Btn) tele.HandlerFunc {
	return func(c tele.Context) error {
		sessionId, err := database.GenerateSession(db, c.Sender().ID)
		if err != nil {
			return err
		}
		card, repeatId, err := database.FindRandomCardToRepeat(db, sessionId)
		if err != nil {
			return err
		}
		if card.Id == 0 {
			return c.Send("Please add some cards to study. Use /add for this")
		}

		reply := &tele.ReplyMarkup{ResizeKeyboard: true}
		showAnswer.Data = strconv.FormatInt(card.Id, 10) + "|" + strconv.FormatInt(repeatId, 10)
		reply.Inline(reply.Row(showAnswer))
		return c.Send(card.Front, reply)
	}
}

func ShowAnswer(db *sql.DB, answers []tele.Btn) tele.HandlerFunc {
	return func(c tele.Context) error {
		ids := strings.Split(c.Data(), "|")
		log.Printf("cardId: %s - repeatId: %s", ids[0], ids[1])
		cardId, err := strconv.ParseInt(ids[0], 10, 64)
		if err != nil {
			return err
		}
		card, err := database.FindCardById(db, cardId, c.Sender().ID)
		if err != nil {
			return err
		}

		reply := &tele.ReplyMarkup{ResizeKeyboard: true}
		for i := range answers {
			answers[i].Data = ids[1]
		}

		reply.Inline(reply.Row(answers...))
		return c.Send(card.Back.String, reply)
	}
}

func SaveAnswer(db *sql.DB, repeatInHours int8, showAnswer tele.Btn) tele.HandlerFunc {
	return func(c tele.Context) error {
		repeatId, err := strconv.ParseInt(c.Data(), 10, 64)
		if err != nil {
			return err
		}
		sessionId, err := database.UpdateRepeatIn(db, repeatId, c.Sender().ID, repeatInHours)
		if err != nil {
			return err
		}
		card, repeatId, err := database.FindRandomCardToRepeat(db, sessionId)
		if err != nil {
			return err
		}
		if card.Id == 0 {
			return c.Send("You finished this session ✅\nNow you can add some more cards to study. Use /add for this")
		}

		reply := &tele.ReplyMarkup{ResizeKeyboard: true}
		showAnswer.Data = strconv.FormatInt(card.Id, 10) + "|" + strconv.FormatInt(repeatId, 10)
		reply.Inline(reply.Row(showAnswer))
		return c.Send(card.Front, reply)
	}
}
