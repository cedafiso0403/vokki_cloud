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

	// Serve Swagger documentation
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// Public routes
	r.HandleFunc(vokki_constants.RouteLogin, handlers.Login).Methods("POST")
	r.HandleFunc(vokki_constants.RouteRegister, handlers.RegisterUser).Methods("POST")
	r.HandleFunc(vokki_constants.RouteVerifyEmail, middleware.EmailVerificationMiddleware(http.HandlerFunc(handlers.VerifyUser))).Methods("GET")

	// API routes with prefix /api/v1
	apiRouter := r.PathPrefix("/api/v1").Subrouter()
	apiRouter.HandleFunc(vokki_constants.RouteAlive, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods("GET")

	return r
}
