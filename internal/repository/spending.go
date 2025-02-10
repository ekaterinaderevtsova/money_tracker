package repository

import (
	"cmd/main.go/internal/domain"
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SpendingRepository struct {
	db *pgxpool.Pool
}

func NewSpendingRepository(db *pgxpool.Pool) *SpendingRepository {
	return &SpendingRepository{db: db}
}

func (sr *SpendingRepository) AddSpending(ctx context.Context, payload *domain.AddSpending) error {
	_, err := sr.db.Exec(ctx, `
		INSERT INTO spendings (date, sum)
		VALUES ($1, $2);
	`, payload.Date, payload.Sum)
	if err != nil {
		return err
	}
	return nil
}

func (sr *SpendingRepository) GetDaySpendings(ctx context.Context, date time.Time) (int32, error) {
	var total int32
	year := date.Year()
	month := date.Month()
	day := date.Day()

	err := sr.db.QueryRow(ctx, `
		SELECT SUM(sum)
		FROM spendings
		WHERE EXTRACT(YEAR FROM date) = $1
		AND EXTRACT(MONTH FROM date) = $2
		AND EXTRACT(DAY FROM date) = $3;
		`, year, month, day).Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (sr *SpendingRepository) GetWeekSpendings(ctx context.Context, date time.Time) (*domain.WeeklySpendings, error) {
	return nil, nil
}

func (sr *SpendingRepository) GetMonthSpendings(ctx context.Context, date time.Time) ([]domain.WeekSpending, error) {
	return nil, nil
}
