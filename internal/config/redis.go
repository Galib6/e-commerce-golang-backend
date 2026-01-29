package config

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client

var ctx = context.Background()

func ConnectRedis() {
	RDB = redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Ideally, load this from env variables
		// Addr:     "redis:6379", // redis is the service name in docker-compose
		Password: "123456",
		Username: "galib",
		DB:       0,
	})

	// Test the connection immediately
	if err := RDB.Ping(ctx).Err(); err != nil {
		panic(fmt.Sprintf("Failed to connect to Redis: %v", err))
	}

	fmt.Println("Connected to Redis successfully")
}
