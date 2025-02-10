package repository

import "github.com/jackc/pgx/v5/pgxpool"

type Repository struct {
	SpendingRepository *SpendingRepository
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{
		SpendingRepository: NewSpendingRepository(pool),
	}
}
