package models

import (
	"database/sql"
	"errors"
	"log"
	vokki_constants "vokki_cloud/internal/constants"
	"vokki_cloud/internal/database"
	"vokki_cloud/internal/shared"
)

type AuthToken struct {
	UserID    int64                     `json:"user_id"`
	Token     string                    `json:"token"`
	RevokedAt sql.NullTime              `json:"revoked_at"`
	TokenType vokki_constants.TokenType `json:"token_type"`
}

func StoreToken(userID int64, token string, tokenType vokki_constants.TokenType) error {

	db := database.GetDB()

	oldtoken := AuthToken{}

	err := database.GetPreparedCurrentTokenQuery().QueryRow(userID, tokenType).Scan(&oldtoken.Token)

	if err != nil && err != sql.ErrNoRows {
		log.Println("Error querying current token: ", err)
		return errors.New("error querying current token")
	}

	if oldtoken.Token == "" || err == sql.ErrNoRows {
		log.Println("No token to revoke")
	}

	log.Println("Token revoked")

	shared.GetTokenManager().RemoveToken(oldtoken.Token)

	//! Make it a prepared statement
	_, err = db.Exec("INSERT INTO user_tokens (user_id, verification_token, token_type) VALUES ($1, $2, $3)", userID, token, tokenType)

	if err != nil {
		log.Println("Error inserting token: ", err)
		return errors.New("")
	}

	log.Printf("Token inserted %s \n", tokenType)

	shared.GetTokenManager().AddToken(token)

	return err
}

func VerifyToken(token string) bool {

	var count int

	err := database.GetPreparedTokenExistsQuery().QueryRow(token).Scan(&count)

	if err != nil {
		log.Println("Error verifying token: ", err)
		return false
	}

	if count == 0 {
		return false
	}

	return true
}

func RevokeToken(token string) error {

	_, err := database.GetPreparedCurrentTokenQuery().Exec(token)

	if err != nil {
		log.Println("Error revoking token: ", err)
	}

	shared.GetTokenManager().RemoveToken(token)

	return err
}
