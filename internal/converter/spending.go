package converter

import (
	"cmd/main.go/internal/domain"
	httpdto "cmd/main.go/internal/transport/http/dto"
	"time"
)

func ToTimeFromString(payload *httpdto.Date) (time.Time, error) {
	parsedTime, err := time.Parse("2006-01-02", payload.Date) // Changed format layout
	if err != nil {
		return time.Time{}, err
	}
	return parsedTime, nil
}

func ToAddSpendingFromHandler(payload *httpdto.DaySpendings) (*domain.AddSpending, error) {
	parsedTime, err := time.Parse("2006-01-02", payload.Day) // Changed format layout
	if err != nil {
		return nil, err
	}
	return &domain.AddSpending{
		Date: parsedTime,
		Sum:  payload.Sum,
	}, nil
}

func ToGetWeekSpendingsHTTPResponseFromServer(domainWeekSpendings *domain.WeekSpendings) *httpdto.WeeklySpendings {
	var days [7]httpdto.DaySpendingsResponse
	daysOfWeek := [7]string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}
	for i, domainDay := range domainWeekSpendings.DaySpendings {
		days[i] = httpdto.DaySpendingsResponse{
			Date:      domainDay.Day,
			DayOfWeek: daysOfWeek[i],
			Sum:       domainDay.Sum,
		}
	}

	return &httpdto.WeeklySpendings{
		Days:  days,
		Total: domainWeekSpendings.Total,
	}
}
