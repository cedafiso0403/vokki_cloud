package main

import (
	"log"
	"vokki_cloud/config"
	_ "vokki_cloud/docs" // This is required for Swagger to find your docs
	"vokki_cloud/internal/database"
	"vokki_cloud/internal/router"
	"vokki_cloud/internal/shared"

	_ "github.com/lib/pq"
)

// @title Vokki Cloud API
// @version 1.0
// @description This is the API for Vokki mobile app
func main() {

	server, err := config.LoadConfig()

	if err != nil {
		log.Fatal("Error loading config: ", err)
		server.Close()
		return
	}

	shared.InitializeTokenManager()

	database.Connect()

	r := router.SetupRouter()

	server.Handler = r

	log.Println(server.ListenAndServe())

}
