package migrations

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

const redisMigrationKey = "migration:redis_version"

func RunRedisMigrations(ctx context.Context, redisDb *redis.Client, currentWeek []string) error {
	version, err := redisDb.Get(ctx, redisMigrationKey).Result()
	if err == redis.Nil {
		version = "0"
	} else if err != nil {
		return fmt.Errorf("failed to get redis migration version: %w", err)
	}

	switch version {
	case "0":
		log.Println("Running Redis migration v1: init week days")

		for _, day := range currentWeek {
			dayKey := "spendings:" + day
			totalKey := "total:" + day

			// только если ключей ещё нет
			if exists, _ := redisDb.Exists(ctx, dayKey).Result(); exists == 0 {
				if err := redisDb.RPush(ctx, dayKey, 0).Err(); err != nil {
					return fmt.Errorf("failed to init %s: %w", dayKey, err)
				}
			}

			if exists, _ := redisDb.Exists(ctx, totalKey).Result(); exists == 0 {
				if err := redisDb.Set(ctx, totalKey, 0, 0).Err(); err != nil {
					return fmt.Errorf("failed to init %s: %w", totalKey, err)
				}
			}
		}

		if err := redisDb.Set(ctx, redisMigrationKey, "1", 0).Err(); err != nil {
			return fmt.Errorf("failed to update redis version: %w", err)
		}

		log.Println("Redis migration to v1 completed")
	}

	// future:
	// case "1":
	//   // инициализация новых ключей

	return nil
}
