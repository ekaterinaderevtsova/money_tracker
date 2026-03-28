package service

import (
	"context"
	"moneytracker/internal/domain"
	mock_repository "moneytracker/internal/repository/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func initTestLogger(t *testing.T) *zap.Logger {
	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatalf("failed to initialize test logger: %v", err)
	}
	return logger
}

func TestAddSpending(t *testing.T) {
	ctx := context.Background()
	logger := initTestLogger(t)
	ctrl := gomock.NewController(t)

	defer func() {
		if err := logger.Sync(); err != nil {
			t.Logf("failed to sync logger: %v", err)
		}
	}()
	defer ctrl.Finish()

	mockExpenseRepo := mock_repository.NewMockIExpenseRepository(ctrl)
	expenseServie := ExpenseService{
		expenseRepository: mockExpenseRepo,
		logger:            logger,
	}

	type mockBehavior func(r *mock_repository.MockIExpenseRepository)

	testTable := []struct {
		name          string
		payload       *domain.DailyExpense
		behavior      mockBehavior
		expectedError error
	}{
		{
			name: "valid input",
			payload: &domain.DailyExpense{
				Date:   time.Now().Format("2006-01-02"),
				Amount: 1000,
			},
			behavior: func(r *mock_repository.MockIExpenseRepository) {
				r.EXPECT().AddExpense(ctx, &domain.DailyExpense{
					Date:   time.Now().Format("2006-01-02"),
					Amount: 1000,
				}).Return(nil)
			},
			expectedError: nil,
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			if tt.behavior != nil {
				tt.behavior(mockExpenseRepo)
			}
			err := expenseServie.AddExpense(ctx, tt.payload)
			if tt.expectedError == nil {
				assert.NoError(t, err)
			} else {
				assert.Equal(t, err, tt.expectedError)
			}
		})
	}
}

func TestGetWeekSpendings(t *testing.T) {
	ctx := context.Background()
	logger := initTestLogger(t)
	ctrl := gomock.NewController(t)

	defer func() {
		if err := logger.Sync(); err != nil {
			t.Logf("failed to sync logger: %v", err)
		}
	}()
	defer ctrl.Finish()

	mockExpenseRepo := mock_repository.NewMockIExpenseRepository(ctrl)
	expenseServie := ExpenseService{
		expenseRepository: mockExpenseRepo,
		logger:            logger,
	}

	type mockBehavior func(r *mock_repository.MockIExpenseRepository)

	testTable := []struct {
		name           string
		date           string
		behavior       mockBehavior
		expectedOutput *domain.WeeklyExpense
		expectedError  error
	}{
		{
			name: "archive week valid input",
			date: "2025-07-07",
			behavior: func(r *mock_repository.MockIExpenseRepository) {
				week := []string{
					"2025-07-07",
					"2025-07-08",
					"2025-07-09",
					"2025-07-10",
					"2025-07-11",
					"2025-07-12",
					"2025-07-13",
				}
				r.EXPECT().GetWeeklyExpenses(ctx, week).Return(&domain.WeeklyExpense{
					DailyExpenses: [7]domain.DailyExpense{
						{
							Date:   "2025-07-07",
							Amount: 1000,
						},
					},
					TotalAmount:  1000,
					AverageDaily: 1000,
				}, nil)
			},
			expectedOutput: &domain.WeeklyExpense{
				DailyExpenses: [7]domain.DailyExpense{
					{
						Date:   "2025-07-07",
						Amount: 1000,
					},
				},
				TotalAmount:  1000,
				AverageDaily: 1000,
			},
			expectedError: nil,
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			if tt.behavior != nil {
				tt.behavior(mockExpenseRepo)
			}
			receivedOutput, err := expenseServie.GetWeeklyExpenses(ctx, tt.date)
			assert.Equal(t, tt.expectedOutput, receivedOutput)
			if tt.expectedError == nil {
				assert.NoError(t, err)
			} else {
				assert.Equal(t, err, tt.expectedError)
			}
		})
	}
}

func TestDeleteExpensesByDate(t *testing.T) {
	ctx := context.Background()
	logger := initTestLogger(t)
	ctrl := gomock.NewController(t)

	defer func() {
		if err := logger.Sync(); err != nil {
			t.Logf("failed to sync logger: %v", err)
		}
	}()
	defer ctrl.Finish()

	mockExpenseRepo := mock_repository.NewMockIExpenseRepository(ctrl)
	expenseService := ExpenseService{
		expenseRepository: mockExpenseRepo,
		logger:            logger,
	}

	type mockBehavior func(r *mock_repository.MockIExpenseRepository)

	testTable := []struct {
		name          string
		date          string
		behavior      mockBehavior
		expectedError error
	}{
		{
			name: "valid date",
			date: "2025-07-07",
			behavior: func(r *mock_repository.MockIExpenseRepository) {
				r.EXPECT().DeleteExpensesByDate(ctx, "2025-07-07").Return(nil)
			},
			expectedError: nil,
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			if tt.behavior != nil {
				tt.behavior(mockExpenseRepo)
			}
			err := expenseService.DeleteyExpensesByDate(ctx, tt.date)
			if tt.expectedError == nil {
				assert.NoError(t, err)
			} else {
				assert.Equal(t, err, tt.expectedError)
			}
		})
	}
}
