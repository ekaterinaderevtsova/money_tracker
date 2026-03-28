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
		expenseGroup.Get("/weekly", h.ExpenseHandler.GetWeeklyExpenses)
		expenseGroup.Post("", h.ExpenseHandler.AddExpense)
		expenseGroup.Delete("", h.ExpenseHandler.DeleteDayExpenses)
	}
}
