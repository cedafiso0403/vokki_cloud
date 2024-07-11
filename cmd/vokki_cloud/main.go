package main

import (
	"log"
	"net/http"
	"vokki_cloud/internal/database"
	"vokki_cloud/internal/router"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	database.Connect()

	r := router.SetupRouter()

	log.Fatal(http.ListenAndServe(":8000", r))
}
