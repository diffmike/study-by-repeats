package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"os"
)

func New() *sql.DB {
	dbPort := os.Getenv("DB_HOST")
	if dbPort == "" {
		dbPort = "5432"
	}
	psqlInfo := fmt.Sprintf("host=%s user=%s port=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_USER"), dbPort, os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	return db
}
