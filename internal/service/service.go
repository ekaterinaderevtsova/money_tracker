package service

import "cmd/main.go/internal/repository"

type Service struct {
	SpendingService *SpendingService
}

func NewService(repository *repository.Repository) *Service {
	return &Service{
		SpendingService: NewSpendingService(repository.SpendingRepository),
	}
}
