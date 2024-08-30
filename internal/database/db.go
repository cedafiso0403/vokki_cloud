package database

import (
	"database/sql"
	"log"
	"os"
)

var (
	db *sql.DB
)

func Connect() {

	var err error
	db, err = sql.Open("postgres", os.Getenv("DB_URL"))

	if err != nil {
		log.Fatal("Cannot connect to the database: ", err)
	}

	_, err = db.Exec("DEALLOCATE ALL")

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to the database")

	err = initPreparedStatements()
	if err != nil {
		log.Fatal("Error initializing prepared statements: ", err)
	}

}

func GetDB() *sql.DB {
	return db
}

func initPreparedStatements() error {
	var err error

	log.Println("Preparing statements...")

	if err = initPreparedTokenStatements(); err != nil {
		log.Println("Error initializing token prepared statements: ", err)
		return err
	}

	if err = initPreparedUserStaments(); err != nil {
		log.Println("Error initializing user prepared statements: ", err)
		return err
	}

	if err = initPreparedWordStatements(); err != nil {
		log.Println("Error initializing word prepared statements: ", err)
		return err
	}

	log.Println("All statements prepared successfully")
	return nil
}

func Close() {
	log.Println("Closing database connections and prepared statements...")

	closePreparedTokenStatements()

	closePreparedUserStatements()

	closePreparedWordStatements()

	if db != nil {
		db.Close()
	}
}
