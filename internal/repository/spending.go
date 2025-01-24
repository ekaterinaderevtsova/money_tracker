package repository

import (
	"cmd/main.go/internal/domain"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type SpendingRepository struct {
	db *mongo.Collection
}

func NewSpendingRepository(db *mongo.Database) *SpendingRepository {
	return &SpendingRepository{db: db.Collection("spendings")}
}

func (sr *SpendingRepository) AddSpending(ctx context.Context, payload *domain.AddSpending) error {
	_, err := sr.db.InsertOne(ctx, payload)
	if err != nil {
		return err
	}

	// insertedID, ok := result.InsertedID.(primitive.ObjectID)
	// if !ok {
	// 	return nil, fmt.Errorf("failed to cast inserted ID to ObjectID")
	// }

	// var spendingInfo domain.AddSpending
	// err = sr.db.FindOne(ctx, bson.M{"_id": insertedID}).Decode(&spendingInfo)
	// if err != nil {
	// 	return nil, err
	// }

	return nil
}

func (sr *SpendingRepository) GetAllSpendings(ctx context.Context) ([]domain.AddSpending, error) {
	cursor, err := sr.db.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch records: %w", err)
	}
	defer cursor.Close(ctx)

	var results []domain.AddSpending
	if err := cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("failed to decode records: %w", err)
	}

	return results, nil
}

func (sr *SpendingRepository) GetDaySpendings(ctx context.Context, date time.Time) (int32, error) {
	date = date.UTC()
	startOfDay := date.Truncate(24 * time.Hour)
	endOfDay := startOfDay.Add(24 * time.Hour)

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"date": bson.M{
					"$gte": startOfDay,
					"$lt":  endOfDay,
				},
			},
		},
		{
			"$group": bson.M{
				"_id": nil,
				"total": bson.M{
					"$sum": "$sum",
				},
			},
		},
	}

	cursor, err := sr.db.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, fmt.Errorf("failed to run aggregation: %w", err)
	}
	defer cursor.Close(ctx)

	var result struct {
		Total int32 `bson:"total"`
	}

	if !cursor.Next(ctx) {
		// No records found
		return 0, nil
	}

	if err := cursor.Decode(&result); err != nil {
		return 0, fmt.Errorf("failed to decode result: %w", err)
	}

	return result.Total, nil
}

func (sr *SpendingRepository) GetWeekSpendings(ctx context.Context, date time.Time) (*domain.WeeklySpendings, error) {
	date = date.UTC()
	startOfWeek := date.AddDate(0, 0, -int(date.Weekday()))
	startOfWeek = startOfWeek.Truncate(24 * time.Hour)
	endOfWeek := startOfWeek.AddDate(0, 0, 7)

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"date": bson.M{
					"$gte": startOfWeek,
					"$lt":  endOfWeek,
				},
			},
		},
		{
			"$group": bson.M{
				"_id": bson.M{
					"$dateToString": bson.M{
						"format": "%Y-%m-%d",
						"date":   "$date",
					},
				},
				"total": bson.M{
					"$sum": "$sum",
				},
			},
		},
		{
			"$sort": bson.M{
				"_id": 1,
			},
		},
	}

	cursor, err := sr.db.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("failed to run aggregation: %w", err)
	}
	defer cursor.Close(ctx)

	weekSpendings := make(map[string]int32)
	var result struct {
		ID    string `bson:"_id"`
		Total int32  `bson:"total"`
	}

	for cursor.Next(ctx) {
		if err := cursor.Decode(&result); err != nil {
			return nil, fmt.Errorf("failed to decode result: %w", err)
		}
		weekSpendings[result.ID] = result.Total
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return &domain.WeeklySpendings{
		Days:  weekSpendings,
		Total: calculateTotal(weekSpendings),
	}, nil
}

func calculateTotal(days map[string]int32) int32 {
	var total int32
	for _, dayTotal := range days {
		total += dayTotal
	}
	return total
}
