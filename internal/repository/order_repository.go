package repository

import (
	"bookstore/internal/model"
	"database/sql"
)

type OrderRepository interface {
	CreateOrder(order *model.Order) error
	GetOrdersByCustomerID(customerID int) ([]model.Order, error)
	AddBookToOrder(orderDetail *model.OrderDetail) error
}

type orderRepository struct {
	db *sql.DB
}

// AddBookToOrder inserts a new book into the order details
func (r *orderRepository) AddBookToOrder(orderDetail *model.OrderDetail) error {
	query := `
		INSERT INTO order_details (order_id, book_id, quantity, subtotal) 
		VALUES ($1, $2, $3, $4)`
	_, err := r.db.Exec(
		query,
		orderDetail.OrderID,
		orderDetail.BookID,
		orderDetail.Quantity,
		orderDetail.Subtotal,
	)
	return err
}
