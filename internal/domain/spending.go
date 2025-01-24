package domain

import "time"

type AddSpending struct {
	Date time.Time `json:"date"`
	Sum  int32     `json:"sum"`
}

type WeeklySpendings struct {
	Days  map[string]int32 `json:"days"`
	Total int32            `json:"total"`
}
