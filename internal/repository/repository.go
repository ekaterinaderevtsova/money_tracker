package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Repository struct {
	SpendingRepository *SpendingRepository
}

func NewRepository(ctx context.Context, db *pgxpool.Pool, redisDb *redis.Client) *Repository {
	return &Repository{
		SpendingRepository: NewSpendingRepository(ctx, db, redisDb),
	}
}
