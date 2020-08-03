package model

//UserCategory ...
type UserCategory struct {
	CategoryId   int64  `json:"category_id,omitempty"`
	CategoryName string `json:"category_name,omitempty"`
}
