package services

import (
	tele "gopkg.in/telebot.v3"
	"strconv"
	"strings"
	"studyAndRepeat/src/database"
)

func Train(db *database.DB, tgId int64, showAnswer tele.Btn) (string, *tele.ReplyMarkup, error) {
	sessionId, err := db.GenerateSession(tgId)
	if err != nil {
		return "", nil, err
	}
	card, repeatId, err := db.FindRandomCardToRepeat(sessionId)
	if err != nil {
		return "", nil, err
	}
	if card.Id == 0 {
		return "Please add some cards to study. Use /add for this", nil, nil
	}

	reply := &tele.ReplyMarkup{ResizeKeyboard: true}
	showAnswer.Data = strconv.FormatInt(card.Id, 10) + "|" + strconv.FormatInt(repeatId, 10)
	reply.Inline(reply.Row(showAnswer))
	return card.Front, reply, nil
}

func ShowAnswer(db *database.DB, tgId int64, data string, answers []tele.Btn) (string, *tele.ReplyMarkup, error) {
	ids := strings.Split(data, "|")
	cardId, err := strconv.ParseInt(ids[0], 10, 64)
	if err != nil {
		return "", nil, err
	}
	card, err := db.FindCardById(cardId, tgId)
	if err != nil {
		return "", nil, err
	}

	reply := &tele.ReplyMarkup{}
	for i := range answers {
		answers[i].Data = ids[1]
	}

	reply.Inline(reply.Row(answers...))
	return card.Back.String, reply, nil
}

func SaveAnswer(db *database.DB, tgId int64, data string, repeatInHours int8, showAnswer tele.Btn) (string, *tele.ReplyMarkup, error) {
	repeatId, err := strconv.ParseInt(data, 10, 64)
	if err != nil {
		return "", nil, err
	}
	sessionId, err := db.UpdateRepeatIn(repeatId, tgId, repeatInHours)
	if err != nil {
		return "", nil, err
	}
	card, repeatId, err := db.FindRandomCardToRepeat(sessionId)
	if err != nil {
		return "", nil, err
	}
	if card.Id == 0 {
		return "You finished this session ✅\nNow you can add some more cards to study. Use /add for this", nil, nil
	}

	reply := &tele.ReplyMarkup{ResizeKeyboard: true}
	showAnswer.Data = strconv.FormatInt(card.Id, 10) + "|" + strconv.FormatInt(repeatId, 10)
	reply.Inline(reply.Row(showAnswer))
	return card.Front, reply, nil
}
