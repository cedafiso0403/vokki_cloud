package services

import (
	"database/sql"
	"errors"
	"vokki_cloud/internal/auth_error"
	vokki_constants "vokki_cloud/internal/constants"
	"vokki_cloud/internal/database"
	"vokki_cloud/internal/models"
	"vokki_cloud/internal/shared"
	"vokki_cloud/internal/utils"
)

type Credentials struct {
	Email    string `json:"email" example:"user@domain.com"`
	Password string `json:"password" example:"password"`
}

type NewPasswordEmailRequest struct {
	Email string `json:"email" binding:"required" example:"user@domain.com"`
}

type NewPasswordRequest struct {
	Password             string `json:"password" binding:"required" example:"password"`
	ConfirmationPassword string `json:"confirmation_password" binding:"required" example:"password"`
	Token                string `json:"token" binding:"required" example:"eyJhbGciOiAiSFMyNTeHBpcmVk1fjr-mL6l7QJbdFfL86D4HK4XsEFPfSb2X8"`
}

func Authenticate(credentials Credentials) (int, string, error) {
	db := database.GetDB()
	var user models.User
	err := db.QueryRow("SELECT id, email, hashed_password FROM users WHERE email = $1", credentials.Email).Scan(&user.ID, &user.Email, &user.Password)
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

func ActivateUser(userID int64, token string) error {

	db := database.GetDB()

	verificationToken := models.AuthToken{}

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

func RequestNewPasswordEmail(email string) error {
	db := database.GetDB()

	user := models.User{}

	preparedQuery, err := db.Prepare("SELECT id, email FROM users WHERE email=$1")

	if err != nil {
		return err
	}

	defer preparedQuery.Close()

	err = preparedQuery.QueryRow(email).Scan(&user.ID, &user.Email)

	if err != nil {
		return err
	}

	resetPasswordToken, err := utils.GenerateJWT(int64(user.ID))

	if err != nil {
		return err
	}

	err = models.StoreToken(int64(user.ID), resetPasswordToken, vokki_constants.ResetPassword)

	if err != nil {
		return err
	}

	shared.GetTokenManager().AddToken(resetPasswordToken)

	SendPasswordResetEmail(user, resetPasswordToken)

	return nil

}
