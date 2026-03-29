package handler

import (
	"context"
	"moneytracker/internal/converter"
	"moneytracker/internal/domain"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type IExpenseService interface {
	AddExpense(ctx context.Context, payload *domain.DailyExpense) error
	GetWeeklyExpenses(ctx context.Context, date string) (*domain.WeeklyExpense, error)
	DeleteExpensesByDate(ctx context.Context, date string) error
	GetDailyExpense(ctx context.Context, date string) ([]domain.DailyExpense, error)
	UpdateDailyExpense(ctx context.Context, uuid string, newAmount int32) error
	DeleteDailyExpense(ctx context.Context, uuid string) error
}

type ExpenseHandler struct {
	ctx            context.Context
	expenseService IExpenseService
	logger         *zap.Logger
}

func NewExpenseHandler(ctx context.Context, expenseService IExpenseService, logger *zap.Logger) *ExpenseHandler {
	return &ExpenseHandler{ctx: ctx, expenseService: expenseService}
}

func (sh *ExpenseHandler) AddExpense(c *fiber.Ctx) error {
	payload := new(domain.DailyExpense)
	if err := c.BodyParser(payload); err != nil {
		// TODO: log
		return c.Status(fiber.StatusBadRequest).JSON("failed to parse input")
	}

	err := sh.expenseService.AddExpense(sh.ctx, payload)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON("failed to add spending")
	}

	return c.Status(fiber.StatusCreated).JSON("spending added")
}

func (sh *ExpenseHandler) GetWeeklyExpenses(c *fiber.Ctx) error {
	date := c.Params("date")
	if _, err := time.Parse("2006-01-02", date); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON("date must be in YYYY-MM-DD format")
	}

	weekSpendings, err := sh.expenseService.GetWeeklyExpenses(sh.ctx, date)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON("failed to fetch spendings")
	}

	return c.Status(fiber.StatusOK).JSON(converter.ToWeeklyExpenseSummaryHTTPResponse(weekSpendings))
}

func (sh *ExpenseHandler) GetDailyExpense(c *fiber.Ctx) error {
	date := c.Params("date")
	if _, err := time.Parse("2006-01-02", date); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON("date must be in YYYY-MM-DD format")
	}

	dailyExpenses, err := sh.expenseService.GetDailyExpense(sh.ctx, date)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON("failed to fetch spendings")
	}

	// TODO: convert to dto
	return c.Status(fiber.StatusOK).JSON(dailyExpenses)
}

func (sh *ExpenseHandler) UpdateDailyExpense(c *fiber.Ctx) error {
	uuid := c.Params("uuid")
	if uuid == "" {
		return c.Status(fiber.StatusBadRequest).JSON("uuid is required")
	}

	payload := new(domain.DailyExpense)
	if err := c.BodyParser(payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON("failed to parse input")
	}

	err := sh.expenseService.UpdateDailyExpense(sh.ctx, uuid, payload.Amount)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON("failed to update spending")
	}

	return c.SendStatus(fiber.StatusOK)
}

func (sh *ExpenseHandler) DeleteDailyExpense(c *fiber.Ctx) error {
	uuid := c.Params("uuid")
	if uuid == "" {
		return c.Status(fiber.StatusBadRequest).JSON("uuid is required")
	}

	err := sh.expenseService.DeleteDailyExpense(sh.ctx, uuid)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON("failed to delete spending")
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (sh *ExpenseHandler) DeleteDayExpenses(c *fiber.Ctx) error {
	date := c.Params("date")
	if _, err := time.Parse("2006-01-02", date); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON("date must be in YYYY-MM-DD format")
	}

	err := sh.expenseService.DeleteExpensesByDate(sh.ctx, date)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON("failed to delete spendings")
	}

	return c.SendStatus(fiber.StatusNoContent)
}
