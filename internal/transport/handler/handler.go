package handler

import (
	"context"
	"moneytracker/internal/service"

	"go.uber.org/zap"
)

type HTTPHandler struct {
	ExpenseHandler *ExpenseHandler
	AuthHandler    *AuthHandler
}

func NewHTTPHandler(ctx context.Context, service *service.Service, logger *zap.Logger) *HTTPHandler {
	return &HTTPHandler{
		ExpenseHandler: NewExpenseHandler(ctx, service.ExpenseService, logger),
		AuthHandler:    NewAuthHandler(),
	}
}
