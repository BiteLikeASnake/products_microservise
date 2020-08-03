package model

//UserProduct используем для user headers
type JsonProduct struct {
	ProductId       int64   `json:"product_id,omitempty"`
	ProductName     string  `json:"product_name,omitempty"`
	ProductPrice    float64 `json:"product_price,omitempty"`
	ProductQuantity int64   `json:"product_quantity,omitempty"`
	UserCategory
}
