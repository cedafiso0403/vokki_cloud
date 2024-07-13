package utils

import (
	"errors"
	"time"
	vokki_constants "vokki_cloud/internal/constants"
	"vokki_cloud/internal/models"

	"github.com/dgrijalva/jwt-go"
)

// !Create a key on env
var jwtKey = []byte("your-secret-key")

func GenerateJWT(userID int64) (string, error) {
	claims := &models.Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
			IssuedAt:  time.Now().Unix(),
			Issuer:    vokki_constants.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func ParseJWT(tokenString string) (*models.Claims, error) {
	claims := &models.Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, errors.New("invalid token signature")
		}
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// func ValidateToken(tokenString string) (*Claims, error) {
// 	// Validate JWT token

// }
