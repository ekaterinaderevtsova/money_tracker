package handler

import (
	"context"
	"moneytracker/internal/service"
)

type HTTPHandler struct {
	ExpenseHandler *ExpenseHandler
	AuthHandler    *AuthHandler
}

func NewHTTPHandler(ctx context.Context, service *service.Service) *HTTPHandler {
	return &HTTPHandler{
		ExpenseHandler: NewExpenseHandler(ctx, service.ExpenseService),
		AuthHandler:    NewAuthHandler(),
	}
}
