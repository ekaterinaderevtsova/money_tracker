package repository

import (
	"context"
	"moneytracker/internal/domain"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type ICurrentSpendingRepository interface {
	InitNewWeek(ctx context.Context, week []string) error
	FlushAll(ctx context.Context) error
	AddSpending(ctx context.Context, payload *domain.DaySpendings) error
	GetWeekSpendings(ctx context.Context, week []string) (*domain.WeekSpendings, error)
}

type IArchiveSpendingRepository interface {
	AddSpending(ctx context.Context, payload *domain.DaySpendings) error
	GetWeekSpendings(ctx context.Context, date time.Time) (*domain.WeekSpendings, error)
}

type Repository struct {
	IArchiveSpendingRepository
	ICurrentSpendingRepository
}

func NewRepository(ctx context.Context, db *pgxpool.Pool, redisDb *redis.Client) *Repository {
	return &Repository{
		IArchiveSpendingRepository: NewArchiveSpendingRepository(db),
		ICurrentSpendingRepository: NewCurrentSpendingRepository(redisDb),
	}
}
