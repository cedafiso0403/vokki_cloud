package models

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	aiagent "vokki_cloud/internal/AI_agent"
	"vokki_cloud/internal/database"

	"github.com/google/generative-ai-go/genai"
)

type Language struct {
	LanguageCode string `json:"language_code" db:"langauge_code"`
	LanguageName string `json:"language_name" db:"language_name"`
}

type Word struct {
	ID   int      `json:"id" db:"id"`
	Word string   `json:"word" db:"word"`
	Lang Language `json:"lang"`
}

type WordTranslations struct {
	InputWord    Word   `json:"input_word"`
	Translations []Word `json:"translations"`
}

func GetWordByText(word string) (Word, error) {

	var w Word

	err := database.GetPreparedGetWordByText().QueryRow(word).Scan(&w.ID, &w.Word, &w.Lang.LanguageCode, &w.Lang.LanguageName)

	return w, err

}

func GetAllLanguage() ([]Language, error) {

	var languages []Language

	rows, err := database.GetPreparedGetAllLanguages().Query()

	if err != nil {
		return languages, err
	}

	defer rows.Close()

	for rows.Next() {
		var l Language
		if err := rows.Scan(&l.LanguageCode, &l.LanguageName); err != nil {
			return languages, err
		}

		languages = append(languages, l)
	}

	return languages, nil

}

func GetWordTranslations(word string, lang string) (WordTranslations, error) {

	var wt WordTranslations

	word = strings.ToLower(strings.TrimSpace(strings.Trim(strings.Trim(word, `"`), `'`)))

	w, err := GetWordByText(word)

	log.Println("Word: ", w)
	log.Println("Query Word: ", word)

	// If the word is not found, insert it into the database
	if err != nil && err != sql.ErrNoRows {
		log.Println("Error getting word by text: ", err)
		return wt, err
	}

	if err == sql.ErrNoRows {

		// If the word is not found, return an empty WordTranslations struct
		// Here we want to implement gemini to translate the word and add it
		log.Println("Word not found: ", word)

		translations := generateTranslation(word)

		wt.InputWord = translations[0]
		wt.Translations = translations[1:]

		errInserting := InsertTranslations(&wt)

		if errInserting != nil {
			log.Println("Error inserting translations: ", errInserting)
			return wt, errInserting
		}

		// Insert the word into the database

		return wt, nil
	} else {

		wt.InputWord = w

		rows, err := database.GetPreparedGetTranslations().Query(w.ID)

		if err != nil {
			log.Println("Error getting translations: ", err)
			return wt, err
		}

		defer rows.Close()

		for rows.Next() {
			var w Word
			if err := rows.Scan(&w.ID, &w.Word, &w.Lang.LanguageCode, &w.Lang.LanguageName); err != nil {
				return wt, err
			}

			wt.Translations = append(wt.Translations, w)
		}

	}

	return wt, nil
}

// ! Create a method to retry linking if fails or deleting the words
func InsertTranslations(words *WordTranslations) error {

	wordsToAdd := []*Word{}

	wordsToAdd = append(wordsToAdd, &words.InputWord)

	for index := range words.Translations {
		wordsToAdd = append(wordsToAdd, &words.Translations[index])
	}

	for index, w := range wordsToAdd {
		err := database.GetPreparedInsertWord().QueryRow(w.Word, w.Lang.LanguageCode).Scan(&wordsToAdd[index].ID)
		if err != nil {
			log.Println("Error inserting word: ", err)
			return err
		}
	}

	var pairToInsert []string

	for i := 0; i < len(wordsToAdd); i++ {
		for j := i + 1; j < len(wordsToAdd); j++ {
			pairToInsert = append(pairToInsert, fmt.Sprintf("(%d, %d)", wordsToAdd[i].ID, wordsToAdd[j].ID))
		}
	}

	log.Println("Pairs to insert: ", pairToInsert)

	_, err := database.GetDB().Exec("INSERT INTO word_translations (word_id, translated_word_id) VALUES " + strings.Join(pairToInsert, ", "))

	if err != nil {
		log.Println("Error linking translations: ", err)
		return err
	}

	return nil

}

func generateTranslation(word string) []Word {

	var words []Word

	ls, err := GetAllLanguage()

	if err != nil {
		log.Println("Error getting all languages: ", err)
		return []Word{}
	}

	//! Might want to move this to AI package
	model := aiagent.GetAIAgent().GenerativeModel("gemini-1.5-flash")

	model.GenerationConfig = genai.GenerationConfig{
		ResponseMIMEType: "application/json",
	}

	prompt := fmt.Sprintf(`Translate the word %s using this JSON schema:
	[{ "type": "object",
		"properties": {
			{ 
				"word": "string",
				"lang": 
					"type": "object",
						properties": {
							"language_code": "string",
							"language_name": "string"
						},
			}
		}
	}]
	Where lang follows the ISO 639-1 standard and langName is the name of the language.	
	`, word)

	prompt += "to for all the following languages: "

	for _, l := range ls {
		prompt += fmt.Sprintf("%s: %s, ", l.LanguageCode, l.LanguageName)
	}

	prompt += "placing the object containing the input word at the first position of the array."

	ctx := context.Background()
	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		log.Fatal(err)
	}

	var contentText string

	for _, c := range resp.Candidates {
		if c.Content != nil {
			for _, part := range c.Content.Parts {
				switch p := part.(type) {
				case genai.Text:
					// Handle Text parts
					contentText += string(p)
				default:
					// Handle unknown parts
					contentText += ""
				}
			}
			words = append(words, Word{Word: contentText})
		}
	}

	err = json.Unmarshal([]byte(contentText), &words)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return []Word{}
	}

	return words
}
