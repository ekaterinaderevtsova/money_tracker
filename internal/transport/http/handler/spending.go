package handler

import (
	"context"
	"moneytracker/internal/converter"
	"moneytracker/internal/domain"
	"moneytracker/internal/service"

	"go.uber.org/zap"

	"github.com/gofiber/fiber/v2"
)

type SpendingHandler struct {
	ctx             context.Context
	logger          *zap.Logger
	spendingService service.ISpendingService
}

func NewSpendingHandler(ctx context.Context, logger *zap.Logger, spendingService service.ISpendingService) *SpendingHandler {
	return &SpendingHandler{ctx: ctx, logger: logger, spendingService: spendingService}
}

func (sh *SpendingHandler) AddSpending(c *fiber.Ctx) error {
	payload := new(domain.DaySpendings)
	if err := c.BodyParser(payload); err != nil {
		sh.logger.Error("Error parsing payload", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON("failed to parse input")
	}

	err := sh.spendingService.AddSpending(sh.ctx, payload)
	if err != nil {
		sh.logger.Error("Error adding spending to db", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON("failed to add spending")
	}

	sh.logger.Info("New spending added", zap.String("date", payload.Day), zap.Int32("added spending", payload.Sum))

	return c.Status(fiber.StatusCreated).JSON("spending added")
}

func (sh *SpendingHandler) GetWeekSpendings(c *fiber.Ctx) error {
	date := c.Query("date")
	if date == "" {
		sh.logger.Error("Date parameter missing", zap.Error(domain.ErrInvalidPayload))
		return c.Status(fiber.StatusBadRequest).JSON("date query parameter is required")
	}

	weekSpendings, err := sh.spendingService.GetWeekSpendings(sh.ctx, date)
	if err != nil {
		sh.logger.Error("Error fetching week spendings", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON("failed to fetch spendings")
	}

	return c.Status(fiber.StatusOK).JSON(converter.ToGetWeekSpendingsHTTPResponseFromServer(weekSpendings))
}
