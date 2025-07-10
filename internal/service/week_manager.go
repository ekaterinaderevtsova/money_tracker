package service

import (
	"context"
	"fmt"
	"moneytracker/internal/repository"
	"slices"
	"time"

	"go.uber.org/zap"
)

type WeekManager struct {
	currentWeek      []string
	currentWeekRepo  repository.ICurrentSpendingRepository
	archiveWeeksRepo repository.IArchiveSpendingRepository
	logger           *zap.Logger
}

func NewWeekManager(ctx context.Context,
	currentWeekRepo repository.ICurrentSpendingRepository,
	archiveWeeksRepo repository.IArchiveSpendingRepository,
	logger *zap.Logger) (*WeekManager, error) {
	wm := &WeekManager{
		currentWeekRepo:  currentWeekRepo,
		archiveWeeksRepo: archiveWeeksRepo,
		logger:           logger,
	}
	wm.calculateCurrentWeek()

	// Initialize the week automatically
	if err := wm.InitializeWeek(ctx); err != nil {
		logger.Error("failed to initialize new week", zap.Error(err))
		return nil, fmt.Errorf("failed to initialize week: %w", err)
	}

	return wm, nil
}

func (wm *WeekManager) calculateCurrentWeek() {
	now := time.Now()
	weekday := int(now.Weekday())
	if weekday == 0 { // Sunday == 0, convert to 7
		weekday = 7
	}

	startOfWeek := now.AddDate(0, 0, -weekday+1).Truncate(24 * time.Hour)
	wm.currentWeek = make([]string, 7)

	for i := 0; i < 7; i++ {
		day := startOfWeek.AddDate(0, 0, i)
		wm.currentWeek[i] = day.Format("2006-01-02")
	}
}

func (s *WeekManager) IsCurrentWeek(date string) bool {
	return slices.Contains(s.currentWeek, date)
}

func (wm *WeekManager) GetCurrentWeek() []string {
	// Return a copy to prevent external modification
	week := make([]string, len(wm.currentWeek))
	copy(week, wm.currentWeek)
	return week
}

func (wm *WeekManager) InitializeWeek(ctx context.Context) error {
	if err := wm.currentWeekRepo.InitNewWeek(ctx, wm.currentWeek); err != nil {
		return fmt.Errorf("failed to initialize new week in repository: %w", err)
	}
	return nil
}

func (wm *WeekManager) ArchiveCurrentWeek(ctx context.Context) error {
	weekSpending, err := wm.currentWeekRepo.GetWeekSpendings(ctx, wm.currentWeek)
	if err != nil {
		wm.logger.Error("Failed to get week spendings from Redis", zap.Error(err))
		return err
	}

	for _, daySpending := range weekSpending.DaySpendings {
		err := wm.archiveWeeksRepo.AddSpending(ctx, &daySpending)
		if err != nil {
			wm.logger.Error("failed to transfer spending to archive", zap.String("date", daySpending.Day), zap.Int32("sum", daySpending.Sum))
			continue
		}
	}

	wm.logger.Info("Successfully archived current week spendings", zap.String("week start", weekSpending.DaySpendings[0].Day), zap.String("week end", weekSpending.DaySpendings[6].Day))
	return nil
}

func (wm *WeekManager) ResetForNewWeek(ctx context.Context) error {
	if err := wm.currentWeekRepo.FlushAll(ctx); err != nil {
		wm.logger.Error("failed to resent for the new week", zap.Error(err))
		return err
	}
	return nil
}
