package handler

import (
	"context"
	"moneytracker/internal/converter"
	"moneytracker/internal/domain"
	"time"

	"github.com/gofiber/fiber/v2"
)

type IExpenseService interface {
	AddExpense(ctx context.Context, payload *domain.DailyExpense) error
	GetWeeklyExpenses(ctx context.Context, date string) (*domain.WeeklyExpense, error)
	DeleteyExpensesByDate(ctx context.Context, date string) error
}

type ExpenseHandler struct {
	ctx            context.Context
	expenseService IExpenseService
}

func NewExpenseHandler(ctx context.Context, expenseService IExpenseService) *ExpenseHandler {
	return &ExpenseHandler{ctx: ctx, expenseService: expenseService}
}

func (sh *ExpenseHandler) AddExpense(c *fiber.Ctx) error {
	payload := new(domain.DailyExpense)
	if err := c.BodyParser(payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON("failed to parse input")
	}

	err := sh.expenseService.AddExpense(sh.ctx, payload)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON("failed to add spending")
	}

	return c.Status(fiber.StatusCreated).JSON("spending added")
}

func (sh *ExpenseHandler) GetWeeklyExpenses(c *fiber.Ctx) error {
	date := c.Query("date")
	if date == "" {
		return c.Status(fiber.StatusBadRequest).JSON("date query parameter is required")
	}

	weekSpendings, err := sh.expenseService.GetWeeklyExpenses(sh.ctx, date)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON("failed to fetch spendings")
	}

	return c.Status(fiber.StatusOK).JSON(converter.ToWeeklyExpenseSummaryHTTPResponse(weekSpendings))
}

func (sh *ExpenseHandler) DeleteDayExpenses(c *fiber.Ctx) error {
	date := c.Query("date")
	if date == "" {
		return c.Status(fiber.StatusBadRequest).JSON("date query parameter is required")
	}
	if _, err := time.Parse("02-01-2006", date); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON("date must be in DD-MM-YYYY format")
	}

	err := sh.expenseService.DeleteyExpensesByDate(sh.ctx, date)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON("failed to delete spendings")
	}

	return c.SendStatus(fiber.StatusNoContent)
}
