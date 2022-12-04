package database

import (
	"database/sql"
)

func StoreUser(db *sql.DB, uuid int64, username string) (int64, error) {
	var userId int64
	err := db.QueryRow("INSERT INTO users (username, uuid) VALUES ($1, $2) RETURNING id", username, uuid).Scan(&userId)

	return userId, err
}

func FindUserId(db *sql.DB, uuid int64) (int64, error) {
	var userId int64
	err := db.QueryRow("SELECT id FROM users WHERE uuid = $1", uuid).Scan(&userId)
	if err == sql.ErrNoRows {
		return 0, nil
	}

	return userId, err
}
