package domain

import "time"

const SpendingsKey = "spendings:"
const TotalKey = "total:"

type AddSpending struct {
	Date time.Time `json:"date"`
	Sum  int32     `json:"sum"`
}

type DaySpendings struct {
	Day string `json:"day"`
	Sum int32  `json:"sum"`
}

type WeekSpendings struct {
	DaySpendings [7]DaySpendings
	Total        int32
	Average      int32
}

type Response struct {
	Status      string      `json:"status,omitempty"`
	Code        int         `json:"code,omitempty"`
	Description string      `json:"description,omitempty"`
	Error       string      `json:"error,omitempty"`
	Content     interface{} `json:"content,omitempty"`
	//Content     Content `json:"content"`
}
