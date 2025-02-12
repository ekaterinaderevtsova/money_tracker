package dto

type DaySpendings struct {
	Day string `json:"day"`
	Sum int32  `json:"sum"`
}

type WeeklySpendings struct {
	Days  [7]DaySpendings `json:"daySpendings"`
	Total int32           `json:"total"`
	Left  int32           `json:"left"`
}

type MonthSpendings struct {
	Weeks []struct {
		Week   int   `json:"week"`
		Amount int32 `json:"amount"`
	} `json:"weeks"`
	Total int32 `json:"total"`
}
