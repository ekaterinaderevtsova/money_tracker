package handler

import "github.com/gofiber/fiber/v2"

func (h *HTTPHandler) SetSpendingRoutes(router fiber.Router) {
	spendings := router.Group("/spending/spendings")
	{
		spendings.Get("", h.SpendingHandler.GetWeekSpendings)
		spendings.Post("", h.SpendingHandler.AddSpending)
	}
}
