package service

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
)

type WeeklyScheduler struct {
	weekManager IWeekManager
	stopCh      chan struct{}
	logger      *zap.Logger
}

func NewWeeklyScheduler(weekManager IWeekManager, logger *zap.Logger) *WeeklyScheduler {
	return &WeeklyScheduler{
		weekManager: weekManager,
		stopCh:      make(chan struct{}),
		logger:      logger,
	}
}

func (ws *WeeklyScheduler) Start(ctx context.Context) error {
	go func() {
		ws.logger.Info("Weekly scheduler started")

		for {
			nextRun := ws.calculateNextMondayMidnight()
			sleepDuration := time.Until(nextRun)
			timer := time.NewTimer(sleepDuration)

			select {
			case <-ctx.Done():
				timer.Stop()
				ws.logger.Info("Weekly scheduler stopped due to context cancellation")
				return
			case <-timer.C:
				if err := ws.performWeeklyArchive(ctx); err != nil {
					ws.logger.Error("Weekly archive failed", zap.Error(err))
				} else {
					ws.logger.Info("Weekly archive completed successfully")
				}
			}
		}
	}()

	return nil
}

func (ws *WeeklyScheduler) calculateNextMondayMidnight() time.Time {
	now := time.Now()
	daysUntilMonday := (int(time.Monday) - int(now.Weekday()) + 7) % 7
	if daysUntilMonday == 0 && now.Hour() >= 0 {
		daysUntilMonday = 7
	}

	nextMonday := now.Truncate(24*time.Hour).AddDate(0, 0, daysUntilMonday)
	return time.Date(
		nextMonday.Year(), nextMonday.Month(), nextMonday.Day(),
		0, 0, 0, 0, nextMonday.Location(),
	)
}

func (ws *WeeklyScheduler) performWeeklyArchive(ctx context.Context) error {
	// Archive current week
	if err := ws.weekManager.ArchiveCurrentWeek(ctx); err != nil {
		return fmt.Errorf("failed to archive current week: %w", err)
	}

	// Reset for new week
	if err := ws.weekManager.ResetForNewWeek(ctx); err != nil {
		return fmt.Errorf("failed to reset for new week: %w", err)
	}

	return nil
}
