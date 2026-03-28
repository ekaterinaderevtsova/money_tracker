package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

//go:generate mockgen -source=repository.go -destination=mocks/mock.go

type Repository struct {
	ExpenseRepository *ExpenseRepository
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		ExpenseRepository: NeExpenseRepository(db),
	}
}
