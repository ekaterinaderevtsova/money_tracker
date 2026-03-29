package service

import (
	"moneytracker/internal/repository"
)

type Service struct {
	ExpenseService *ExpenseService
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		ExpenseService: NewExpenseService(repo.ExpenseRepository),
	}
}
