package model

import (
	"time"
)

type Order struct {
	ID         int       `json:"id"`
	CustomerID int       `json:"customer_id"`
	UpdatedAt  time.Time `json:"updated_at"`
	OrderState int       `json:"order_state"`
	Total      float64   `json:"total"`
}

type OrderDetail struct {
	ID       int     `json:"id"`
	OrderID  int     `json:"order_id"`
	BookID   int     `json:"book_id"`
	Quantity int     `json:"quantity"`
	Subtotal float64 `json:"subtotal"`
}

type OrderResponse struct {
	ID          int                   `json:"id"`
	OrderDetail []OrderDetailResponse `json:"orderDetails"`
	Total       float64               `json:"total"`
}

type OrderDetailResponse struct {
	ID       int     `json:"id"`
	Book     []Book  `json:"books"`
	Quantity int     `json:"quantity"`
	Subtotal float64 `json:"subtotal"`
}
