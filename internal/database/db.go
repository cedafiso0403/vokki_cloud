package database

import (
	"database/sql"
	"log"
	"os"
)

var db *sql.DB

func Connect() (*sql.DB, error) {
	var err error
	db, err = sql.Open("postgres", os.Getenv("DB_URL"))

	if err != nil {
		log.Fatal("Cannot connect to the database: ", err)
	}

	log.Println("Connected to the database")

	return db, nil
}

func GetDB() *sql.DB {
	return db
}
