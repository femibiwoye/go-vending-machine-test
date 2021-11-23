package models

type DepositRequest struct {
	Amount int `json:"amount" validate:"required"`
}
type BuyRequest struct {
	ProductID int `json:"product_id" validate:"required"`
	Quantity  int `json:"quantity" validate:"required"`
}

type BuyResponse struct {
	ProductID         int `json:"product_id"`
	QuantityPurchased int `json:"quantity_purchased"`
	AmountSpent       int `json:"amount_spent"`
	Change            int `json:"change"`
}
