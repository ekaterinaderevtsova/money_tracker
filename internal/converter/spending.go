package converter

import (
	"cmd/main.go/internal/domain"
	"cmd/main.go/internal/transport/telegram/dto"
)

func ToGetWeekSpendingsResponseFromServer(domainWeekSpendings *domain.WeeklySpendings) *dto.WeeklySpendings {
	return &dto.WeeklySpendings{
		Days:  domainWeekSpendings.Days,
		Total: domainWeekSpendings.Total,
		Left:  int32(7000) - domainWeekSpendings.Total,
	}
}

func ToGetMonthSpendingsResponse(domainWeekSpendings []domain.WeekSpending) *dto.MonthSpendings {
	var total int32
	weeks := make([]struct {
		Week   int   `json:"week"`
		Amount int32 `json:"amount"`
	}, len(domainWeekSpendings))

	for i, weekSpending := range domainWeekSpendings {
		weeks[i] = struct {
			Week   int   `json:"week"`
			Amount int32 `json:"amount"`
		}{
			Week:   weekSpending.Week,
			Amount: weekSpending.Amount,
		}
		total += weekSpending.Amount
	}

	return &dto.MonthSpendings{
		Weeks: weeks,
		Total: total,
	}
}
