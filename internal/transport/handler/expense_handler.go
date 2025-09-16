package handler

import (
	"context"
	"moneytracker/internal/converter"
	"moneytracker/internal/service"
	httpdto "moneytracker/internal/transport/dto"

	"github.com/gofiber/fiber/v2"
)

type ExpenseHandler struct {
	ctx            context.Context
	expenseService service.IExpenseService
}

func NewExpenseHandler(ctx context.Context, expenseService service.IExpenseService) *ExpenseHandler {
	return &ExpenseHandler{ctx: ctx, expenseService: expenseService}
}

func (sh *ExpenseHandler) AddExpense(c *fiber.Ctx) error {
	payload := new(httpdto.DailyExpense)
	if err := c.BodyParser(payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON("failed to parse input")
	}

	err := sh.expenseService.AddExpense(sh.ctx, converter.ToDailyExpenseDomain(payload))
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
