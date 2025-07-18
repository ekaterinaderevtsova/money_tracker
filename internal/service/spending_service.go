package service

import (
	"context"
	"moneytracker/internal/domain"
	"moneytracker/internal/repository"
)

type SpendingService struct {
	weekManager               IWeekManager
	currentSpendingRepository repository.ICurrentSpendingRepository
	archiveSpendingRepository repository.IArchiveSpendingRepository
}

func NewSpendingService(
	weekManager IWeekManager,
	currentSpendingRepository repository.ICurrentSpendingRepository,
	archiveSpendingRepository repository.IArchiveSpendingRepository,
) *SpendingService {
	return &SpendingService{
		weekManager:               weekManager,
		currentSpendingRepository: currentSpendingRepository,
		archiveSpendingRepository: archiveSpendingRepository,
	}
}

func (s *SpendingService) AddSpending(ctx context.Context, payload *domain.DaySpendings) error {
	if s.weekManager.IsCurrentWeek(payload.Day) {
		return s.currentSpendingRepository.AddSpending(ctx, payload)
	}
	return s.archiveSpendingRepository.AddSpending(ctx, payload)
}

func (s *SpendingService) GetWeekSpendings(ctx context.Context, date string) (*domain.WeekSpendings, error) {
	if s.weekManager.IsCurrentWeek(date) {
		weekSpendings, err := s.currentSpendingRepository.GetWeekSpendings(ctx, s.weekManager.GetCurrentWeek())
		if err != nil {
			return nil, err
		}
		return weekSpendings, nil
	}

	week, err := s.weekManager.GetArchiveWeek(date)
	if err != nil {
		return nil, err
	}

	weekSpendings, err := s.archiveSpendingRepository.GetWeekSpendings(ctx, week)
	if err != nil {
		return nil, err
	}

	return weekSpendings, nil
}
