package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"vokki_cloud/internal/database"
	"vokki_cloud/internal/router"

	_ "github.com/lib/pq"
)

func main() {
	// Open a file for logging
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer file.Close()

	multiWriter := io.MultiWriter(os.Stdout, file)

	// Set the log output to the file
	log.SetOutput(multiWriter)
	log.Println("Server started on port 8000")

	database.Connect()

	r := router.SetupRouter()

	log.Fatal(http.ListenAndServe(":8000", r))

}
