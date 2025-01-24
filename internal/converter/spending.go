package converter

import (
	"cmd/main.go/internal/domain"
	"cmd/main.go/internal/transport/telegram/dto"
)

func ToGetWeekSpendingsResponseFromServer(domainWeekSpendings *domain.WeeklySpendings) *dto.WeeklySpendings {
	return &dto.WeeklySpendings{
		Days:  domainWeekSpendings.Days,
		Total: domainWeekSpendings.Total,
	}
}
