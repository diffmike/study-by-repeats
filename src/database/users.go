package database

import (
	"database/sql"
)

func StoreUser(db *sql.DB, tgId int64, username string) (int64, error) {
	var userId int64
	err := db.QueryRow("INSERT INTO users (username, tg_id) VALUES ($1, $2) RETURNING id", username, tgId).Scan(&userId)

	return userId, err
}

func FindUserById(db *sql.DB, tgId int64) (int64, error) {
	var userId int64
	err := db.QueryRow("SELECT id FROM users WHERE tg_id = $1", tgId).Scan(&userId)
	if err == sql.ErrNoRows {
		return 0, nil
	}

	return userId, err
}
