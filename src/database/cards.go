package database

import (
	"database/sql"
	"time"
)

type Card struct {
	Id          int64
	Front       string
	Back        sql.NullString
	CreatedAt   time.Time
	RepeatAfter sql.NullTime
}

func StoreCard(db *sql.DB, tgId int64, front string) (int64, error) {
	var cardId int64
	err := db.QueryRow("INSERT INTO cards (tg_id, front) VALUES ($1, $2) RETURNING id", tgId, front).Scan(&cardId)

	return cardId, err
}

func DeleteCard(db *sql.DB, tgId int64, front string) error {
	_, err := db.Exec("DELETE FROM cards WHERE tg_id = $1 AND front = $2", tgId, front)

	return err
}

func SetBackForCard(db *sql.DB, tgId int64, back string, cardId int64) error {
	row := db.QueryRow("UPDATE cards SET back = $1 WHERE tg_id = $2 and id = $3", back, tgId, cardId)
	if row.Err() == sql.ErrNoRows {
		return nil
	}

	return row.Err()
}

func FindUserCards(db *sql.DB, tgId int64) (cards []Card, err error) {
	rows, err := db.Query("SELECT front, back, created_at, repeat_after FROM cards WHERE tg_id = $1 ORDER BY created_at DESC", tgId)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var card Card
		err = rows.Scan(&card.Front, &card.Back, &card.CreatedAt, &card.RepeatAfter)
		if err != nil {
			return nil, err
		}
		cards = append(cards, card)
	}
	return cards, rows.Err()
}

func FindCardByFront(db *sql.DB, tgId int64, front string) (int64, error) {
	var cardId int64
	err := db.QueryRow("SELECT id FROM cards WHERE tg_id = $1 and front = $2", tgId, front).Scan(&cardId)
	if err == sql.ErrNoRows {
		return 0, nil
	}

	return cardId, err
}

func FindLatestUserCard(db *sql.DB, tgId int64) (int64, string, error) {
	var cardId int64
	var front string
	query := "SELECT id, front FROM cards WHERE tg_id = $1 and back is NULL ORDER BY created_at DESC LIMIT 1"
	err := db.QueryRow(query, tgId).Scan(&cardId, &front)
	if err == sql.ErrNoRows {
		return 0, "", nil
	}

	return cardId, front, err
}

func FindCardById(db *sql.DB, cardId int64, tgId int64) (card Card, err error) {
	query := "SELECT id, front, back, created_at FROM cards WHERE id = $1 and tg_id = $2 LIMIT 1"
	err = db.QueryRow(query, cardId, tgId).Scan(&card.Id, &card.Front, &card.Back, &card.CreatedAt)
	if err == sql.ErrNoRows {
		return card, nil
	}

	return card, err
}

func UpdateRepeatIn(db *sql.DB, repeatId int64, tgId int64, repeatIn int8) (sessionId int64, err error) {
	var cardId int64
	err = db.QueryRow("UPDATE repeats SET repeat_in = $1 WHERE id = $2 RETURNING card_id, session_id", repeatIn, repeatId).Scan(&cardId, &sessionId)
	if err == sql.ErrNoRows {
		return 0, nil
	}

	repeatAfter := time.Now().Add(time.Hour * time.Duration(repeatIn))
	row := db.QueryRow("UPDATE cards SET repeat_after = $1 WHERE tg_id = $2 and id = $3", repeatAfter, tgId, cardId)
	if row.Err() == sql.ErrNoRows {
		return 0, nil
	}

	return sessionId, row.Err()
}
