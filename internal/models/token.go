package models

import (
	"database/sql"
	vokki_constants "vokki_cloud/internal/constants"
	"vokki_cloud/internal/database"

	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	UserID int64 `json:"user_id"`
	jwt.StandardClaims
}

type AuthToken struct {
	UserID    int64                     `json:"user_id"`
	Token     string                    `json:"token"`
	RevokedAt sql.NullTime              `json:"revoked_at"`
	TokenType vokki_constants.TokenType `json:"token_type"`
}

func StoreToken(userID int64, token string, tokenType vokki_constants.TokenType) error {

	db := database.GetDB()

	_, err := db.Exec("INSERT INTO user_tokens (user_id, verification_token, token_type) VALUES ($1, $2, $3)", userID, token, tokenType)
	return err
}
