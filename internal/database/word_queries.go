package database

import (
	"database/sql"
	"log"
)

var (
	preparedGetWordByText   *sql.Stmt
	preparedGetTranslations *sql.Stmt
	preparedGetAllLanguages *sql.Stmt
	preparedInsertWord      *sql.Stmt
)

func initPreparedWordStatements() error {

	var err error

	log.Println("Preparing word statements")

	if preparedGetWordByText, err = db.Prepare("SELECT words.id, word, languages.language_code, language_name FROM words INNER JOIN languages on words.language_id = languages.id WHERE words.word=$1"); err != nil {
		log.Println("Error preparing get word by text query: ", err)
		return err
	}

	if preparedGetTranslations, err = db.Prepare("SELECT DISTINCT words.id, words.word, languages.language_code, languages.language_name FROM words INNER JOIN word_translations ON (words.id = word_translations.word_id AND word_translations.translated_word_id = $1) OR (words.id = word_translations.translated_word_id AND word_translations.word_id = $1) INNER JOIN languages ON languages.id = words.language_id"); err != nil {
		log.Println("Error preparing get translations query: ", err)
		return err
	}

	if preparedGetAllLanguages, err = db.Prepare("SELECT language_code, language_name FROM languages"); err != nil {
		log.Println("Error preparing get all languages query: ", err)
		return err
	}

	if preparedInsertWord, err = db.Prepare("WITH lang AS (SELECT id FROM languages WHERE language_code = $2) INSERT INTO words (word, language_id) VALUES ($1, (SELECT id FROM lang)) RETURNING id"); err != nil {
		log.Println("Error preparing insert word query: ", err)
		return err
	}

	log.Println("All word statements prepared successfully")

	return nil
}

func closePreparedWordStatements() {

	if preparedGetWordByText != nil {
		preparedGetWordByText.Close()
	}

	if preparedGetTranslations != nil {
		preparedGetTranslations.Close()
	}

	if preparedGetAllLanguages != nil {
		preparedGetAllLanguages.Close()
	}

	if preparedInsertWord != nil {
		preparedInsertWord.Close()
	}
}

func GetPreparedGetWordByText() *sql.Stmt {
	return preparedGetWordByText
}

func GetPreparedGetTranslations() *sql.Stmt {
	return preparedGetTranslations
}

func GetPreparedGetAllLanguages() *sql.Stmt {
	return preparedGetAllLanguages
}

func GetPreparedInsertWord() *sql.Stmt {
	return preparedInsertWord
}
