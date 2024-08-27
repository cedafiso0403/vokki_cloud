package router

import (
	"net/http"
	vokki_constants "vokki_cloud/internal/constants"
	"vokki_cloud/internal/handlers"
	"vokki_cloud/internal/middleware"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

func SetupRouter() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc(vokki_constants.RouteLandingPage, handlers.LandingPage).Methods("GET")
	r.HandleFunc(vokki_constants.RouteTermAndConditions, handlers.TermAndConditions).Methods("GET")

	// Serve Swagger documentation
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// API routes with prefix /api/v1
	apiRouterPublic := r.PathPrefix("/api/v1").Subrouter()
	apiRouterPrivate := r.PathPrefix("/api/v1").Subrouter()

	apiRouterPrivate.Use(middleware.AuthMiddleware)

	// Public routes
	apiRouterPublic.HandleFunc(vokki_constants.RouteLogin, handlers.Login).Methods("POST")
	apiRouterPublic.HandleFunc(vokki_constants.RouteRegister, handlers.RegisterUser).Methods("POST")
	apiRouterPublic.HandleFunc(vokki_constants.RouteVerifyEmail, middleware.EmailVerificationMiddleware(http.HandlerFunc(handlers.VerifyUser))).Methods("GET")
	apiRouterPublic.HandleFunc(vokki_constants.RouteResetPassword, handlers.RequestResetPassword).Methods("POST")

	apiRouterPublic.HandleFunc(vokki_constants.RouteAlive, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Alive"))
	}).Methods("GET")

	// Private routes
	// User routes
	apiRouterPrivate.HandleFunc(vokki_constants.RouteUser, handlers.GetUser).Methods("GET")
	apiRouterPrivate.HandleFunc(vokki_constants.RouteUser, handlers.UpdateUser).Methods("PUT")

	return r
}
