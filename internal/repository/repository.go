package repository

import (
	"context"
	"moneytracker/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:generate mockgen -source=repository.go -destination=mocks/mock.go

type IExpenseRepository interface {
	AddExpense(ctx context.Context, payload *domain.DailyExpense) error
	GetWeeklyExpenses(ctx context.Context, week []string) (*domain.WeeklyExpense, error)
}

type Repository struct {
	ExpenseRepository IExpenseRepository
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		ExpenseRepository: NeExpenseRepository(db),
	}
}
