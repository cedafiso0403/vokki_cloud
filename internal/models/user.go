package models

import (
	"errors"
	vokki_constants "vokki_cloud/internal/constants"
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

func ActivateUser(userID int64, token string) error {

	db := database.GetDB()

	verificationToken := AuthToken{}

	preparedTokenExistsQuery, err := db.Prepare("SELECT verification_token, revoked_at, token_type, user_id FROM user_tokens WHERE user_id=$1 AND verification_token=$2")

	if err != nil {
		return err
	}

	defer preparedTokenExistsQuery.Close()

	err = preparedTokenExistsQuery.QueryRow(userID, token).Scan(&verificationToken.Token, &verificationToken.RevokedAt, &verificationToken.TokenType, &verificationToken.UserID)

	if err != nil {
		return err
	}

	if verificationToken.RevokedAt.Valid {
		return errors.New("token has been revoked")
	}

	if verificationToken.TokenType != vokki_constants.EmailToken {
		return errors.New("invalid token type")
	}

	preparedUpdateTokenQuery, err := db.Prepare("UPDATE user_tokens SET revoked_at=$1 WHERE user_id=$2 AND verification_token=$3")

	if err != nil {
		return err
	}

	defer preparedUpdateTokenQuery.Close()

	_, err = preparedUpdateTokenQuery.Exec("now()", userID, token)

	if err != nil {
		return err
	}

	preparedActivateUserQuery, err := db.Prepare("UPDATE users SET activated=true WHERE id=$1")

	if err != nil {
		return err
	}

	defer preparedActivateUserQuery.Close()

	_, err = preparedActivateUserQuery.Exec(userID)

	if err != nil {
		return err
	}

	return nil
}
