package main

import (
	"cmd/main.go/internal/config"
	"cmd/main.go/internal/repository"
	"cmd/main.go/internal/service"
	"cmd/main.go/internal/transport/telegram/handler"
	mongodb "cmd/main.go/pkg/database"
	"context"
	"fmt"
	"os"

	"gopkg.in/telebot.v4"
)

func main() {
	newConfig, err := config.NewConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create config: %v\n", err)
		os.Exit(1)
	}

	bot, err := telebot.NewBot(telebot.Settings{
		Token: newConfig.TelegramBotToken,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to start Telegram bot: %v\n", err.Error())
		os.Exit(1)
	}

	mongoClient, err := mongodb.NewClient(newConfig.MongoURI, newConfig.MongoUsername, newConfig.MongoPassword)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect to MongoDB: %v\n", err)
		os.Exit(1)
	}

	defer func() {
		if err := mongoClient.Disconnect(context.Background()); err != nil {
			fmt.Fprintf(os.Stderr, "Error disconnecting MongoDB: %v\n", err)
		}
	}()

	database := mongoClient.Database(newConfig.MongoDatabase)

	newRepository := repository.NewRepository(database)
	newService := service.NewService(newRepository)
	newTelegramHandler := handler.NewTelegramHandler(newService)

	newTelegramHandler.SetCommands(bot)
	bot.Start()
}
