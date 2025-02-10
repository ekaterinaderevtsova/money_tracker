package repository

// type SpendingRepository struct {
// 	//db *mongo.Collection
// 	db *pgxpool.Pool
// }

// func NewSpendingRepository(db *pgxpool.Pool) *SpendingRepository {
// 	return &SpendingRepository{db: db}
// }

// func (sr *SpendingRepository) AddSpending(ctx context.Context, payload *domain.AddSpending) error {
// 	// _, err := sr.db.InsertOne(ctx, payload)
// 	// if err != nil {
// 	// 	return err
// 	// }

// 	return nil
// }

// func (sr *SpendingRepository) GetAllSpendings(ctx context.Context) ([]domain.AddSpending, error) {
// 	// cursor, err := sr.db.Find(ctx, bson.M{})
// 	// if err != nil {
// 	// 	return nil, fmt.Errorf("failed to fetch records: %w", err)
// 	// }
// 	// defer cursor.Close(ctx)

// 	// var results []domain.AddSpending
// 	// if err := cursor.All(ctx, &results); err != nil {
// 	// 	return nil, fmt.Errorf("failed to decode records: %w", err)
// 	// }

// 	return nil, nil
// }

// func (sr *SpendingRepository) GetDaySpendings(ctx context.Context, date time.Time) (int32, error) {
// 	// 	date = date.UTC()
// 	// 	startOfDay := date.Truncate(24 * time.Hour)
// 	// 	endOfDay := startOfDay.Add(24 * time.Hour)

// 	// 	pipeline := []bson.M{
// 	// 		{
// 	// 			"$match": bson.M{
// 	// 				"date": bson.M{
// 	// 					"$gte": startOfDay,
// 	// 					"$lt":  endOfDay,
// 	// 				},
// 	// 			},
// 	// 		},
// 	// 		{
// 	// 			"$group": bson.M{
// 	// 				"_id": nil,
// 	// 				"total": bson.M{
// 	// 					"$sum": "$sum",
// 	// 				},
// 	// 			},
// 	// 		},
// 	// 	}

// 	// 	cursor, err := sr.db.Aggregate(ctx, pipeline)
// 	// 	if err != nil {
// 	// 		return 0, fmt.Errorf("failed to run aggregation: %w", err)
// 	// 	}
// 	// 	defer cursor.Close(ctx)

// 	// 	var result struct {
// 	// 		Total int32 `bson:"total"`
// 	// 	}

// 	// 	if !cursor.Next(ctx) {
// 	// 		// No records found
// 	// 		return 0, nil
// 	// 	}

// 	// 	if err := cursor.Decode(&result); err != nil {
// 	// 		return 0, fmt.Errorf("failed to decode result: %w", err)
// 	// 	}

// 	return 0, nil
// }

// func (sr *SpendingRepository) GetWeekSpendings(ctx context.Context, date time.Time) (*domain.WeeklySpendings, error) {
// 	// 	date = date.UTC()
// 	// 	today := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)

// 	// 	startOfWeek := today.AddDate(0, 0, -int(today.Weekday()))

// 	// 	endOfWeek := time.Date(
// 	// 		startOfWeek.Year(), startOfWeek.Month(), startOfWeek.Day()+6,
// 	// 		23, 59, 59, 999000000,
// 	// 		time.UTC,
// 	// 	)

// 	// 	pipeline := []bson.M{
// 	// 		{
// 	// 			"$match": bson.M{
// 	// 				"date": bson.M{
// 	// 					"$gte": startOfWeek,
// 	// 					"$lte": endOfWeek,
// 	// 				},
// 	// 			},
// 	// 		},
// 	// 		{
// 	// 			"$group": bson.M{
// 	// 				"_id": bson.M{
// 	// 					"$dateToString": bson.M{
// 	// 						"format": "%Y-%m-%d",
// 	// 						"date":   "$date",
// 	// 					},
// 	// 				},
// 	// 				"total": bson.M{
// 	// 					"$sum": "$sum",
// 	// 				},
// 	// 			},
// 	// 		},
// 	// 		{
// 	// 			"$sort": bson.M{
// 	// 				"_id": 1,
// 	// 			},
// 	// 		},
// 	// 	}

// 	// 	cursor, err := sr.db.Aggregate(ctx, pipeline)
// 	// 	if err != nil {
// 	// 		return nil, fmt.Errorf("failed to run aggregation: %w", err)
// 	// 	}
// 	// 	defer cursor.Close(ctx)

// 	// 	weekSpendings := make(map[string]int32)
// 	// 	var result struct {
// 	// 		ID    string `bson:"_id"`
// 	// 		Total int32  `bson:"total"`
// 	// 	}

// 	// 	for cursor.Next(ctx) {
// 	// 		if err := cursor.Decode(&result); err != nil {
// 	// 			return nil, fmt.Errorf("failed to decode result: %w", err)
// 	// 		}
// 	// 		weekSpendings[result.ID] = result.Total
// 	// 	}

// 	// 	if err := cursor.Err(); err != nil {
// 	// 		return nil, fmt.Errorf("cursor error: %w", err)
// 	// 	}

// 	return nil, nil
// }

// func calculateTotal(days map[string]int32) int32 {
// 	var total int32
// 	for _, dayTotal := range days {
// 		total += dayTotal
// 	}
// 	return total
// }

// func (sr *SpendingRepository) GetMonthSpendings(ctx context.Context, date time.Time) ([]domain.WeekSpending, error) {
// 	// date = date.UTC()
// 	// today := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
// 	// startOfMonth := today.AddDate(0, 0, -int(today.Day())+1)
// 	// endOfMonth := time.Date(
// 	// 	startOfMonth.Year(), startOfMonth.Month(), getMonthDuration(startOfMonth.Month(), startOfMonth.Year()),
// 	// 	23, 59, 59, 999000000, // 999000000 nanoseconds = 999 milliseconds
// 	// 	time.UTC,
// 	// )
// 	// fmt.Println(startOfMonth, endOfMonth)

// 	// pipeline := []bson.M{
// 	// 	{
// 	// 		"$match": bson.M{
// 	// 			"date": bson.M{
// 	// 				"$gte": startOfMonth,
// 	// 				"$lte": endOfMonth,
// 	// 			},
// 	// 		},
// 	// 	},
// 	// 	// Change grouping to week
// 	// 	{
// 	// 		"$group": bson.M{
// 	// 			"_id": bson.M{
// 	// 				"$dateToString": bson.M{
// 	// 					"format": "%Y-%m-%d",
// 	// 					"date":   "$date",
// 	// 				},
// 	// 			},
// 	// 			"total": bson.M{
// 	// 				"$sum": "$sum",
// 	// 			},
// 	// 		},
// 	// 	},
// 	// 	{
// 	// 		"$sort": bson.M{
// 	// 			"_id": 1,
// 	// 		},
// 	// 	},
// 	// }

// 	// cursor, err := sr.db.Aggregate(ctx, pipeline)
// 	// if err != nil {
// 	// 	return nil, fmt.Errorf("failed to run aggregation: %w", err)
// 	// }
// 	// defer cursor.Close(ctx)

// 	// weeklySpendingsMap := make(map[int]int32)

// 	// var result struct {
// 	// 	ID    string `bson:"_id"`
// 	// 	Total int32  `bson:"total"`
// 	// }

// 	// for cursor.Next(ctx) {
// 	// 	if err := cursor.Decode(&result); err != nil {
// 	// 		return nil, fmt.Errorf("failed to decode result: %w", err)
// 	// 	}

// 	// 	// Parse the date from result.ID (format: "YYYY-MM-DD")
// 	// 	spendingDate, err := time.Parse("2006-01-02", result.ID)
// 	// 	if err != nil {
// 	// 		return nil, fmt.Errorf("failed to parse date: %w", err)
// 	// 	}

// 	// 	// Calculate week number based on the first day of the month
// 	// 	firstDay := time.Date(spendingDate.Year(), spendingDate.Month(), 1, 0, 0, 0, 0, time.UTC)
// 	// 	_, firstWeek := firstDay.ISOWeek()
// 	// 	_, spendingWeek := spendingDate.ISOWeek()
// 	// 	monthWeek := spendingWeek - firstWeek + 1

// 	// 	// Add spending to corresponding week
// 	// 	weeklySpendingsMap[monthWeek] += result.Total
// 	// }

// 	// if err := cursor.Err(); err != nil {
// 	// 	return nil, fmt.Errorf("cursor error: %w", err)
// 	// }

// 	// // Convert map to slice of WeekSpending
// 	// weekSpendings := make([]domain.WeekSpending, 0, len(weeklySpendingsMap))
// 	// for weekNum, amount := range weeklySpendingsMap {
// 	// 	weekSpendings = append(weekSpendings, domain.WeekSpending{
// 	// 		Week:   weekNum,
// 	// 		Amount: amount,
// 	// 	})
// 	// }

// 	// // Sort weeks in ascending order
// 	// sort.Slice(weekSpendings, func(i, j int) bool {
// 	// 	return weekSpendings[i].Week < weekSpendings[j].Week
// 	// })
// 	// date = date.UTC()
// 	// today := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
// 	// startOfMonth := today.AddDate(0, 0, -int(today.Day())+1)
// 	// endOfMonth := time.Date(
// 	// 	startOfMonth.Year(), startOfMonth.Month(), getMonthDuration(startOfMonth.Month(), startOfMonth.Year()),
// 	// 	23, 59, 59, 999000000, // 999000000 nanoseconds = 999 milliseconds
// 	// 	time.UTC,
// 	// )
// 	// fmt.Println(startOfMonth, endOfMonth)

// 	// pipeline := []bson.M{
// 	// 	{
// 	// 		"$match": bson.M{
// 	// 			"date": bson.M{
// 	// 				"$gte": startOfMonth,
// 	// 				"$lte": endOfMonth,
// 	// 			},
// 	// 		},
// 	// 	},
// 	// 	// Change grouping to week
// 	// 	{
// 	// 		"$group": bson.M{
// 	// 			"_id": bson.M{
// 	// 				"$dateToString": bson.M{
// 	// 					"format": "%Y-%m-%d",
// 	// 					"date":   "$date",
// 	// 				},
// 	// 			},
// 	// 			"total": bson.M{
// 	// 				"$sum": "$sum",
// 	// 			},
// 	// 		},
// 	// 	},
// 	// 	{
// 	// 		"$sort": bson.M{
// 	// 			"_id": 1,
// 	// 		},
// 	// 	},
// 	// }

// 	// cursor, err := sr.db.Aggregate(ctx, pipeline)
// 	// if err != nil {
// 	// 	return nil, fmt.Errorf("failed to run aggregation: %w", err)
// 	// }
// 	// defer cursor.Close(ctx)

// 	// weeklySpendingsMap := make(map[int]int32)

// 	// var result struct {
// 	// 	ID    string `bson:"_id"`
// 	// 	Total int32  `bson:"total"`
// 	// }

// 	// for cursor.Next(ctx) {
// 	// 	if err := cursor.Decode(&result); err != nil {
// 	// 		return nil, fmt.Errorf("failed to decode result: %w", err)
// 	// 	}

// 	// 	// Parse the date from result.ID (format: "YYYY-MM-DD")
// 	// 	spendingDate, err := time.Parse("2006-01-02", result.ID)
// 	// 	if err != nil {
// 	// 		return nil, fmt.Errorf("failed to parse date: %w", err)
// 	// 	}

// 	// 	// Calculate week number based on the first day of the month
// 	// 	firstDay := time.Date(spendingDate.Year(), spendingDate.Month(), 1, 0, 0, 0, 0, time.UTC)
// 	// 	_, firstWeek := firstDay.ISOWeek()
// 	// 	_, spendingWeek := spendingDate.ISOWeek()
// 	// 	monthWeek := spendingWeek - firstWeek + 1

// 	// 	// Add spending to corresponding week
// 	// 	weeklySpendingsMap[monthWeek] += result.Total
// 	// }

// 	// if err := cursor.Err(); err != nil {
// 	// 	return nil, fmt.Errorf("cursor error: %w", err)
// 	// }

// 	// // Convert map to slice of WeekSpending
// 	// weekSpendings := make([]domain.WeekSpending, 0, len(weeklySpendingsMap))
// 	// for weekNum, amount := range weeklySpendingsMap {
// 	// 	weekSpendings = append(weekSpendings, domain.WeekSpending{
// 	// 		Week:   weekNum,
// 	// 		Amount: amount,
// 	// 	})
// 	// }

// 	// // Sort weeks in ascending order
// 	// sort.Slice(weekSpendings, func(i, j int) bool {
// 	// 	return weekSpendings[i].Week < weekSpendings[j].Week
// 	//	})
// 	// date = date.UTC()
// 	// today := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
// 	// startOfMonth := today.AddDate(0, 0, -int(today.Day())+1)
// 	// endOfMonth := time.Date(
// 	// 	startOfMonth.Year(), startOfMonth.Month(), getMonthDuration(startOfMonth.Month(), startOfMonth.Year()),
// 	// 	23, 59, 59, 999000000, // 999000000 nanoseconds = 999 milliseconds
// 	// 	time.UTC,
// 	// )
// 	// fmt.Println(startOfMonth, endOfMonth)

// 	// pipeline := []bson.M{
// 	// 	{
// 	// 		"$match": bson.M{
// 	// 			"date": bson.M{
// 	// 				"$gte": startOfMonth,
// 	// 				"$lte": endOfMonth,
// 	// 			},
// 	// 		},
// 	// 	},
// 	// 	// Change grouping to week
// 	// 	{
// 	// 		"$group": bson.M{
// 	// 			"_id": bson.M{
// 	// 				"$dateToString": bson.M{
// 	// 					"format": "%Y-%m-%d",
// 	// 					"date":   "$date",
// 	// 				},
// 	// 			},
// 	// 			"total": bson.M{
// 	// 				"$sum": "$sum",
// 	// 			},
// 	// 		},
// 	// 	},
// 	// 	{
// 	// 		"$sort": bson.M{
// 	// 			"_id": 1,
// 	// 		},
// 	// 	},
// 	// }

// 	// cursor, err := sr.db.Aggregate(ctx, pipeline)
// 	// if err != nil {
// 	// 	return nil, fmt.Errorf("failed to run aggregation: %w", err)
// 	// }
// 	// defer cursor.Close(ctx)

// 	// weeklySpendingsMap := make(map[int]int32)

// 	// var result struct {
// 	// 	ID    string `bson:"_id"`
// 	// 	Total int32  `bson:"total"`
// 	// }

// 	// for cursor.Next(ctx) {
// 	// 	if err := cursor.Decode(&result); err != nil {
// 	// 		return nil, fmt.Errorf("failed to decode result: %w", err)
// 	// 	}

// 	// 	// Parse the date from result.ID (format: "YYYY-MM-DD")
// 	// 	spendingDate, err := time.Parse("2006-01-02", result.ID)
// 	// 	if err != nil {
// 	// 		return nil, fmt.Errorf("failed to parse date: %w", err)
// 	// 	}

// 	// 	// Calculate week number based on the first day of the month
// 	// 	firstDay := time.Date(spendingDate.Year(), spendingDate.Month(), 1, 0, 0, 0, 0, time.UTC)
// 	// 	_, firstWeek := firstDay.ISOWeek()
// 	// 	_, spendingWeek := spendingDate.ISOWeek()
// 	// 	monthWeek := spendingWeek - firstWeek + 1

// 	// 	// Add spending to corresponding week
// 	// 	weeklySpendingsMap[monthWeek] += result.Total
// 	// }

// 	// if err := cursor.Err(); err != nil {
// 	// 	return nil, fmt.Errorf("cursor error: %w", err)
// 	// }

// 	// // Convert map to slice of WeekSpending
// 	// weekSpendings := make([]domain.WeekSpending, 0, len(weeklySpendingsMap))
// 	// for weekNum, amount := range weeklySpendingsMap {
// 	// 	weekSpendings = append(weekSpendings, domain.WeekSpending{
// 	// 		Week:   weekNum,
// 	// 		Amount: amount,
// 	// 	})
// 	// }

// 	// // Sort weeks in ascending order
// 	// sort.Slice(weekSpendings, func(i, j int) bool {
// 	// 	return weekSpendings[i].Week < weekSpendings[j].Week
// 	// })

// 	// fmt.Println(weekSpendings)
// 	return nil, nil
// }

// func getMonthDuration(month time.Month, year int) int {
// 	nextMonth := time.Date(year, month+1, 1, 0, 0, 0, 0, time.UTC)
// 	lastDay := nextMonth.AddDate(0, 0, -1)

// 	return lastDay.Day()
// }
