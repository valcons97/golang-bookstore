package request

type LoginRequest struct {
	Email    string `json:"email"    binding:"required,email"` // Email field with validation
	Password string `json:"password" binding:"required"`       // Password field with validation
}

type AddToCartRequest struct {
	BookId   int64   `json:"bookId"   binding:"required"`
	Quantity int64   `json:"quantity" binding:"required,gte=1"`
	Price    float64 `json:"price"    binding:"required"`
}

type RemoveItemFromCartRequest struct {
	BookId int64 `json:"bookId" binding:"required"`
}

type HistoryRequest struct {
	Page  int `json:"page"  binding:"gte=0"`
	Limit int `json:"limit" binding:"gte=0"`
}
