package httpdto

type DaySpendings struct {
	Day string `json:"day"`
	Sum int32  `json:"sum"`
}

type DaySpendingsResponse struct {
	Date      string `json:"date"`
	DayOfWeek string `json:"dayOfWeek"`
	Sum       int32  `json:"sum"`
}

type WeeklySpendings struct {
	Days  [7]DaySpendingsResponse `json:"daySpendings"`
	Total int32                   `json:"total"`
}

type MonthSpendings struct {
	Weeks []struct {
		Week   int   `json:"week"`
		Amount int32 `json:"amount"`
	} `json:"weeks"`
	Total int32 `json:"total"`
}
