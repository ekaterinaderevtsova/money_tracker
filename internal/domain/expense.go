package domain

type DailyExpense struct {
	ID     string `json:"id"`
	Date   string `json:"date"`
	Amount int32  `json:"amount"`
}

type WeeklyExpense struct {
	DailyExpenses [7]DailyExpense
	TotalAmount   int32
	AverageDaily  int32
}
