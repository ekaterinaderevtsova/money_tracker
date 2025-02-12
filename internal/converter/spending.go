package converter

import (
	"cmd/main.go/internal/domain"
	"cmd/main.go/internal/transport/telegram/dto"
)

func ToGetWeekSpendingsResponseFromServer(domainWeekSpendings *domain.WeekSpendings) *dto.WeeklySpendings {
	var days [7]dto.DaySpendings
	for i, domainDay := range domainWeekSpendings.DaySpendings {
		days[i] = dto.DaySpendings{
			Day: domainDay.Day,
			Sum: domainDay.Sum,
		}
	}

	return &dto.WeeklySpendings{
		Days:  days,
		Total: domainWeekSpendings.Total,
		Left:  7000 - domainWeekSpendings.Total,
	}
}

func ToGetMonthSpendingsResponse(domainWeekSpendings []domain.WeekTotalSpending) *dto.MonthSpendings {
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
