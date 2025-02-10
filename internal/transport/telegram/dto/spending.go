package dto

type WeeklySpendings struct {
	Days  map[string]int32 `json:"days"`
	Total int32            `json:"total"`
	Left  int32            `json:"left"`
}

type MonthSpendings struct {
	Weeks []struct {
		Week   int   `json:"week"`
		Amount int32 `json:"amount"`
	} `json:"weeks"`
	Total int32 `json:"total"`
}
