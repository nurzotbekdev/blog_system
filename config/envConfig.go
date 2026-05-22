package config

import (
	"blog_system/logging"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func EnvConfig() {
	logging.Log.Info("Loading environment variables from .env file")

	err := godotenv.Load()
	if err != nil {
		logging.Log.Fatal("Failed to load .env file", zap.Error(err))
	}

	logging.Log.Info(".env file loaded successfully")
}
