package main

import (
	"log"
	"vokki_cloud/config"
	"vokki_cloud/docs" // This is required for Swagger to find your docs
	"vokki_cloud/internal/database"
	"vokki_cloud/internal/router"
	"vokki_cloud/internal/shared"

	_ "github.com/lib/pq"
)

func main() {

	server, err := config.LoadConfig()

	docs.SwaggerInfo.Title = "Vokki Cloud API"
	docs.SwaggerInfo.Description = "This is the API for Vokki mobile app."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "vokki.net"
	docs.SwaggerInfo.BasePath = "/v1"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

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
