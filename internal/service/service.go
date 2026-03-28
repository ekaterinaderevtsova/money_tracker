package service

import (
	"moneytracker/internal/repository"

	"go.uber.org/zap"
)

type Service struct {
	ExpenseService *ExpenseService
}

func NewService(repo *repository.Repository, logger *zap.Logger) *Service {
	return &Service{
		ExpenseService: NewExpenseService(repo.ExpenseRepository, logger),
	}
}
