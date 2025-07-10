package main

import (
	"context"
	"moneytracker/internal/config"
	"moneytracker/internal/repository"
	"moneytracker/internal/service"
	httpHandler "moneytracker/internal/transport/http/handler"
	"moneytracker/pkg/database"
	"moneytracker/pkg/logger"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	//telegramHandler "cmd/main.go/internal/transport/telegram/handler"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	zapLogger, err := logger.NewLogger(zapcore.InfoLevel, "money_tracker.log")
	if err != nil {
		panic(fmt.Sprintf("error initializing logger: %v", err))
	}

	config, err := config.NewConfig()
	if err != nil {
		zapLogger.Fatal("Error creating config", zap.Error(err))
	}

	err = database.RunMigrations(config.DBSource)
	if err != nil {
		zapLogger.Fatal("Error executing migrations", zap.Error(err))
	}

	db, err := database.NewPGXPool(config.DBSource)
	if err != nil {
		zapLogger.Fatal("Error creating connection pool", zap.Error(err))
	}
	defer db.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	redisDb, err := database.NewRedisConn(ctx, config.RedisAddress, config.RedisPassword)
	if err != nil {
		zapLogger.Fatal("Error creating redis connection", zap.Error(err))
	}
	defer redisDb.Close()

	newRepository := repository.NewRepository(ctx, db, redisDb)
	newService := service.NewService(newRepository)
	newHTTPHandler := httpHandler.NewHTTPHandler(ctx, zapLogger, newService)

	app := startHTTPServer(newHTTPHandler)

	go func() {
		err := app.Listen(":8000")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to start HTTP server: %v\n", err)
			os.Exit(1)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	cancel()
	app.Shutdown()

	zapLogger.Info("Service stopped")
}

func startHTTPServer(handler *httpHandler.HTTPHandler) *fiber.App {
	app := fiber.New()
	//	app.Use(logger)

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "https://myspendingstracker.netlify.app, http://localhost:5173, http://localhost:4173, http://localhost:3004",
		AllowMethods:     "GET,POST,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization,ngrok-skip-browser-warning",
		AllowCredentials: true,
	}))

	handler.SetSpendingRoutes(app)

	return app
}
