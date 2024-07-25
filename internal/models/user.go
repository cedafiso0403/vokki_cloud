package models

import (
	"database/sql"
	"log"
	"regexp"
	"vokki_cloud/internal/database"
	"vokki_cloud/internal/utils"
)

type User struct {
	ID        int    `json:"id" db:"id"`
	Email     string `json:"email" db:"email"`
	Password  string `json:"password" db:"password"`
	CreatedAt string `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt string `json:"updated_at,omitempty" db:"updated_at"`
}

type NewUserRequest struct {
	Email                string `json:"email"`
	Password             string `json:"password"`
	ConfirmationPassword string `json:"confirmation_password"`
}

func GetUser(email string) (User, error) {

	db := database.GetDB()

	var user = User{}

	preparedQuery, err := db.Prepare("SELECT id, created_at, email, updated_at FROM users WHERE email=$1")

	if err != nil {

		return user, err
	}

	defer preparedQuery.Close()

	err = preparedQuery.QueryRow(email).Scan(&user.ID, &user.CreatedAt, &user.Email, &user.UpdatedAt)

	if err != nil {

		return user, err
	}

	return user, nil
}

func (newUser *NewUserRequest) CreateUser() (User, error) {

	db := database.GetDB()

	user := User{}

	preparedCreateUserQuery, err := db.Prepare("INSERT INTO users (email, hashed_password) VALUES ($1, $2) RETURNING id, email")

	if err != nil {
		log.Print("Error preparing create user: ", err)
		return User{}, err
	}

	defer preparedCreateUserQuery.Close()

	hashedPassword, err := utils.HashPassword(newUser.Password)

	if err != nil {
		log.Print("Error hashing password: ", err)
		return User{}, err
	}

	err = preparedCreateUserQuery.QueryRow(&newUser.Email, hashedPassword).Scan(&user.ID, &user.Email)

	if err != nil && err != sql.ErrNoRows {
		log.Print("Error creating user: ", err)
		return User{}, err
	}

	return user, nil
}

func (user *NewUserRequest) IsValidEmail() bool {
	var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(user.Email)
}

//! Missing is password valid
