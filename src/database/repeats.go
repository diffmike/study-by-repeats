package database

import "database/sql"

func GenerateSession(db *sql.DB, tgId int64) (sessionId int64, err error) {
	err = db.QueryRow("INSERT INTO sessions (tg_id) VALUES ($1) RETURNING id", tgId).Scan(&sessionId)
	if err != nil {
		return sessionId, err
	}
	cardIds, err := findCardsForNewSession(db, tgId)
	if err != nil {
		return sessionId, err
	}
	for _, id := range cardIds {
		row := db.QueryRow("INSERT INTO repeats (card_id, session_id) VALUES ($1, $2) RETURNING id", id, sessionId)
		if row.Err() != nil {
			return sessionId, err
		}
	}

	return sessionId, nil
}

func findCardsForNewSession(db *sql.DB, tgId int64) (cardIds []int64, err error) {
	rows, err := db.Query("SELECT id FROM cards WHERE tg_id = $1 AND back IS NOT NULL AND (repeat_after <= NOW() or repeat_after IS NULL) LIMIT 20", tgId)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var cardId int64
		err = rows.Scan(&cardId)
		if err != nil {
			return nil, err
		}
		cardIds = append(cardIds, cardId)
	}
	return cardIds, rows.Err()
}

func FindRandomCardToRepeat(db *sql.DB, sessionId int64) (card Card, repeatId int64, err error) {
	err = db.QueryRow("SELECT c.id, c.front, c.back, r.id FROM cards AS c "+
		"JOIN repeats AS r ON r.card_id = c.id "+
		"WHERE (r.repeat_in IS NULL or r.repeat_in = 0) AND r.session_id = $1"+
		"ORDER BY random() LIMIT 1", sessionId).Scan(&card.Id, &card.Front, &card.Back, &repeatId)
	if err == sql.ErrNoRows {
		return card, repeatId, nil
	}

	return card, repeatId, err
}
