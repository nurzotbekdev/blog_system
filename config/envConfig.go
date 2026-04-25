package config

import (
	"log"

	"github.com/joho/godotenv"
)

func EnvConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Failed to load .env file")
	}
}
