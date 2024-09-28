package model

type AddOrderRequest struct {
	BookId   string  `json:"bookId"   binding:"required"`
	Quantity int     `json:"quantity" binding:"required"`
	Price    float64 `json:"price"    binding:"required"`
}
