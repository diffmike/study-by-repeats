package services

import (
	"studyAndRepeat/src/database"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestStoreNewUser(t *testing.T) {
	mockDb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDb.Close()

	columns := []string{"id"}
	mock.ExpectQuery("SELECT id FROM users WHERE tg_id = (.+)").
		WithArgs(1).
		WillReturnRows(mock.NewRows(columns))
	mock.ExpectQuery("INSERT INTO users").
		WithArgs("username", 1).
		WillReturnRows(mock.NewRows(columns).AddRow(1))

	db := &database.DB{DB: mockDb}
	result, err := StoreUser(db, 1, "username")
	if err != nil {
		t.Errorf("error was not expected while updating stats: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	if result != "Thank you! Now lets add some cards to study. Use /add for this" {
		t.Errorf("results is wrong: %s", result)
	}
}

func TestStoreExistedUser(t *testing.T) {
	mockDb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDb.Close()

	columns := []string{"id"}
	mock.ExpectQuery("SELECT id FROM users WHERE tg_id = (.+)").
		WithArgs(1).
		WillReturnRows(mock.NewRows(columns).AddRow(1))

	db := &database.DB{DB: mockDb}
	if _, err := StoreUser(db, 1, "username"); err != nil {
		t.Errorf("error was not expected while updating stats: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
