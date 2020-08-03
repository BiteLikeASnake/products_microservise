package model

//AdminProduct используется для Admin handlers
type AdminProduct struct {
	ProductId          int64   `json:"product_id"`
	ProductName        string  `json:"product_name"`
	SupplyProductPrice float64 `json:"base_product_price"`
	SaleProductPrice   float64 `json:"sale_product_price"`
	CategoryId         int64   `json:"category_id"`
	ProductQuantity    int64   `json:"product_quantity"`
	ProductActive      bool    `json:"product_active"`
}
