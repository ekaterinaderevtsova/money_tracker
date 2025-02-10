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
	GetMonthSpendings(ctx context.Context, date time.Time) ([]domain.WeekSpending, error)
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
	weekdays := []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}
	var message strings.Builder
	message.WriteString("Week Spendings:\n \n")

	today := time.Now().UTC()

	// Calculate the start of the week (Monday) for the current week
	startOfWeek := today.AddDate(0, 0, -int(today.Weekday())) // Adjust for Sun
	// if today.Weekday() == 0 {                                   // If today is Sunday, set the start of the week to Monday of the current week
	// 	startOfWeek = today.AddDate(0, 0, -6)
	// }

	for i, weekday := range weekdays {
		date := startOfWeek.AddDate(0, 0, i)
		formattedDate := date.Format("2006-01-02")

		amount := ws.Days[formattedDate]

		message.WriteString(fmt.Sprintf("%s %s: %d din\n", weekday, date.Format("02-01"), amount))
	}

	message.WriteString(fmt.Sprintf("\nSpent: %d din", ws.Total))
	message.WriteString(fmt.Sprintf("\nLeft: %d din", ws.Left))
	return message.String()
}

func (sh *SpendingHandler) GetMonthSpendings(c telebot.Context) error {
	telegramUser := c.Sender()
	if telegramUser.ID != int64(625034947) && telegramUser.ID != int64(481899825) {
		return nil
	}

	monthSpendings, err := sh.spendingService.GetMonthSpendings(context.Background(), time.Now())
	if err != nil {
		log.Printf("Error getting month spendings: %v", err)
		return c.Send("Failed to get month spendings. Please try again.")
	}

	response := converter.ToGetMonthSpendingsResponse(monthSpendings)
	return c.Send(formatMonthSpendingsMessage(response))
}

func formatMonthSpendingsMessage(ws *dto.MonthSpendings) string {
	var message strings.Builder
	message.WriteString("Month Spendings:\n \n")

	for _, week := range ws.Weeks {
		message.WriteString(fmt.Sprintf("Week %d: %d din\n", week.Week, week.Amount))
	}

	message.WriteString(fmt.Sprintf("\nSpent: %d din", ws.Total))
	return message.String()
}
