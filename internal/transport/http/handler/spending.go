package handler

import (
	"cmd/main.go/internal/converter"
	"cmd/main.go/internal/domain"
	httpdto "cmd/main.go/internal/transport/http/dto"
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
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

func (sh *SpendingHandler) AddSpending(c *fiber.Ctx) error {
	payload := new(httpdto.DaySpendings)
	if err := c.BodyParser(payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON("failed to parse input")
	}

	fmt.Println(payload.Day)
	fmt.Println(payload.Sum)

	input, err := converter.ToAddSpendingFromHandler(payload)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON("failed to convert input")
	}

	_, err = sh.spendingService.AddSpending(context.Background(), input)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON("failed to add spending")
	}

	return c.Status(fiber.StatusCreated).JSON("spending added")
}

func (sh *SpendingHandler) GetWeekSpendings(c *fiber.Ctx) error {
	weekSpendings, err := sh.spendingService.GetWeekSpendings(context.Background(), time.Now())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON("failed to fetch spendings")
	}

	return c.Status(fiber.StatusOK).JSON(converter.ToGetWeekSpendingsHTTPResponseFromServer(weekSpendings))
}
