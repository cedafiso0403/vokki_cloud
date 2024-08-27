package models

import (
	"database/sql"
	"vokki_cloud/internal/database"
)

type Word struct {
	ID       int    `json:"id" db:"id"`
	Word     string `json:"word" db:"word"`
	Lang     string `json:"lang" db:"language_code"`
	LangName string `json:"lang_name" db:"language_name"`
}

type WordTranslations struct {
	InputWord    Word   `json:"input_word"`
	Lang         string `json:"lang"`
	Translations []Word `json:"translations"`
}

func GetWordByText(word string) (Word, error) {

	var w Word

	err := database.GetPreparedGetWordByText().QueryRow(word).Scan(&w.ID, &w.Word, &w.Lang, &w.LangName)

	if err != nil && err != sql.ErrNoRows {
		return w, err
	}

	return w, nil

}

func GetWordTranslations(word string, lang string) (WordTranslations, error) {

	var wt WordTranslations

	w, err := GetWordByText(word)

	if err != nil {
		return wt, err
	}

	wt.InputWord = w
	wt.Lang = lang

	rows, err := database.GetPreparedGetTranslations().Query(w.ID, lang)

	if err != nil {
		return wt, err
	}

	defer rows.Close()

	for rows.Next() {
		var w Word
		if err := rows.Scan(&w.ID, &w.Word, &w.Lang, &w.LangName); err != nil {
			return wt, err
		}

		wt.Translations = append(wt.Translations, w)
	}

	return wt, nil
}
