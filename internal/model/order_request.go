package model

type AddToCartRequest struct {
	BookId   int64   `json:"bookId"   binding:"required"`
	Quantity int64   `json:"quantity" binding:"required,gte=1"`
	Price    float64 `json:"price"    binding:"required"`
}

type RemoveItemFromCartRequest struct {
	BookId int64 `json:"bookId" binding:"required"`
}
