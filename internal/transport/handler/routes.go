package handler

import (
	"github.com/gofiber/fiber/v2"
)

func (h *HTTPHandler) SetSpendingRoutes(router fiber.Router) {
	authGroup := router.Group("/auth")

	authGroup.Post("/login", h.AuthHandler.Login)
	authGroup.Post("/refresh", h.AuthHandler.Refresh)

	expenseGroup := router.Group("/expenses")
	expenseGroup.Use(JWTMiddleware)
	{
		expenseGroup.Post("", h.ExpenseHandler.AddExpense)
		expenseGroup.Get("/weekly/:date", h.ExpenseHandler.GetWeeklyExpenses)
		expenseGroup.Put("/:uuid", h.ExpenseHandler.UpdateDailyExpense)
		expenseGroup.Delete("/:uuid", h.ExpenseHandler.DeleteDailyExpense)
	}

	dailyExpensesGroup := expenseGroup.Group("/daily")
	{
		dailyExpensesGroup.Get("/:date", h.ExpenseHandler.GetDailyExpense)
		dailyExpensesGroup.Delete("/:date", h.ExpenseHandler.DeleteDayExpenses)
	}

}
