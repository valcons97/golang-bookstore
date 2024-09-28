package repository

import (
	"bookstore/internal/model"
	converter "bookstore/internal/utils"
	"database/sql"
	"fmt"
)

type OrderRepository interface {
	AddBookToOrder(orderDetail *model.OrderDetail) error
	GetCart(customerID int) (model.OrderResponse, error)
	GetPaidOrder(customerID int) ([]model.OrderResponse, error)
	CreateOrderIfNotExists(customerID int, total float64) (int, error)
}

type orderRepository struct {
	db *sql.DB
}

// NewOrderRepository creates a new instance of orderRepository
func NewOrderRepository(db *sql.DB) OrderRepository {
	return &orderRepository{db: db}
}

// GetCart implements OrderRepository.
func (r *orderRepository) GetCart(customerID int) (model.OrderResponse, error) {
	query := `SELECT o.id, o.total,
			  d.id AS detail_id, d.book_id, d.quantity, d.subtotal,
			  b.title, b.author, b.price
			  FROM orders o
			  JOIN order_details d ON o.id = d.order_id
			  JOIN books b ON d.book_id = b.id
			  WHERE o.customer_id = $1 AND o.order_state = 1`

	rows, err := r.db.Query(query, customerID)
	if err != nil {
		return model.OrderResponse{}, fmt.Errorf("could not retrieve cart: %w", err)
	}

	cart, err := productConverter(rows)
	if err != nil {
		return model.OrderResponse{}, err
	}

	// Check if the cart slice is empty
	if len(cart) == 0 {
		return model.OrderResponse{}, fmt.Errorf("cart is empty")
	}

	// Return the first OrderResponse
	return cart[0], nil
}

// AddBookToOrder adds a book to the specified order
func (r *orderRepository) AddBookToOrder(orderDetail *model.OrderDetail) error {
	query := `INSERT INTO order_details (order_id, book_id, quantity, subtotal)
				VALUES ($1, $2, $3, $4)`

	_, err := r.db.Exec(
		query,
		orderDetail.OrderID,
		orderDetail.BookID,
		orderDetail.Quantity,
		*converter.ConvertStorePrice(&orderDetail.Subtotal),
	)
	if err != nil {
		return err
	}

	return nil
}

// GetPaidOrder retrieves all paid orders for a specific customer
func (r *orderRepository) GetPaidOrder(customerID int) ([]model.OrderResponse, error) {
	query := `SELECT o.id, o.total,
			  d.id AS detail_id, d.book_id, d.quantity, d.subtotal,
			  b.title, b.author, b.price
			  FROM orders o
			  JOIN order_details d ON o.id = d.order_id
			  JOIN books b ON d.book_id = b.id
			  WHERE o.customer_id = $1 AND o.order_state != 1`

	rows, err := r.db.Query(query, customerID)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve paid orders: %w", err)
	}

	// Use the productConverter function
	return productConverter(rows)
}

// CreateOrderIfNotExists checks if an order exists and creates a new one if it doesn't
func (r *orderRepository) CreateOrderIfNotExists(customerID int, total float64) (int, error) {
	var id int
	query := `
    INSERT INTO orders (customer_id, updated_at, total)
    SELECT ?, NOW(), ?
    WHERE NOT EXISTS (
        SELECT 1 FROM orders
        WHERE customer_id = ? AND order_state = 1
    )
	RETURNING id;`

	if err := r.db.QueryRow(query, customerID, total, customerID).Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func productConverter(rows *sql.Rows) ([]model.OrderResponse, error) {
	defer rows.Close()

	var orders []model.OrderResponse
	orderMap := make(map[int]*model.OrderResponse)

	for rows.Next() {
		var orderID, detailID, bookID, quantity int
		var total, subtotal, price float64
		var title, author string

		err := rows.Scan(
			&orderID,
			&total,
			&detailID,
			&bookID,
			&quantity,
			&subtotal,
			&title,
			&author,
			&price,
		)
		if err != nil {
			return nil, fmt.Errorf("could not scan order row: %w", err)
		}

		// If the order is not in the map, create a new OrderResponse entry
		if _, ok := orderMap[orderID]; !ok {
			orderMap[orderID] = &model.OrderResponse{
				ID:          orderID,
				OrderDetail: []model.OrderDetailResponse{},
				Total:       total,
			}
		}

		// Create a new Book entry
		book := model.Book{
			ID:     bookID,
			Title:  title,
			Author: author,
			Price:  price,
		}

		// Create a new OrderDetailResponse entry
		orderDetail := model.OrderDetailResponse{
			ID:       detailID,
			Book:     []model.Book{book}, // Assuming multiple books could be part of order details
			Quantity: quantity,
			Subtotal: subtotal,
		}

		// Append the order details to the corresponding order
		orderMap[orderID].OrderDetail = append(orderMap[orderID].OrderDetail, orderDetail)
	}

	// Convert the map to a slice
	for _, order := range orderMap {
		orders = append(orders, *order)
	}

	// Check for errors during row iteration
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("could not retrieve orders: %w", err)
	}

	return orders, nil
}
