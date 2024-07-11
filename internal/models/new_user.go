package models

import (
	"regexp"
	"vokki_cloud/internal/database"

	"golang.org/x/crypto/bcrypt"
)

type NewUser struct {
	Email                string `json:"email"`
	Password             string `json:"password"`
	ConfirmationPassword string `json:"confirmation_password"`
}

func (user *NewUser) CreateUser() error {

	db := database.GetDB()

	preparedCreateUserQuery, err := db.Prepare("INSERT INTO users (email, password_hash) VALUES ($1, $2)")

	if err != nil {
		return err
	}

	hashedPassword, err := hashPassword(user.Password)

	if err != nil {
		return err
	}

	defer preparedCreateUserQuery.Close()

	_, err = preparedCreateUserQuery.Exec(&user.Email, hashedPassword)

	if err != nil {
		return err
	}

	return nil
}

func (user *NewUser) IsValidEmail() bool {
	var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(user.Email)
}

//! Missing is password valid

func hashPassword(password string) (string, error) {
	// bcripts hash and salt password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}
