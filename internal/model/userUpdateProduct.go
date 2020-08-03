package model

//UserUpdateProduct ...
type UserUpdateProduct struct {
	ProductId       int64 `json:"product_id,omitempty"`
	AmountPurchased int64 `json:"amount_purchased,omitempty"`
}
