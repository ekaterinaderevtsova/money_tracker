package handler

import (
	"github.com/gofiber/fiber/v2"
)

func (h *HTTPHandler) SetSpendingRoutes(router fiber.Router) {
	authGroup := router.Group("/auth")

	authGroup.Post("/login", h.AuthHandler.Login)
	authGroup.Post("/refresh", h.AuthHandler.Refresh)

	spendingsGroup := router.Group("/spending/spendings")
	spendingsGroup.Use(JWTMiddleware)
	{
		spendingsGroup.Get("", h.SpendingHandler.GetWeekSpendings)
		spendingsGroup.Post("", h.SpendingHandler.AddSpending)
	}
}
