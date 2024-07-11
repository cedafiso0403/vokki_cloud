package config

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadConfig() error {
	if err := godotenv.Load(); err != nil {
		log.Printf("Error loading .env file: %v", err)
		return err
	}

	//Add set port later
	return nil
}
