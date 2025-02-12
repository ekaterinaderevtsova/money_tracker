package service

import (
	"cmd/main.go/internal/domain"
	"context"
	"time"
)

type ISpendingRepository interface {
	AddSpending(ctx context.Context, payload *domain.AddSpending) error
	//	GetAllSpendings(ctx context.Context) ([]domain.AddSpending, error)
	GetDaySpendings(ctx context.Context, date time.Time) (int32, error)
	GetWeekSpendings(ctx context.Context, date time.Time) (*domain.WeekSpendings, error)
	GetMonthSpendings(ctx context.Context, date time.Time) ([]domain.WeekTotalSpending, error)
}

type SpendingService struct {
	spendingRepository ISpendingRepository
}

func NewSpendingService(spendingRepository ISpendingRepository) *SpendingService {
	return &SpendingService{spendingRepository: spendingRepository}
}

func (ss *SpendingService) AddSpending(ctx context.Context, payload *domain.AddSpending) (int32, error) {
	err := ss.spendingRepository.AddSpending(ctx, payload)
	if err != nil {
		return 0, err
	}

	daySpendings, err := ss.spendingRepository.GetDaySpendings(ctx, payload.Date)
	if err != nil {
		return 0, err
	}

	return daySpendings, nil
}

func (ss *SpendingService) GetWeekSpendings(ctx context.Context, date time.Time) (*domain.WeekSpendings, error) {
	weekSpendings, err := ss.spendingRepository.GetWeekSpendings(ctx, date)
	if err != nil {
		return nil, err
	}

	return weekSpendings, nil
}

func (ss *SpendingService) GetMonthSpendings(ctx context.Context, date time.Time) ([]domain.WeekTotalSpending, error) {
	monthSpendings, err := ss.spendingRepository.GetMonthSpendings(ctx, date)
	if err != nil {
		return nil, err
	}

	return monthSpendings, nil
}
