package httpdto

type DailyExpense struct {
	Date   string `json:"day"`
	Amount int32  `json:"sum"`
}

type WeeklyExpense struct {
	DailyExpenses []struct {
		Date      string `json:"date"`
		DayOfWeek string `json:"dayOfWeek"`
		Amount    int32  `json:"amount"`
	} `json:"dailyExpenses"`
	TotalAmount  int32 `json:"totalAmount"`
	AverageDaily int32 `json:"averageDaily"`
}
