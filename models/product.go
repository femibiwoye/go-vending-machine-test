package models

type Product struct {
	ID          uint   `gorm:"primaryKey" json:"id,omitempty"`
	Cost        int    `json:"cost" validate:"required"`
	ProductName string `json:"product_name" validate:"required"`
	SellerId    uint   `json:"seller_id,omitempty"`
}

type ProductUpdate struct {
	Cost        int    `json:"cost"`
	ProductName string `json:"product_name"`
}
