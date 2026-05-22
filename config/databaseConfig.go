package config

import (
	"blog_system/logging"
	"fmt"
	"os"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func DatabaseConfig() {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	logging.Log.Info("Initializing database connection", zap.String("host", host), zap.String("port", port), zap.String("database", dbname), zap.String("user", user))
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable ", host, user, password, dbname, port)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logging.Log.Fatal("Failed to connect to database", zap.String("host", host), zap.String("port", port), zap.String("database", dbname), zap.Error(err))
	}

	sqlDB, err := DB.DB()
	if err != nil {
		logging.Log.Fatal("Failed to get generic database object", zap.Error(err))
	}

	if err = sqlDB.Ping(); err != nil {
		logging.Log.Fatal("Database ping failed", zap.Error(err))
	}

	logging.Log.Info("Database connected successfully", zap.String("host", host), zap.String("port", port), zap.String("database", dbname))
}
