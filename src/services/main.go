package services

import (
	"studyAndRepeat/src/database"
)

func StoreUser(db *database.DB, tgId int64, username string) (string, error) {
	id, err := db.FindUserById(tgId)
	if err != nil {
		return "", err
	}
	if id == 0 {
		id, err = db.StoreUser(tgId, username)
	}

	return "Thank you! Now lets add some cards to study. Use /add for this", err
}
