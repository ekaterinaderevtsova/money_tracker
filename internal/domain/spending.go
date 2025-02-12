package domain

import "time"

type AddSpending struct {
	Date time.Time `json:"date"`
	Sum  int32     `json:"sum"`
}

type DaySpendings struct {
	Day string
	Sum int32
}

type WeekSpendings struct {
	DaySpendings [7]DaySpendings
	Total        int32
}

type WeekTotalSpending struct {
	Week   int
	Amount int32
}

type MonthlySpendings struct {
	Weeks map[string]int32
	Total int32
}
