package model

import (
	"time"
)

type Order struct {
	ID         int           `json:"id"`
	CustomerID int           `json:"customer_id"`
	UpdatedAt  time.Time     `json:"updated_at"`
	OrderState int           `json:"order_state"`
	Total      float64       `json:"total"`
	Details    []OrderDetail `json:"details"`
}
