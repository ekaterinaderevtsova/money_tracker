package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

func NewRedisConn(ctx context.Context, address, password string) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis_test:6379",
		Password: "password",
		DB:       0,
	})

	time.Sleep(1 * time.Second)

	// Verify connection
	if _, err := rdb.Ping(ctx).Result(); err != nil {
		fmt.Println("failed to ping redis")
		return nil, fmt.Errorf("failed to connect to Redis: %v", err)
	}

	log.Println("сonnected to redis")

	return rdb, nil
}
