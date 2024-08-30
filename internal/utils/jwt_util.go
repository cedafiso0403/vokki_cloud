package utils

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"
	"vokki_cloud/internal/auth_error"
	vokki_constants "vokki_cloud/internal/constants"

	"github.com/dgrijalva/jwt-go"
)

// !Create a key on env
var jwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))

type Claims struct {
	UserID int `json:"user_id"`
	jwt.StandardClaims
}

func GenerateJWT(userID int) (string, error) {

	nonce := time.Now().UTC()

	claims := &Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: nonce.Add(time.Hour * 24).Unix(), // Token expires in 24 hours
			IssuedAt:  nonce.Unix(),
			Issuer:    vokki_constants.TokenIssuer,
			Id:        fmt.Sprintf("%d", nonce.Unix()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func ParseJWT(tokenString string) (*Claims, error) {
	claims := Claims{}

	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {

		if []byte(jwtKey) == nil {
			log.Println("JWT_SECRET_KEY environment variable not set")
			return nil, errors.New("")
		}
		return []byte(jwtKey), nil
	})

	if err != nil {
		// Use type assertion to check for specific JWT errors
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				// Token is malformed
				return nil, auth_error.ErrInvalidToken
			} else if ve.Errors&jwt.ValidationErrorSignatureInvalid != 0 {
				// Signature is invalid
				return nil, auth_error.ErrInvalidToken

			} else if ve.Errors&jwt.ValidationErrorIssuer != 0 {
				// Issuer is invalid
				return nil, auth_error.ErrInvalidTokenIssuer
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, auth_error.ErrExpiredToken
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				// Token is not yet valid
				return nil, auth_error.ErrInvalidToken
				// Other validation errors
			} else {
				return nil, auth_error.ErrInvalidToken
			}
		}
		// Return other errors that are not validation errors
		return nil, auth_error.ErrInvalidToken
	}

	// Check if the token is valid
	if !token.Valid {
		return nil, auth_error.ErrInvalidToken
	}

	return &claims, err
}

func ValidateToken(tokenString string) (*Claims, error) {

	decodedToken, err := ParseJWT(tokenString)

	if err != nil {
		return nil, err
	}

	return decodedToken, nil

}
