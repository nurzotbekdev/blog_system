package config

import (
	"log"

	"github.com/joho/godotenv"
)

func EnvConfig() {
	log.Println("Loading environment variables from .env file")

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Failed to load .env file: %v", err)
	}

	log.Println(".env file loaded successfully")
}
