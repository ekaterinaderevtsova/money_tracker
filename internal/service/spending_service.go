package service

import (
	"context"
	"fmt"
	"log"
	"moneytracker/internal/domain"
	"moneytracker/internal/repository"
	"time"
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

	fmt.Println("date: ", date)

	dateParsed, err := time.Parse("2006-01-02", date)
	if err != nil {
		log.Println("Invalid date format")
		return nil, err
	}

	weekSpendings, err := s.archiveSpendingRepository.GetWeekSpendings(ctx, dateParsed)

	if err != nil {
		return nil, err
	}

	return weekSpendings, nil
}
