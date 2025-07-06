package service

import "cmd/main.go/internal/repository"

type Service struct {
	SpendingService *SpendingService
	ScriptService   *ScriptService
}

func NewService(repository *repository.Repository) *Service {
	return &Service{
		SpendingService: NewSpendingService(repository.SpendingRepository),
		ScriptService:   NewScriptService(repository.SpendingRepository),
	}
}
