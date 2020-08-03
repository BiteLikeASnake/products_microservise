package model

type Product struct {
	Product_id           int64 `gorm:"primary_key"`
	Product_name         string
	Supply_product_price float64
	Sale_product_price   float64
	Category_id          *int64
	Product_quantity     int64
	Product_active       bool
}

type UserProduct struct {
	Product_id           int64 `gorm:"primary_key"`
	Product_name         string
	Supply_product_price float64
	Sale_product_price   float64
	Category
	Product_quantity int64
	Product_active   bool
}
