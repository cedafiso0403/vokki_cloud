package services

import (
	"database/sql"
	"vokki_cloud/internal/auth_error"
	"vokki_cloud/internal/database"
	"vokki_cloud/internal/models"
	"vokki_cloud/internal/utils"
)

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Authenticate(credentials Credentials) (int, string, error) {
	db := database.GetDB()
	var user models.User
	err := db.QueryRow("SELECT id, email, password_hash FROM users WHERE email = $1", credentials.Email).Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, "", auth_error.ErrUserNotFound
		}
		return 0, "", err
	}

	if !utils.CheckPasswordHash(credentials.Password, user.Password) {
		return 0, "", auth_error.ErrIncorrectCredentials
	}

	token, err := utils.GenerateJWT(int64(user.ID))
	if err != nil {
		return 0, "", err
	}

	return user.ID, token, nil
}
