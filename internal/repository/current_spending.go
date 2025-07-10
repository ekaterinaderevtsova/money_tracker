package repository

import (
	"context"
	"fmt"
	"moneytracker/internal/domain"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type CurrentSpendingRepository struct {
	redisDb *redis.Client
}

func NewCurrentSpendingRepository(redisDb *redis.Client) *CurrentSpendingRepository {
	return &CurrentSpendingRepository{
		redisDb: redisDb,
	}
}

func (r *CurrentSpendingRepository) InitNewWeek(ctx context.Context, week []string) error {
	for _, day := range week {
		dayKey := domain.SpendingsKey + day
		totalKey := domain.TotalKey + day

		exists, err := r.redisDb.Exists(ctx, dayKey).Result()
		if err != nil {
			return err
		}
		if exists == 0 {
			if err := r.redisDb.RPush(ctx, dayKey, 0).Err(); err != nil {
				return err
			}
		}

		exists, err = r.redisDb.Exists(ctx, totalKey).Result()
		if err != nil {
			return err
		}
		if exists == 0 {
			if err := r.redisDb.Set(ctx, totalKey, 0, 0).Err(); err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *CurrentSpendingRepository) FlushAll(ctx context.Context) error {
	// Get all keys that match the spending patterns
	spendingKeys, err := r.redisDb.Keys(ctx, domain.SpendingsKey+"*").Result()
	if err != nil {
		return fmt.Errorf("failed to get spending keys: %w", err)
	}

	totalKeys, err := r.redisDb.Keys(ctx, domain.TotalKey+"*").Result()
	if err != nil {
		return fmt.Errorf("failed to get total keys: %w", err)
	}

	// Combine all keys to delete
	allKeys := append(spendingKeys, totalKeys...)

	// Delete all keys if any exist
	if len(allKeys) > 0 {
		if err := r.redisDb.Del(ctx, allKeys...).Err(); err != nil {
			return fmt.Errorf("failed to delete keys: %w", err)
		}
	}

	return nil
}

func (sr *CurrentSpendingRepository) AddSpending(ctx context.Context, payload *domain.DaySpendings) error {
	if err := sr.redisDb.LPush(ctx, domain.SpendingsKey+payload.Day, payload.Sum).Err(); err != nil {
		return fmt.Errorf("failed to set new sum in redis: %w", err)
	}

	if err := sr.redisDb.IncrBy(ctx, domain.TotalKey+payload.Day, int64(payload.Sum)).Err(); err != nil {
		return fmt.Errorf("failed to set new sum in redis: %w", err)
	}

	return nil
}

func (sr *CurrentSpendingRepository) GetWeekSpendings(ctx context.Context, week []string) (*domain.WeekSpendings, error) {
	var weekSpendings domain.WeekSpendings
	var total int32
	var daysCount int32
	today := time.Now().Format("2006-01-02")

	for i, day := range week {

		sumStr, err := sr.redisDb.Get(ctx, domain.TotalKey+day).Result()
		if err != nil {
			// TODO: log
			continue
		}

		sum, _ := strconv.Atoi(sumStr)

		weekSpendings.DaySpendings[i] = domain.DaySpendings{
			Day: day,
			Sum: int32(sum),
		}

		total += int32(sum)

		if today >= day {
			daysCount++
		}
	}

	weekSpendings.Total = total
	weekSpendings.Average = 0
	if daysCount > 0 {
		weekSpendings.Average = total / daysCount
	}
	return &weekSpendings, nil
}
