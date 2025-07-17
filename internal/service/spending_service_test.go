package service

import (
	"context"
	"moneytracker/internal/domain"
	"moneytracker/internal/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func initTestLogger() *zap.Logger {
	logger, _ := zap.NewDevelopment()
	return logger
}

func TestAddSpending(t *testing.T) {
	ctx := context.Background()
	logger := initTestLogger()
	ctrl := gomock.NewController(t)

	defer logger.Sync()
	defer ctrl.Finish()

	currentWeekRepoMock := mocks.NewMockICurrentSpendingRepository(ctrl)
	archiveWeeksRepoMock := mocks.NewMockIArchiveSpendingRepository(ctrl)

	currentWeekRepoMock.EXPECT().InitNewWeek(ctx, gomock.Any()).Return(nil)

	weekManager, err := NewWeekManager(ctx, currentWeekRepoMock, archiveWeeksRepoMock, logger)
	assert.NoError(t, err, "failed to create weekManager")
	spendingService := NewSpendingService(weekManager, currentWeekRepoMock, archiveWeeksRepoMock)

	type mockCurrentBehavior func(r *mocks.MockICurrentSpendingRepository)
	type mockArchiveBehavior func(r *mocks.MockIArchiveSpendingRepository)

	testTable := []struct {
		name            string
		payload         *domain.DaySpendings
		behaviorCurrent mockCurrentBehavior
		behaviorArchive mockArchiveBehavior
		expectedError   error
	}{
		{
			name: "current week valid input",
			payload: &domain.DaySpendings{
				Day: time.Now().Format("2006-01-02"),
				Sum: 1000,
			},
			behaviorCurrent: func(r *mocks.MockICurrentSpendingRepository) {
				r.EXPECT().AddSpending(ctx, &domain.DaySpendings{
					Day: time.Now().Format("2006-01-02"),
					Sum: 1000,
				}).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "old week valid input",
			payload: &domain.DaySpendings{
				Day: "2025-07-07",
				Sum: 1000,
			},
			behaviorArchive: func(r *mocks.MockIArchiveSpendingRepository) {
				r.EXPECT().AddSpending(ctx, &domain.DaySpendings{
					Day: "2025-07-07",
					Sum: 1000,
				}).Return(nil)
			},
			expectedError: nil,
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			if tt.behaviorCurrent != nil {
				tt.behaviorCurrent(currentWeekRepoMock)
			} else {
				tt.behaviorArchive(archiveWeeksRepoMock)
			}
			err := spendingService.AddSpending(ctx, tt.payload)
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
	logger := initTestLogger()
	ctrl := gomock.NewController(t)

	defer logger.Sync()
	defer ctrl.Finish()

	currentWeekRepoMock := mocks.NewMockICurrentSpendingRepository(ctrl)
	archiveWeeksRepoMock := mocks.NewMockIArchiveSpendingRepository(ctrl)

	currentWeekRepoMock.EXPECT().InitNewWeek(ctx, gomock.Any()).Return(nil)

	weekManager, err := NewWeekManager(ctx, currentWeekRepoMock, archiveWeeksRepoMock, logger)
	assert.NoError(t, err, "failed to create weekManager")
	spendingService := NewSpendingService(weekManager, currentWeekRepoMock, archiveWeeksRepoMock)

	type mockCurrentBehavior func(r *mocks.MockICurrentSpendingRepository)
	type mockArchiveBehavior func(r *mocks.MockIArchiveSpendingRepository)

	testTable := []struct {
		name            string
		date            string
		behaviorCurrent mockCurrentBehavior
		behaviorArchive mockArchiveBehavior
		expectedOutput  *domain.WeekSpendings
		expectedError   error
	}{
		{
			name: "current week valid input",
			date: time.Now().Format("2006-01-02"),
			behaviorCurrent: func(r *mocks.MockICurrentSpendingRepository) {
				week := []string{
					time.Now().AddDate(0, 0, -3).Format("2006-01-02"),
					time.Now().AddDate(0, 0, -2).Format("2006-01-02"),
					time.Now().AddDate(0, 0, -1).Format("2006-01-02"),
					time.Now().Format("2006-01-02"),
					time.Now().AddDate(0, 0, 1).Format("2006-01-02"),
					time.Now().AddDate(0, 0, 2).Format("2006-01-02"),
					time.Now().AddDate(0, 0, 3).Format("2006-01-02"),
				}

				r.EXPECT().GetWeekSpendings(ctx, week).Return(&domain.WeekSpendings{
					DaySpendings: [7]domain.DaySpendings{
						{
							Day: time.Now().Format("2006-01-02"),
							Sum: 1000,
						},
					},
					Total:   1000,
					Average: 1000,
				}, nil)
			},
			expectedOutput: &domain.WeekSpendings{
				DaySpendings: [7]domain.DaySpendings{
					{
						Day: time.Now().Format("2006-01-02"),
						Sum: 1000,
					},
				},
				Total:   1000,
				Average: 1000,
			},
			expectedError: nil,
		},
		{
			name: "archive week valid input",
			date: "2025-07-07",
			behaviorArchive: func(r *mocks.MockIArchiveSpendingRepository) {
				week := []string{
					"2025-07-07",
					"2025-07-08",
					"2025-07-09",
					"2025-07-10",
					"2025-07-11",
					"2025-07-12",
					"2025-07-13",
				}
				r.EXPECT().GetWeekSpendings(ctx, week).Return(&domain.WeekSpendings{
					DaySpendings: [7]domain.DaySpendings{
						{
							Day: "2025-07-07",
							Sum: 1000,
						},
					},
					Total:   1000,
					Average: 1000,
				}, nil)
			},
			expectedOutput: &domain.WeekSpendings{
				DaySpendings: [7]domain.DaySpendings{
					{
						Day: "2025-07-07",
						Sum: 1000,
					},
				},
				Total:   1000,
				Average: 1000,
			},
			expectedError: nil,
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			if tt.behaviorCurrent != nil {
				tt.behaviorCurrent(currentWeekRepoMock)
			} else {
				tt.behaviorArchive(archiveWeeksRepoMock)
			}
			receivedOutput, err := spendingService.GetWeekSpendings(ctx, tt.date)
			assert.Equal(t, tt.expectedOutput, receivedOutput)
			if tt.expectedError == nil {
				assert.NoError(t, err)
			} else {
				assert.Equal(t, err, tt.expectedError)
			}
		})
	}
}
