package handler

import (
	"cmd/main.go/internal/service"
	"context"

	"go.uber.org/zap"
)

type HTTPHandler struct {
	SpendingHandler *SpendingHandler
}

func NewHTTPHandler(ctx context.Context, logger *zap.Logger, service *service.Service) *HTTPHandler {
	return &HTTPHandler{
		SpendingHandler: NewSpendingHandler(ctx, logger, service.SpendingService),
	}
}
