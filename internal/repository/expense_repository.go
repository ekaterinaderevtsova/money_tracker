package repository

import (
	"context"
	"fmt"
	"moneytracker/internal/domain"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ExpenseRepository struct {
	db *pgxpool.Pool
}

func NeExpenseRepository(db *pgxpool.Pool) *ExpenseRepository {
	return &ExpenseRepository{
		db: db,
	}
}

func (sr *ExpenseRepository) AddExpense(ctx context.Context, payload *domain.DailyExpense) error {
	date, err := time.Parse("2006-01-02", payload.Date)
	if err != nil {
		return fmt.Errorf("invalid date format: %w", err)
	}

	_, err = sr.db.Exec(ctx, `
		INSERT INTO spendings (date, sum)
		VALUES ($1, $2)
	`, date, payload.Amount)

	if err != nil {
		return fmt.Errorf("failed to insert/update spending: %w", err)
	}

	return nil
}

func (sr *ExpenseRepository) DeleteExpensesByDate(ctx context.Context, date string) error {
	_, err := sr.db.Exec(ctx, `DELETE FROM spendings WHERE date = $1::date`, date)
	if err != nil {
		return fmt.Errorf("failed to delete spendings: %w", err)
	}
	return nil
}

func (sr *ExpenseRepository) GetWeeklyExpenses(ctx context.Context, week []string) (*domain.WeeklyExpense, error) {
	var weeklyExpenses domain.WeeklyExpense

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
		var dayExpenses domain.DailyExpense
		var date time.Time
		err := rows.Scan(&date, &dayExpenses.Amount)
		if err != nil {
			return nil, err
		}
		dayExpenses.Date = date.Format("02-01-2006")
		weeklyExpenses.DailyExpenses[i] = dayExpenses
		weeklyExpenses.TotalAmount += dayExpenses.Amount
		i++
	}

	weeklyExpenses.AverageDaily = weeklyExpenses.TotalAmount / 7

	return &weeklyExpenses, nil
}

func (sr *ExpenseRepository) GetDailyExpense(ctx context.Context, date string) ([]domain.DailyExpense, error) {
	rows, err := sr.db.Query(ctx, `
		SELECT uuid, sum
		FROM spendings
		WHERE date = $1
		ORDER BY created_at;
		`, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var expenses []domain.DailyExpense
	for rows.Next() {
		var e domain.DailyExpense
		err := rows.Scan(&e.ID, &e.Amount)
		if err != nil {
			return nil, err
		}
		expenses = append(expenses, e)
	}

	return expenses, nil
}

func (sr *ExpenseRepository) UpdateDailyExpense(ctx context.Context, uuid string, newAmount int32) error {
	_, err := sr.db.Exec(ctx, `
		UPDATE spendings SET sum = $1
		WHERE uuid = $2;
	`, newAmount, uuid)

	if err != nil {
		return fmt.Errorf("failed to update spending: %w", err)
	}

	return nil
}

func (sr *ExpenseRepository) DeleteDailyExpense(ctx context.Context, uuid string) error {
	_, err := sr.db.Exec(ctx, `
		DELETE FROM spendings
		WHERE uuid = $1;
	`, uuid)

	if err != nil {
		return fmt.Errorf("failed to delete spending: %w", err)
	}

	return nil
}
