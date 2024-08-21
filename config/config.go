package config

import (
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func LoadConfig() (*http.Server, *os.File, error) {

	// Open log file for writing
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	multiWriter := io.MultiWriter(os.Stdout, file)
	log.SetOutput(multiWriter)

	// Load environment variables from .env file
	// if err := godotenv.Load(); err != nil {
	// 	log.Printf("Error loading .env file: %v", err)
	// 	return nil, nil, err
	// }

	log.Println("Successfully loaded .env file")

	// Setup server configuration
	server := &http.Server{
		Addr:              ":8000", // Consider reading this from an env variable if needed
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Println("Server loaded")
	log.Println("Server running on port 8000")

	// Return the server and the log file so that it can be closed later
	return server, file, nil
}
