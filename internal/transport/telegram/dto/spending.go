package dto

type WeeklySpendings struct {
	Days  map[string]int32 `json:"days"`
	Total int32            `json:"total"`
}
