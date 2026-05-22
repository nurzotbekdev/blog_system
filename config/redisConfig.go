package config

import (
	"blog_system/logging"
	"context"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var RedisClient *redis.Client
var Ctx = context.Background()

func RedisConfig() {
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")
	password := os.Getenv("REDIS_PASSWORD")
	dbStr := os.Getenv("REDIS_DB")

	logging.Log.Info("Initializing Redis connection", zap.String("host", host), zap.String("port", port), zap.String("db", dbStr))

	db, err := strconv.Atoi(dbStr)
	if err != nil {
		logging.Log.Warn("Invalid REDIS_DB value, using default database 0", zap.String("redis_db", dbStr), zap.Error(err))
		db = 0
	}

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: password,
		DB:       db,
	})

	_, err = RedisClient.Ping(Ctx).Result()
	if err != nil {
		logging.Log.Fatal("Redis connection failed", zap.String("host", host), zap.String("port", port), zap.Int("db", db), zap.Error(err))
	}

	logging.Log.Info("Redis connected successfully", zap.String("host", host), zap.String("port", port), zap.Int("db", db))
}
