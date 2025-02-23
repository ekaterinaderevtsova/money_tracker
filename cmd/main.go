package main

import (
	"cmd/main.go/internal/config"
	"cmd/main.go/internal/repository"
	"cmd/main.go/internal/service"
	httpHandler "cmd/main.go/internal/transport/http/handler"

	//telegramHandler "cmd/main.go/internal/transport/telegram/handler"
	"cmd/main.go/pkg/database"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	fmt.Println("hit main")
	newConfig, err := config.NewConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create config: %v\n", err)
		os.Exit(1)
	}

	// bot, err := telebot.NewBot(telebot.Settings{
	// 	Token: newConfig.TelegramBotToken,
	// })
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "Failed to start Telegram bot: %v\n", err.Error())
	// 	os.Exit(1)
	// }

	err = database.RunMigrations(newConfig.DBSource)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to execute migrations", "err", err)
		os.Exit(1)
	}

	newPool, err := database.NewPGXPool(newConfig.DBSource)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create connection pool", "err", err)
		os.Exit(1)
	}
	defer newPool.Close()

	// mongoClient, err := mongodb.NewClient(newConfig.MongoURI, newConfig.MongoUsername, newConfig.MongoPassword)
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "Failed to connect to MongoDB: %v\n", err)
	// 	os.Exit(1)
	// }

	// defer func() {
	// 	if err := mongoClient.Disconnect(context.Background()); err != nil {
	// 		fmt.Fprintf(os.Stderr, "Error disconnecting MongoDB: %v\n", err)
	// 	}
	// }()

	// database := mongoClient.Database(newConfig.MongoDatabase)

	newRepository := repository.NewRepository(newPool)
	newService := service.NewService(newRepository)
	//	newTelegramHandler := telegramHandler.NewTelegramHandler(newService)
	newHTTPHandler := httpHandler.NewHTTPHandler(newService)

	fmt.Println("created repository, service and handler")

	startHTTPServer(newHTTPHandler)

	// newTelegramHandler.SetCommands(bot)
	// bot.Start()
}

func startHTTPServer(handler *httpHandler.HTTPHandler) {
	app := fiber.New()
	app.Use(logger.New())
	app.Use(cors.New())

	handler.SetSpendingRoutes(app)
	err := app.Listen(":8000")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to start HTTP server: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("server started")
}
