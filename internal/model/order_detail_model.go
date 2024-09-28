package model

type OrderDetail struct {
	ID       int64   `json:"id"`
	OrderID  int     `json:"order_id"`
	BookID   int     `json:"book_id"`
	Quantity int     `json:"quantity"`
	Subtotal float64 `json:"subtotal"`
}
