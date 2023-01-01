package services

import (
	tele "gopkg.in/telebot.v3"
	"log"
	"os"
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
	var rows []tele.Row
	for i := range answers {
		answers[i].Data = ids[1]
		rows = append(rows, reply.Row(answers[i]))
	}

	reply.Inline(rows...)
	return card.Back.String, reply, nil
}

func SaveAnswer(db *database.DB, tgId int64, data string, repeatInHours int64, showAnswer tele.Btn) (string, *tele.ReplyMarkup, error) {
	repeatId, err := strconv.ParseInt(data, 10, 64)
	if err != nil {
		return "", nil, err
	}
	currentRepeatIn, err := db.FindPreviousRepeatInHoursById(repeatId)
	if err != nil {
		return "", nil, err
	}
	log.Printf("currentRepeatIn %d, repeatInHours: %d", currentRepeatIn.Int64, repeatInHours)
	if currentRepeatIn.Valid && repeatInHours < 0 {
		repeatInHours, err = defineRepeatInHours(currentRepeatIn.Int64)
		if err != nil {
			return "", nil, err
		}
	}
	if repeatInHours < 0 {
		repeatsConfig, err := getRepeatsConfiguration()
		if err != nil {
			return "", nil, err
		}
		repeatInHours = repeatsConfig[1]
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

func defineRepeatInHours(repeatInHours int64) (int64, error) {
	parameters, err := getRepeatsConfiguration()
	if err != nil {
		return 0, err
	}

	for k, parameter := range parameters {
		log.Printf("repeatInHours %d, parameterInt: %d, k: %d", repeatInHours, parameter, k)
		if repeatInHours == parameter && len(parameters) > k+1 {
			result := parameters[k+1]
			return result, nil
		}
	}

	return repeatInHours, nil
}

func getRepeatsConfiguration() (results []int64, err error) {
	configuration := "24,72,120,168,336,744"
	if os.Getenv("REPEATS_CONFIGURATION") != "" {
		configuration = os.Getenv("REPEATS_CONFIGURATION")
	}

	for _, parameter := range strings.Split(configuration, ",") {
		result, err := strconv.ParseInt(parameter, 10, 64)
		if err != nil {
			return results, err
		}

		results = append(results, result)
	}

	return results, nil
}
