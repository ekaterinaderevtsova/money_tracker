package handler

import "cmd/main.go/internal/service"

type TelegramHandler struct {
	SpendingHandler *SpendingHandler
}

func NewTelegramHandler(service *service.Service) *TelegramHandler {
	return &TelegramHandler{
		SpendingHandler: NewSpendingHandler(service.SpendingService),
	}
}
