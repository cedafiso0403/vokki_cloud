package config

import (
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func LoadConfig() (*http.Server, error) {

	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer file.Close()

	multiWriter := io.MultiWriter(os.Stdout, file)

	log.SetOutput(multiWriter)

	if err := godotenv.Load(); err != nil {
		log.Printf("Error loading .env file: %v", err)
		return nil, err
	}

	log.Println("Successfully loaded .env file")

	server := &http.Server{
		Addr:              ":8000",
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Println("Server loaded")
	log.Println("Server running on port 8000")

	return server, nil
}
