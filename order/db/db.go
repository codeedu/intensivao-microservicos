package db

import (
	"github.com/go-redis/redis/v7"
	"os"
)

func Connect() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST"),
		Password: "",
		DB:       0,
	})
	return client
}
