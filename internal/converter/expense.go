package converter

import (
	"moneytracker/internal/domain"
	httpdto "moneytracker/internal/transport/dto"
)

func ToDailyExpenseDomain(dtoDailyExpense *httpdto.DailyExpense) *domain.DailyExpense {
	return &domain.DailyExpense{
		Date:   dtoDailyExpense.Date,
		Amount: dtoDailyExpense.Amount,
	}
}

func ToWeeklyExpenseSummaryHTTPResponse(domainWeeklyExpenses *domain.WeeklyExpense) *httpdto.WeeklyExpense {
	daysOfWeek := [7]string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}

	dailyExpenses := make([]struct {
		Date      string `json:"date"`
		DayOfWeek string `json:"dayOfWeek"`
		Amount    int32  `json:"amount"`
	}, len(domainWeeklyExpenses.DailyExpenses))

	for i, domainDay := range domainWeeklyExpenses.DailyExpenses {
		dailyExpenses[i] = struct {
			Date      string `json:"date"`
			DayOfWeek string `json:"dayOfWeek"`
			Amount    int32  `json:"amount"`
		}{
			Date:      domainDay.Date,
			DayOfWeek: daysOfWeek[i],
			Amount:    domainDay.Amount,
		}
	}

	return &httpdto.WeeklyExpense{
		DailyExpenses: dailyExpenses,
		TotalAmount:   domainWeeklyExpenses.TotalAmount,
		AverageDaily:  domainWeeklyExpenses.AverageDaily,
	}
}
