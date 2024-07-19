package router

import (
	"net/http"
	"vokki_cloud/internal/handlers"
	"vokki_cloud/internal/middleware"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

func SetupRouter() *mux.Router {

	r := mux.NewRouter()

	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	//apiRouter := r.PathPrefix("/api/v1").Subrouter()

	//Public routes
	r.HandleFunc("/login", handlers.Login).Methods("POST")
	r.HandleFunc("/register", handlers.RegisterUser).Methods("POST")
	r.HandleFunc("/verify", middleware.EmailVerificationMiddleware(http.HandlerFunc(handlers.VerifyUser))).Methods("GET")

	return r
}
