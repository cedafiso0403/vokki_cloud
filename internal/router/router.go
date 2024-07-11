package router

import (
	"vokki_cloud/internal/handlers"

	"github.com/gorilla/mux"
)

func SetupRouter() *mux.Router {

	r := mux.NewRouter()
	r.HandleFunc("/v1/register", handlers.RegisterUser).Methods("POST")

	return r
}
