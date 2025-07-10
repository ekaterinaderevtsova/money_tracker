package service

import (
	"context"
	"moneytracker/internal/domain"
	"moneytracker/internal/repository"

	"go.uber.org/zap"
)

type IWeekManager interface {
	IsCurrentWeek(date string) bool
	GetCurrentWeek() []string
	InitializeWeek(ctx context.Context) error
	ArchiveCurrentWeek(ctx context.Context) error
	ResetForNewWeek(ctx context.Context) error
}

type IWeeklyScheduler interface {
	Start(ctx context.Context) error
}

type ISpendingService interface {
	AddSpending(ctx context.Context, payload *domain.DaySpendings) error
	GetWeekSpendings(ctx context.Context, date string) (*domain.WeekSpendings, error)
}

type Service struct {
	IWeekManager
	IWeeklyScheduler
	ISpendingService
}

func NewService(
	ctx context.Context,
	repo *repository.Repository,
	logger *zap.Logger,
) (*Service, error) {
	weekManager, err := NewWeekManager(ctx,
		repo.ICurrentSpendingRepository,
		repo.IArchiveSpendingRepository,
		logger,
	)

	if err != nil {
		return nil, err
	}

	spendingService := NewSpendingService(
		weekManager,
		repo.ICurrentSpendingRepository,
		repo.IArchiveSpendingRepository,
	)

	weeklyScheduler := NewWeeklyScheduler(weekManager, logger)

	return &Service{
		IWeekManager:     weekManager,
		IWeeklyScheduler: weeklyScheduler,
		ISpendingService: spendingService,
	}, nil
}

func (ss *Service) Start(ctx context.Context) error {
	return ss.IWeeklyScheduler.Start(ctx)
}
