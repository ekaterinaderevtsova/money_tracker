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
	GetWeekSpendings(ctx context.Context, date time.Time) (*domain.WeekSpendings, error)
	GetMonthSpendings(ctx context.Context, date time.Time) ([]domain.WeekTotalSpending, error)
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

func (sh *SpendingHandler) AddOldSpending(c telebot.Context) error {
	telegramUser := c.Sender()
	if telegramUser.ID != int64(625034947) && telegramUser.ID != int64(481899825) {
		return nil
	}

	input := c.Message().Payload
	if input == "" {
		return c.Send("Please provide the amount to add. Usage: /add <amount>")
	}

	parts := strings.Fields(input)
	if len(parts) != 2 {
		return c.Send("Invalid format. Please provide both date and sum as '/add <date> <sum>'")
	}

	dateStr := parts[0]
	sumStr := parts[1]

	sum, err := strconv.Atoi(sumStr)
	if err != nil {
		return c.Send("Invalid sum. Please provide a valid integer.")
	}

	parsedDate, err := time.Parse("02.01", dateStr)
	if err != nil {
		return c.Send("Invalid date format. Use DD.MM format like 24.06")
	}

	if parsedDate.Format("02.01") != dateStr {
		return c.Send("Invalid date. Please check the day and month values.")
	}

	payload := &domain.AddSpending{
		Sum: int32(sum),
		Date: time.Date(
			2025,
			parsedDate.Month(),
			parsedDate.Day(),
			0, 0, 0, 0, // Set time to midnight
			time.Local,
		),
	}

	daySpendings, err := sh.spendingService.AddSpending(context.Background(), payload)
	if err != nil {
		log.Printf("Error adding spending: %v", err)
		return c.Send("Failed to add spending. Please try again.")
	}

	return c.Send(fmt.Sprintf("Spending added successfully! Overall, spent on %s: %d", dateStr, daySpendings))
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
	message.WriteString("Week Spendings:\n\n")

	for i, weekday := range weekdays {
		daySpending := ws.Days[i] // Access DaySpendings struct from array
		message.WriteString(fmt.Sprintf("%s %s: %d din\n", weekday, daySpending.Day, daySpending.Sum))
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
