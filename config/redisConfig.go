package config

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client
var Ctx = context.Background()

func RedisConfig() {
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")
	password := os.Getenv("REDIS_PASSWORD")
	dbStr := os.Getenv("REDIS_DB")

	log.Printf("Initializing Redis connection (host: %s, port: %s, db: %s)", host, port, dbStr)

	db, err := strconv.Atoi(dbStr)
	if err != nil {
		log.Printf("WARNING: Invalid REDIS_DB value '%s', using default database 0: %v", dbStr, err)
		db = 0
	}

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: password,
		DB:       db,
	})

	_, err = RedisClient.Ping(Ctx).Result()
	if err != nil {
		log.Fatalf("Redis connection failed (host: %s, port: %s, db: %d): %v", host, port, db, err)
	}

	log.Printf("Redis connected successfully (host: %s, port: %s, db: %d)", host, port, db)
}
