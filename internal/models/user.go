package models

import (
	"vokki_cloud/internal/database"
)

type User struct {
	ID      int    `json:"id"`
	Email   string `json:"email"`
	Created string `json:"created,omitempty"`
	Updated string `json:"updated,omitempty"`
}

type NewUser struct {
	Email                string `json:"email"`
	Password             string `json:"password"`
	ConfirmationPassword string `json:"confirmation_password"`
}

func GetUser(email string) (User, error) {

	db := database.GetDB()

	var user = User{}

	row := db.QueryRow("SELECT id, created_at, email, updated_at FROM users WHERE email=$1", email)

	err := row.Scan(&user.ID, &user.Created, &user.Email, &user.Updated)

	if err != nil {
		return user, err
	}

	return user, nil
}
