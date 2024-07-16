package auth_error

import "errors"

var (
	ErrUserNotFound         = errors.New("user not found")
	ErrIncorrectCredentials = errors.New("incorrect credentials")
	ErrInvalidToken         = errors.New("invalid token")
	ErrExpiredToken         = errors.New("expired token")
	ErrEmailAlreadyInUse    = errors.New("email already in use")
	ErrInvalidEmailFormat   = errors.New("invalid email format")
	ErrPasswordsDoNotMatch  = errors.New("passwords do not match")
	ErrUnauthorized         = errors.New("unauthorized")
	ErrInternalServer       = errors.New("internal server error")
	ErrInvalidTokenIssuer   = errors.New("invalid token issuer")
)
