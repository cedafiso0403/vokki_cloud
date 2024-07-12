package models

import (
	"vokki_cloud/internal/database"

	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	UserID int64 `json:"user_id"`
	jwt.StandardClaims
}

func StoreToken(userID int64, token string) error {

	db := database.GetDB()

	_, err := db.Exec("INSERT INTO user_auth (user_id, verification_token) VALUES ($1, $2)", userID, token)
	return err
}
