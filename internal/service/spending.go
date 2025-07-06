package service

import (
	"cmd/main.go/internal/domain"
	"context"
	"fmt"
	"time"
)

type ISpendingRepository interface {
	IsCurrentWeek(date string) bool
	AddSpending(ctx context.Context, payload *domain.DaySpendings) error
	AddCurrentWeekSpending(ctx context.Context, payload *domain.DaySpendings) error
	GetWeekSpendings(ctx context.Context, date time.Time) (*domain.WeekSpendings, error)
	GetCurrentWeekSpendings(ctx context.Context) (*domain.WeekSpendings, error)
}

type SpendingService struct {
	spendingRepository ISpendingRepository
}

func NewSpendingService(spendingRepository ISpendingRepository) *SpendingService {
	return &SpendingService{spendingRepository: spendingRepository}
}

func (ss *SpendingService) AddSpending(ctx context.Context, payload *domain.DaySpendings) error {
	if ss.spendingRepository.IsCurrentWeek(payload.Day) {
		return ss.spendingRepository.AddCurrentWeekSpending(ctx, payload)
	}
	return ss.spendingRepository.AddSpending(ctx, payload)
}

func (ss *SpendingService) GetWeekSpendings(ctx context.Context, date string) (*domain.WeekSpendings, error) {
	if ss.spendingRepository.IsCurrentWeek(date) {
		fmt.Println("CURRENT WEEK")
		weekSpendings, err := ss.spendingRepository.GetCurrentWeekSpendings(ctx)
		if err != nil {
			return nil, err
		}
		return weekSpendings, nil
	}

	dateParsed, err := time.Parse("2006-01-02", date)
	if err != nil {
		//	sh.logger.Error("Invalid date format", zap.Error(err))
		return nil, err
	}

	weekSpendings, err := ss.spendingRepository.GetWeekSpendings(ctx, dateParsed)
	if err != nil {
		return nil, err
	}
	return weekSpendings, nil
}
