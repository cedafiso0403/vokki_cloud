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

func GetUser(email string) (User, error) {

	db := database.GetDB()

	var user = User{}

	preparedQuery, err := db.Prepare("SELECT id, created_at, email, updated_at FROM users WHERE email=$1")

	if err != nil {

		return user, err
	}

	defer preparedQuery.Close()

	err = preparedQuery.QueryRow(email).Scan(&user.ID, &user.Created, &user.Email, &user.Updated)

	if err != nil {

		return user, err
	}

	return user, nil
}
