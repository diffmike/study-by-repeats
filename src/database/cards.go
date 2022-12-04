package database

import (
	"database/sql"
	"time"
)

type Card struct {
	Front     string
	Back      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func StoreCard(db *sql.DB, userId int64, front string) (int64, error) {
	var cardId int64
	err := db.QueryRow("INSERT INTO cards (user_id, front) VALUES ($1, $2) RETURNING id", userId, front).Scan(&cardId)

	return cardId, err
}

func DeleteCard(db *sql.DB, userId int64, front string) error {
	_, err := db.Exec("DELETE FROM cards WHERE user_id = $1 AND front = $2", userId, front)

	return err
}

func SetBackForCard(db *sql.DB, user_id int64, back string, cardId int64) error {
	row := db.QueryRow("UPDATE cards SET back = $1 WHERE user_id = $2 and id = $3", back, user_id, cardId)
	if row.Err() == sql.ErrNoRows {
		return nil
	}

	return row.Err()
}

func FindUserCards(db *sql.DB, user_id int64) (cards []Card, err error) {
	rows, err := db.Query("SELECT front, back, created_at, updated_at FROM cards WHERE user_id = $1 ORDER BY updated_at DESC", user_id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var card Card
		err = rows.Scan(&card.Front, &card.Back, &card.CreatedAt, &card.UpdatedAt)
		if err != nil {
			return nil, err
		}
		cards = append(cards, card)
	}
	return cards, rows.Err()
}

func FindCardByFront(db *sql.DB, user_id int64, front string) (int64, error) {
	var cardId int64
	err := db.QueryRow("SELECT id FROM cards WHERE user_id = $1 and front = $2", user_id, front).Scan(&cardId)
	if err == sql.ErrNoRows {
		return 0, nil
	}

	return cardId, err
}

func FindLatestUserCard(db *sql.DB, user_id int64) (int64, string, error) {
	var cardId int64
	var front string
	query := "SELECT id, front FROM cards WHERE user_id = $1 and back is NULL ORDER BY updated_at DESC LIMIT 1"
	err := db.QueryRow(query, user_id).Scan(&cardId, &front)
	if err == sql.ErrNoRows {
		return 0, "", nil
	}

	return cardId, front, err
}
