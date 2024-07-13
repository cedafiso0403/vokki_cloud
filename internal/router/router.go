package router

import (
	"net/http"
	"vokki_cloud/internal/handlers"
	"vokki_cloud/internal/middleware"

	"github.com/gorilla/mux"
)

func SetupRouter() *mux.Router {

	r := mux.NewRouter()

	r.HandleFunc("/v1/register", handlers.RegisterUser).Methods("POST")
	r.HandleFunc("/v1/verify", middleware.EmailVerificationMiddleware(http.HandlerFunc(handlers.VerifyUser))).Methods("GET")

	return r
}
