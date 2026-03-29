package service

import (
	"context"
	"moneytracker/internal/domain"
	"time"
)

type IExpenseRepository interface {
	AddExpense(ctx context.Context, payload *domain.DailyExpense) error
	GetWeeklyExpenses(ctx context.Context, week []string) (*domain.WeeklyExpense, error)
	DeleteExpensesByDate(ctx context.Context, date string) error
	GetDailyExpense(ctx context.Context, date string) ([]domain.DailyExpense, error)
	UpdateDailyExpense(ctx context.Context, uuid string, newAmount int32) error
	DeleteDailyExpense(ctx context.Context, uuid string) error
}

type ExpenseService struct {
	expenseRepository IExpenseRepository
}

func NewExpenseService(
	expenseRepository IExpenseRepository,
) *ExpenseService {
	return &ExpenseService{
		expenseRepository: expenseRepository,
	}
}

func (s *ExpenseService) AddExpense(ctx context.Context, payload *domain.DailyExpense) error {
	err := s.expenseRepository.AddExpense(ctx, payload)
	if err != nil {
		// TODO: wrap error
		return err
	}

	return nil
}

func (s *ExpenseService) GetWeeklyExpenses(ctx context.Context, date string) (*domain.WeeklyExpense, error) {
	week, err := s.getWeek(date)
	if err != nil {
		return nil, err
	}

	weekExpenses, err := s.expenseRepository.GetWeeklyExpenses(ctx, week)
	if err != nil {
		return nil, err
	}

	return weekExpenses, nil
}

func (s *ExpenseService) DeleteExpensesByDate(ctx context.Context, date string) error {
	err := s.expenseRepository.DeleteExpensesByDate(ctx, date)
	if err != nil {
		return err
	}

	return nil
}

func (s *ExpenseService) GetDailyExpense(ctx context.Context, date string) ([]domain.DailyExpense, error) {
	dailyExpenses, err := s.expenseRepository.GetDailyExpense(ctx, date)
	if err != nil {
		return nil, err
	}

	return dailyExpenses, nil
}

func (s *ExpenseService) UpdateDailyExpense(ctx context.Context, uuid string, newAmount int32) error {
	err := s.expenseRepository.UpdateDailyExpense(ctx, uuid, newAmount)
	if err != nil {
		return err
	}

	return nil
}

func (s *ExpenseService) DeleteDailyExpense(ctx context.Context, uuid string) error {
	err := s.expenseRepository.DeleteDailyExpense(ctx, uuid)
	if err != nil {
		return err
	}

	return nil
}

func (s *ExpenseService) getWeek(dateStr string) ([]string, error) {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
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
