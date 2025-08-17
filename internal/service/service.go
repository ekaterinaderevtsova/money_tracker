package service

import (
	"context"
	"moneytracker/internal/domain"
	"moneytracker/internal/repository"

	"go.uber.org/zap"
)

type IExpenseService interface {
	AddExpense(ctx context.Context, payload *domain.DailyExpense) error
	GetWeeklyExpenses(ctx context.Context, date string) (*domain.WeeklyExpense, error)
}

type Service struct {
	ExpenseService IExpenseService
}

func NewService(repo *repository.Repository, logger *zap.Logger) *Service {
	return &Service{
		ExpenseService: NewExpenseService(repo.ExpenseRepository, logger),
	}
}
