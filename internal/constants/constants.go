package vokki_constants

type ContextKey string

type TokenType string

// Base path for the API

const (
	BaseHost   string = "vokki.net"
	BasePath   string = "/api/v1"
	BaseScheme string = "https"
)

const (
	EmailToken    TokenType = "email"
	AuthToken     TokenType = "auth"
	RefreshToken  TokenType = "refresh"
	ResetPassword TokenType = "reset"
)

const (
	UserIDKey   ContextKey = "userID"
	TokenKey    ContextKey = "token"
	TokenIssuer string     = "Vokki"
)

// Routes
const (
	RouteLogin             = "/login"
	RouteRegister          = "/register"
	RouteVerifyEmail       = "/verify"
	RouteResetPassword     = "/reset-password"
	RouteAlive             = "/alive"
	RouteTermAndConditions = "/terms-and-conditions"
	RouteCreateNewPassword = "/create-new-password"
	RouteLandingPage       = "/"
	RouteUser              = "/user"
	RouteWords             = "/words"
)
