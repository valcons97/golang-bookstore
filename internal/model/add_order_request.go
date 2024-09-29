package model

type AddOrderRequest struct {
	BookId   int64   `json:"bookId"   binding:"required"`
	Quantity int64   `json:"quantity" binding:"required"`
	Price    float64 `json:"price"    binding:"required"`
}
