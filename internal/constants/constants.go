package vokki_constants

type ContextKey string

type TokenType string

const (
	UserIDKey    ContextKey = "userID"
	TokenKey     ContextKey = "token"
	EmailToken   TokenType  = "email"
	AuthToken    TokenType  = "auth"
	RefreshToken TokenType  = "refresh"
	Issuer       string     = "Vokki"
)
