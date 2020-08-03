package model

type AdminCategory struct {
	CategoryId     int64  `json:"category_id"`
	CategoryName   string `json:"category_name"`
	CategoryActive bool   `json:"category_active"`
}
