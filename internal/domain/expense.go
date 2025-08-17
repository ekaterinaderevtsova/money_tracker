package domain

type DailyExpense struct {
	Date   string
	Amount int32
}

type WeeklyExpense struct {
	DailyExpenses [7]DailyExpense
	TotalAmount   int32
	AverageDaily  int32
}
