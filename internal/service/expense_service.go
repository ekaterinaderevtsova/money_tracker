package service

import (
	"context"
	"moneytracker/internal/domain"
	"moneytracker/internal/repository"
	"time"

	"go.uber.org/zap"
)

type ExpenseService struct {
	expenseRepository repository.IExpenseRepository
	logger            *zap.Logger
}

func NewExpenseService(
	expenseRepository repository.IExpenseRepository,
	logger *zap.Logger,
) *ExpenseService {
	return &ExpenseService{
		expenseRepository: expenseRepository,
		logger:            logger,
	}
}

func (s *ExpenseService) AddExpense(ctx context.Context, payload *domain.DailyExpense) error {
	err := s.expenseRepository.AddExpense(ctx, payload)
	if err != nil {
		s.logger.Error("Error adding expense to db",
			zap.Error(err),
			zap.String("date", payload.Date),
			zap.Int32("amount", payload.Amount),
		)
		return err
	}

	s.logger.Info("New expense added",
		zap.String("date", payload.Date),
		zap.Int32("amount", payload.Amount),
	)
	return nil
}

func (s *ExpenseService) GetWeeklyExpenses(ctx context.Context, date string) (*domain.WeeklyExpense, error) {
	week, err := s.getWeek(date)
	if err != nil {
		s.logger.Error("Error getting week by date",
			zap.Error(err),
			zap.String("date", date),
		)
		return nil, err
	}

	weekExpenses, err := s.expenseRepository.GetWeeklyExpenses(ctx, week)
	if err != nil {
		s.logger.Error("Error getting week expenses",
			zap.Error(err),
			zap.String("date", date),
		)
		return nil, err
	}

	return weekExpenses, nil
}

func (s *ExpenseService) getWeek(dateStr string) ([]string, error) {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		s.logger.Error("failed to parse date", zap.Error(err), zap.String("date", dateStr))
		return nil, err
	}

	daysSinceMonday := (int(date.Weekday()) + 6) % 7
	startOfWeek := date.AddDate(0, 0, -int(daysSinceMonday))

	week := make([]string, 0, 7)

	for i := 0; i < 7; i++ {
		day := startOfWeek.AddDate(0, 0, i)
		dayStr := day.Format("2006-01-02")
		week = append(week, dayStr)
	}

	return week, nil
}
