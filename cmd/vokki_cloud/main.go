package main

import (
	"log"
	"vokki_cloud/config"
	"vokki_cloud/internal/database"
	"vokki_cloud/internal/router"
	"vokki_cloud/internal/shared"

	_ "github.com/lib/pq"
)

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
