package repository

import (
	"cmd/main.go/internal/domain"
	"context"
	"fmt"
	"slices"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type SpendingRepository struct {
	db          *pgxpool.Pool
	redisDb     *redis.Client
	currentWeek []string
}

func NewSpendingRepository(ctx context.Context, db *pgxpool.Pool, redisDb *redis.Client) *SpendingRepository {
	spendingRepo := &SpendingRepository{
		db:      db,
		redisDb: redisDb,
	}

	spendingRepo.initCurrentWeek()

	// if err := migrations.RunRedisMigrations(ctx, redisDb, spendingRepo.currentWeek); err != nil {
	// 	fmt.Printf("failed to run redis migrations: %v\n", err)
	// }

	err := spendingRepo.initRedisWeekDays(ctx)
	if err != nil {
		fmt.Printf("failed to initialize redis data: %v\n", err)
	}

	go func() {
		// TODO: Create ticker
		for {
			now := time.Now()
			daysUntilMonday := (int(time.Monday) - int(now.Weekday()) + 7) % 7
			if daysUntilMonday == 0 && now.Hour() >= 0 {
				daysUntilMonday = 7
			}

			nextMonday := now.Truncate(24*time.Hour).AddDate(0, 0, daysUntilMonday)
			nextRun := time.Date(
				nextMonday.Year(), nextMonday.Month(), nextMonday.Day(),
				0, 0, 0, 0, nextMonday.Location(),
			)

			sleepDuration := time.Until(nextRun)
			timer := time.NewTimer(sleepDuration)

			select {
			case <-ctx.Done():
				timer.Stop()
				return
			case <-timer.C:
				err := spendingRepo.transferRedisDataToPostgres(ctx)
				if err != nil {
					// TODO: log
				}

				spendingRepo.initCurrentWeek()
				spendingRepo.initRedisWeekDays(ctx)
			}
		}
	}()

	return spendingRepo
}

func (r *SpendingRepository) initCurrentWeek() {
	now := time.Now()
	weekday := int(now.Weekday())
	if weekday == 0 { // Sunday == 0
		weekday = 7
	}

	startOfWeek := now.AddDate(0, 0, -weekday+1).Truncate(24 * time.Hour)
	r.currentWeek = make([]string, 7)

	for i := 0; i < 7; i++ {
		day := startOfWeek.AddDate(0, 0, i)
		r.currentWeek[i] = day.Format("2006-01-02")
	}
}

func (r *SpendingRepository) initRedisWeekDays(ctx context.Context) error {
	log.Print("Init redis")
	for _, day := range r.currentWeek {
		dayKey := domain.SpendingsKey + day
		totalKey := domain.TotalKey + day

		exists, err := r.redisDb.Exists(ctx, dayKey).Result()
		if err != nil {
			return err
		}
		if exists == 0 {
			if err := r.redisDb.RPush(ctx, dayKey, 0).Err(); err != nil {
				return err
			}
		}

		exists, err = r.redisDb.Exists(ctx, totalKey).Result()
		if err != nil {
			return err
		}
		if exists == 0 {
			if err := r.redisDb.Set(ctx, totalKey, 0, 0).Err(); err != nil {
				return err
			}
		}
	}
	return nil
}

func (sr *SpendingRepository) transferRedisDataToPostgres(ctx context.Context) error {
	for _, day := range sr.currentWeek {
		// TODO: use MGET or scan
		sumStr, err := sr.redisDb.Get(ctx, domain.TotalKey+day).Result()
		if err != nil {
			// TODO: log
			continue
		}
		// TODO: process error
		sum, _ := strconv.Atoi(sumStr)

		spendingInfo := domain.DaySpendings{
			Day: day,
			Sum: int32(sum),
		}

		// TODO: write AddSpendings func
		err = sr.AddSpending(ctx, &spendingInfo)
		if err != nil {
			// TODO: log
			continue
		}
	}

	err := sr.redisDb.FlushAll(ctx)
	if err != nil {
		// TODO: log
	}

	return nil
}

func (r *SpendingRepository) IsCurrentWeek(date string) bool {
	return slices.Contains(r.currentWeek, date)
}

func (sr *SpendingRepository) AddCurrentWeekSpending(ctx context.Context, payload *domain.DaySpendings) error {

	if err := sr.redisDb.LPush(ctx, domain.SpendingsKey+payload.Day, payload.Sum).Err(); err != nil {
		return fmt.Errorf("failed to set new sum in redis: %w", err)
	}

	if err := sr.redisDb.IncrBy(ctx, domain.TotalKey+payload.Day, int64(payload.Sum)).Err(); err != nil {
		return fmt.Errorf("failed to set new sum in redis: %w", err)
	}

	return nil
}

func (sr *SpendingRepository) GetCurrentWeekSpendings(ctx context.Context) (*domain.WeekSpendings, error) {
	var weekSpendings domain.WeekSpendings
	var total int32
	var daysCount int32
	today := time.Now().Format("2006-01-02")

	for i, day := range sr.currentWeek {
		//keys = append(keys, string(day))
		sumStr, err := sr.redisDb.Get(ctx, domain.TotalKey+day).Result()
		if err != nil {
			// TODO: log
			continue
		}

		sum, _ := strconv.Atoi(sumStr)

		weekSpendings.DaySpendings[i] = domain.DaySpendings{
			Day: day,
			Sum: int32(sum),
		}

		total += int32(sum)

		if today >= day {
			daysCount++
		}
	}

	weekSpendings.Total = total
	weekSpendings.Average = 0
	if daysCount > 0 {
		weekSpendings.Average = total / daysCount
	}
	return &weekSpendings, nil
}

// TODO: transfer to a separate repo
func (sr *SpendingRepository) AddSpending(ctx context.Context, payload *domain.DaySpendings) error {
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

func (sr *SpendingRepository) GetWeekSpendings(ctx context.Context, date time.Time) (*domain.WeekSpendings, error) {
	var weekSpendings domain.WeekSpendings
	date = date.UTC()
	today := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)

	daysSinceMonday := (today.Weekday() - time.Monday) % 7
	if daysSinceMonday < 0 {
		daysSinceMonday += 7 // Ensure non-negative
	}
	startOfWeek := today.AddDate(0, 0, -int(daysSinceMonday))

	endOfWeek := time.Date(
		startOfWeek.Year(), startOfWeek.Month(), startOfWeek.Day()+6,
		23, 59, 59, 0,
		time.UTC,
	)

	rows, err := sr.db.Query(ctx, `
	   SELECT date_series.date AS date,
       COALESCE(SUM(spendings.sum), 0) AS total
       FROM generate_series($1::date, $2::date, '1 day'::interval) AS date_series
       LEFT JOIN spendings ON date_series.date = spendings.date
       GROUP BY date_series.date
       ORDER BY date_series.date;
		`, startOfWeek.Format("2006-01-02"), endOfWeek.Format("2006-01-02"))
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
