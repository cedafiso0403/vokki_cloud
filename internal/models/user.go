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
	Email                string `json:"email" binding:"required" example:"user@domain.com"`
	Password             string `json:"password" binding:"required" example:"password"`
	ConfirmationPassword string `json:"confirmation_password" binding:"required" example:"password"`
}

type UserProfile struct {
	ID        int    `json:"id" db:"id" example:"1"`
	Email     string `json:"email" db:"email" example:"user@domain.com"`
	FirstName string `json:"first_name" db:"first_name" example:"John"`
	LastName  string `json:"last_name" db:"last_name" example:"Doe"`
}

type UpdateUserProfileRequest struct {
	FirstName *string `json:"first_name" db:"first_name" example:"John"`
	LastName  *string `json:"last_name" db:"last_name" example:"Doe"`
}

func GetUser(email string) (User, error) {

	user := User{}

	err := database.GetPreparedGetUserQuery().QueryRow(email).Scan(&user.ID, &user.CreatedAt, &user.Email, &user.UpdatedAt)

	if err != nil {

		return user, err
	}

	return user, nil
}

func (newUser *NewUserRequest) CreateUser() (User, error) {

	user := User{}

	hashedPassword, err := utils.HashPassword(newUser.Password)

	if err != nil {
		log.Print("Error hashing password: ", err)
		return User{}, err
	}

	err = database.GetPreparedCreateUserQuery().QueryRow(&newUser.Email, hashedPassword).Scan(&user.ID, &user.Email)

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

func GetUserProfile(userID int) (UserProfile, error) {

	var userProfile = UserProfile{}

	var firstName, lastName sql.NullString

	err := database.GetPreparedGetUserProfile().QueryRow(userID).Scan(&userProfile.ID, &userProfile.Email, &firstName, &lastName)

	if err != nil && err != sql.ErrNoRows {
		return userProfile, err
	}

	userProfile.FirstName = utils.ConvertNullString(firstName)
	userProfile.LastName = utils.ConvertNullString(lastName)

	return userProfile, nil
}

func UpdateUserProfile(userID int, userProfile UpdateUserProfileRequest) (UserProfile, error) {

	updatedUser := UserProfile{}

	err := database.GetPreparedUpdateUserProfile().QueryRow(userProfile.FirstName, userProfile.LastName, userID).Scan(&updatedUser.ID, &updatedUser.Email, &updatedUser.FirstName, &updatedUser.LastName)

	if err != nil {
		log.Println("Error updating user profile: ", err)
		return updatedUser, err
	}

	return updatedUser, nil
}
