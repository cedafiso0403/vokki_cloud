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

	prepareCurrentTokenQuery, err := db.Prepare("UPDATE user_tokens set revoked_at = now() WHERE user_id=$1 AND token_type=$2 AND revoked_at IS NULL RETURNING verification_token")

	if err != nil {
		log.Println("Error preparing current token query: ", err)
		return errors.New("")
	}

	defer prepareCurrentTokenQuery.Close()

	oldtoken := AuthToken{}

	err = prepareCurrentTokenQuery.QueryRow(userID, tokenType).Scan(&oldtoken.Token)

	if err != nil && err != sql.ErrNoRows {
		log.Println("Error querying current token: ", err)
		return errors.New("error querying current token")
	}

	if oldtoken.Token == "" || err == sql.ErrNoRows {
		log.Println("No token to revoke")
	}

	log.Println("Token revoked")

	shared.GetTokenManager().RemoveToken(oldtoken.Token)

	_, err = db.Exec("INSERT INTO user_tokens (user_id, verification_token, token_type) VALUES ($1, $2, $3)", userID, token, tokenType)

	if err != nil {
		log.Println("Error inserting token: ", err)
		return errors.New("")
	}

	log.Println("Token inserted")

	shared.GetTokenManager().AddToken(token)

	return err
}

func VerifyToken(token string) bool {

	db := database.GetDB()

	var count int

	preparedTokenExistsQuery, err := db.Prepare("SELECT COUNT(*) FROM user_tokens WHERE verification_token=$1 AND revoked_at IS NULL")

	if err != nil {
		log.Println("Error preparing query for verifying token: ", err)
		return false
	}

	err = preparedTokenExistsQuery.QueryRow(token).Scan(&count)

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

	db := database.GetDB()

	prepareCurrentTokenQuery, err := db.Prepare("UPDATE user_tokens set revoked_at = now() WHERE verification_token=$1 AND revoked_at IS NULL")

	if err != nil {
		log.Println("Error preparing query to revoke token: ", err)
		return errors.New("")
	}

	defer prepareCurrentTokenQuery.Close()

	_, err = prepareCurrentTokenQuery.Exec(token)

	if err != nil {
		log.Println("Error revoking token: ", err)
	}

	shared.GetTokenManager().RemoveToken(token)

	return err
}
