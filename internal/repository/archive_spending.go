package repository

import (
	"context"
	"fmt"
	"moneytracker/internal/domain"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ArchiveSpendingRepository struct {
	db *pgxpool.Pool
}

func NewArchiveSpendingRepository(db *pgxpool.Pool) *ArchiveSpendingRepository {
	return &ArchiveSpendingRepository{
		db: db,
	}
}

func (sr *ArchiveSpendingRepository) AddSpending(ctx context.Context, payload *domain.DaySpendings) error {
	date, err := time.Parse("2006-01-02", payload.Day)
	if err != nil {
		return fmt.Errorf("invalid date format: %w", err)
	}

	_, err = sr.db.Exec(ctx, `
		INSERT INTO spendings (date, sum)
		VALUES ($1, $2)
		ON CONFLICT (date)
		DO UPDATE SET sum = spendings.sum + EXCLUDED.sum
	`, date, payload.Sum)

	if err != nil {
		return fmt.Errorf("failed to insert/update spending: %w", err)
	}

	return nil
}

func (sr *ArchiveSpendingRepository) GetWeekSpendings(ctx context.Context, week []string) (*domain.WeekSpendings, error) {
	var weekSpendings domain.WeekSpendings
	// date = date.UTC()
	// today := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)

	// daysSinceMonday := (today.Weekday() - time.Monday) % 7
	// if daysSinceMonday < 0 {
	// 	daysSinceMonday += 7 // Ensure non-negative
	// }
	// startOfWeek := today.AddDate(0, 0, -int(daysSinceMonday))

	// endOfWeek := time.Date(
	// 	startOfWeek.Year(), startOfWeek.Month(), startOfWeek.Day()+6,
	// 	23, 59, 59, 0,
	// 	time.UTC,
	// )

	rows, err := sr.db.Query(ctx, `
	   SELECT date_series.date AS date,
       COALESCE(SUM(spendings.sum), 0) AS total
       FROM generate_series($1::date, $2::date, '1 day'::interval) AS date_series
       LEFT JOIN spendings ON date_series.date = spendings.date
       GROUP BY date_series.date
       ORDER BY date_series.date;
		`, week[0], week[6])
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	i := 0
	for rows.Next() {
		var daySpendings domain.DaySpendings
		var date time.Time
		err := rows.Scan(&date, &daySpendings.Sum)
		if err != nil {
			return nil, err
		}
		daySpendings.Day = date.Format("02-01")
		weekSpendings.DaySpendings[i] = daySpendings
		weekSpendings.Total += daySpendings.Sum
		i++
	}

	weekSpendings.Average = weekSpendings.Total / 7

	return &weekSpendings, nil
}
