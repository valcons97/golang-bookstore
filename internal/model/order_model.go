package model

import (
	"time"
)

type Order struct {
	ID         int64     `json:"id"`
	CustomerID int64     `json:"customer_id"`
	UpdatedAt  time.Time `json:"updated_at"`
	OrderState int64     `json:"order_state"`
	Total      float64   `json:"total"`
}

type OrderDetail struct {
	ID       int64   `json:"id"`
	OrderID  int64   `json:"order_id"`
	BookID   int64   `json:"book_id"`
	Quantity int64   `json:"quantity"`
	Subtotal float64 `json:"subtotal"`
}

type OrderResponse struct {
	ID          int64                 `json:"id"`
	OrderDetail []OrderDetailResponse `json:"orderDetails"`
	Total       float64               `json:"total"`
}

type OrderDetailResponse struct {
	ID       int64   `json:"id"`
	Book     []Book  `json:"books"`
	Quantity int64   `json:"quantity"`
	Subtotal float64 `json:"subtotal"`
}

type OrderState int

const OrderState_One OrderState = 1  // 
const OrderState_Two OrderState = 2
