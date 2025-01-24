package handler

import (
	"cmd/main.go/internal/converter"
	"cmd/main.go/internal/domain"
	"cmd/main.go/internal/transport/telegram/dto"
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"gopkg.in/telebot.v4"
)

type ISpendingService interface {
	AddSpending(ctx context.Context, payload *domain.AddSpending) (int32, error)
	GetWeekSpendings(ctx context.Context, date time.Time) (*domain.WeeklySpendings, error)
}

type SpendingHandler struct {
	spendingService ISpendingService
}

func NewSpendingHandler(spendingService ISpendingService) *SpendingHandler {
	return &SpendingHandler{spendingService: spendingService}
}

func (sh *SpendingHandler) AddSpending(c telebot.Context) error {
	telegramUser := c.Sender()
	if telegramUser.ID != int64(625034947) && telegramUser.ID != int64(481899825) {
		return nil
	}

	input := c.Message().Payload
	if input == "" {
		return c.Send("Please provide the amount to add. Usage: /add <amount>")
	}

	sum, err := strconv.Atoi(strings.TrimSpace(input))
	if err != nil {
		return c.Send("Invalid amount. Please provide a valid integer.")
	}

	payload := &domain.AddSpending{
		Sum:  int32(sum),
		Date: time.Now(),
	}

	daySpendings, err := sh.spendingService.AddSpending(context.Background(), payload)
	if err != nil {
		log.Printf("Error adding spending: %v", err)
		return c.Send("Failed to add spending. Please try again.")
	}

	return c.Send(fmt.Sprintf("Spending added successfully! Overall, spent today: %d", daySpendings))
}

func (sh *SpendingHandler) GetWeekSpendings(c telebot.Context) error {
	telegramUser := c.Sender()
	if telegramUser.ID != int64(625034947) && telegramUser.ID != int64(481899825) {
		return nil
	}

	weekSpendings, err := sh.spendingService.GetWeekSpendings(context.Background(), time.Now())
	if err != nil {
		log.Printf("Error getting week spendings: %v", err)
		return c.Send("Failed to get week spendings. Please try again.")
	}

	response := converter.ToGetWeekSpendingsResponseFromServer(weekSpendings)
	return c.Send(formatWeekSpendingsMessage(response))
}

func formatWeekSpendingsMessage(ws *dto.WeeklySpendings) string {
	weekdays := []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}
	var message strings.Builder
	message.WriteString("Weekly Spendings:\n \n")

	for i, weekday := range weekdays {
		date := time.Now().AddDate(0, 0, -int(time.Now().Weekday())+i+1)
		formattedDate := date.Format("2006-01-02")
		amount := ws.Days[formattedDate]
		message.WriteString(fmt.Sprintf("%s %s: %d din\n", weekday, date.Format("02-01"), amount))
	}

	message.WriteString(fmt.Sprintf("\nTotal: %d din", ws.Total))
	return message.String()
}

// // Helper function to get index of weekday
// func indexOf(slice []string, item string) int {
// 	for i, v := range slice {
// 		if v == item {
// 			return i
// 		}
// 	}
// 	return -1
// }
