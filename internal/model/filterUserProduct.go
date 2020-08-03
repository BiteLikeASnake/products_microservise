package model

type FilterUProduct struct {
	ProductName     *string  `json:"product_name"`
	MinProductPrice *float64 `json:"min_product_price"`
	MaxProductPrice *float64 `json:"max_product_price"`
	CategoryId      *int64   `json:"category_id"`
}
