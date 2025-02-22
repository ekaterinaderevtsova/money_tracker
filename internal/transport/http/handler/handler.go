package handler

import "cmd/main.go/internal/service"

type HTTPHandler struct {
	SpendingHandler *SpendingHandler
}

func NewHTTPHandler(service *service.Service) *HTTPHandler {
	return &HTTPHandler{
		SpendingHandler:                 NewSpendingHandler(service.SpendingService),
	}
}
