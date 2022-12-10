package services

import (
	"database/sql"
	"fmt"
	"strings"
	"studyAndRepeat/src/database"
	"time"
)

func AddCard(db *database.DB, tgId int64, text string) (string, error) {
	var back sql.NullString
	if strings.Contains(text, ";") {
		pieces := strings.Split(text, ";")
		_ = back.Scan(strings.Trim(pieces[1], " "))
		text = strings.Trim(pieces[0], " ")
	}
	id, err := db.FindCardByFront(tgId, text)
	if err != nil {
		return "", err
	}
	if id == 0 {
		id, err = db.StoreCard(tgId, text, back)
	}
	if err != nil {
		return "", err
	}

	if back.Valid {
		return fmt.Sprintf("✅ Card: *%s - %s* was completed", text, back.String), nil
	}

	return fmt.Sprintf("*%s* was added!\nNow write the definition", text), nil
}

func DeleteCard(db *database.DB, tgId int64, frontText string) (string, error) {
	id, err := db.FindCardByFront(tgId, frontText)
	if err != nil {
		return "", err
	}
	if id == 0 {
		return "Hmm 🤨... Such cards wasn't found in your dictionary", nil
	}

	err = db.DeleteCard(tgId, frontText)

	return fmt.Sprintf("*%s* was deleted", frontText), err
}

func SetDefinition(db *database.DB, tgId int64, backText string) (string, error) {
	id, frontText, err := db.FindLatestUserCard(tgId)
	if err != nil {
		return "", err
	}
	if id == 0 {
		return "Hmm 🤨... It seems you need to use /add <phrase> beforehand", nil
	}

	err = db.SetBackForCard(tgId, backText, id)

	return fmt.Sprintf("✅ Card: *%s - %s* was completed", frontText, backText), err
}

func GetDictionary(db *database.DB, tgId int64) (string, error) {
	cards, err := db.FindUserCards(tgId)
	if err != nil {
		return "", err
	}

	results := []string{}
	for k, card := range cards {
		results = append(results, fmt.Sprintf("%d. *%s - %s*. %s",
			k+1, card.Front, card.Back.String, readableAfter(card.RepeatAfter)))
	}

	return strings.Join(results, "\n"), nil
}

func readableAfter(after sql.NullTime) string {
	if !after.Valid || after.Time.Before(time.Now()) {
		return "Available"
	}

	duration := after.Time.Sub(time.Now()).Round(time.Minute)
	hours := duration / time.Hour
	duration -= hours * time.Hour
	minutes := duration / time.Minute

	return fmt.Sprintf("In %02dh:%02dm", hours, minutes)
}
