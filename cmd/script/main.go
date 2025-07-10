package main

// import (
// 	"context"
// 	"fmt"
// 	"moneytracker/internal/config"
// 	"moneytracker/internal/repository"
// 	"moneytracker/internal/service"
// 	"moneytracker/pkg/database"
// 	"moneytracker/pkg/logger"
// 	"os"

// 	"go.uber.org/zap"
// 	"go.uber.org/zap/zapcore"
// )

// func main() {
// 	zapLogger, err := logger.NewLogger(zapcore.InfoLevel, "money_tracker.log")
// 	if err != nil {
// 		panic(fmt.Sprintf("error initializing logger: %v", err))
// 	}

// 	dbSource := "postgresql://user:password@localhost:5432/spendings_db?sslmode=disable"

// 	config, err := config.NewConfig()
// 	if err != nil {
// 		zapLogger.Fatal("Error creating config", zap.Error(err))
// 	}

// 	db, err := database.NewPGXPool(dbSource)
// 	if err != nil {
// 		zapLogger.Fatal("Error creating connection pool", zap.Error(err))
// 	}
// 	defer db.Close()

// 	ctx, cancel := context.WithCancel(context.Background())
// 	defer cancel()

// 	redisDb, err := database.NewRedisConn(ctx, "localhost:6379", config.RedisPassword)
// 	if err != nil {
// 		zapLogger.Fatal("Error creating redis connection", zap.Error(err))
// 	}
// 	defer redisDb.Close()

// 	repository := repository.NewRepository(ctx, db, redisDb)
// 	service := service.NewService(repository)

// 	if len(os.Args) > 1 {
// 		if os.Args[1] == "—transfer-to-redis" {
// 			err := service.ScriptService.TransferToRedis(ctx)
// 			if err != nil {
// 				zapLogger.Error("Error transferring data to redis", zap.Error(err))
// 			}
// 		}
// 	}
// }
