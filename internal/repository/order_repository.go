package repository

import (
	"bookstore/internal/model"
	"bookstore/pkg/utils"
	"log"

	"database/sql"
	"fmt"
)

type OrderRepository interface {
	AddOrUpdateCart(orderID, bookID, quantity int, subtotal int64) error
	RemoveFromCart(orderId int, bookId int) error
	GetCart(customerID int) (model.OrderResponse, error)
	GetPaidOrder(customerID int) ([]model.OrderResponse, error)
	CreateOrderIfNotExists(customerID int) (int, error)
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

	cart, err := utils.ConvertToDetailResponse(rows)
	if err != nil {
		return model.OrderResponse{}, err
	}

	// Check if the cart slice is empty
	if len(cart) == 0 {
		return model.OrderResponse{}, nil
	}

	// Return the first OrderResponse
	return cart[0], nil
}

// AddToCart adds a book to the specified order

func (r *orderRepository) AddOrUpdateCart(orderID, bookID, quantity int, subtotal int64) error {
	// Begin a transaction
	tx, err := r.db.Begin()
	if err != nil {
		log.Printf("could not start transaction: %v", err)
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			log.Println("Recovered from panic, rolling back transaction")
			tx.Rollback()
		}
	}()

	_, err = tx.Exec(`
	INSERT INTO order_details (order_id, book_id, quantity, subtotal)
	VALUES ($1, $2, $3, $4)
	ON CONFLICT (order_id, book_id) 
	DO UPDATE SET 
		quantity = $3,
		subtotal = $4;
	`, orderID, bookID, quantity, subtotal)
	if err != nil {
		tx.Rollback()
		log.Printf("error updating order_details: %v", err)
		return err
	}

	// Recalculate the total for the order
	if err := r.RecalculateTotalPrice(tx, orderID); err != nil {
		tx.Rollback()
		log.Printf("error updating order total: %v", err)
		return err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		log.Printf("could not commit transaction: %v", err)
		return err
	}

	// Return the updated cart after the book was added
	return nil
}

func (r *orderRepository) RemoveFromCart(orderID, bookID int) error {
	// Begin a transaction to handle potential rollback in case of errors
	tx, err := r.db.Begin()
	if err != nil {
		log.Printf("could not start transaction: %v", err)
		return fmt.Errorf("could not start transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			log.Println("Recovered from panic, rolling back transaction")
			tx.Rollback() // Rollback in case of panic
		}
	}()

	// Delete the book from order_details based on orderID and bookID
	_, err = tx.Exec(
		`DELETE FROM order_details WHERE order_id = $1 AND book_id = $2`,
		orderID,
		bookID,
	)
	if err != nil {
		tx.Rollback()
		log.Printf("error deleting from order_details: %v", err)
		return err
	}

	// Recalculate the total for the order
	if err := r.RecalculateTotalPrice(tx, orderID); err != nil {
		tx.Rollback()
		log.Printf("error updating order total: %v", err)
		return err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		log.Printf("could not commit transaction: %v", err)
		return fmt.Errorf("could not commit transaction: %w", err)
	}

	// Return the updated cart (after the book was removed)
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
	return utils.ConvertToDetailResponse(rows)
}

// CreateOrderIfNotExists checks if an order exists and creates a new one if it doesn't
func (r *orderRepository) CreateOrderIfNotExists(customerID int) (int, error) {
	var id int

	// First, check if an order already exists for the customer
	err := r.db.QueryRow(`
		SELECT id FROM orders
		WHERE customer_id = $1 AND order_state = 1
	`, customerID).Scan(&id)

	if err == nil {
		// If no error, return the existing order ID
		return id, nil
	} else if err != sql.ErrNoRows {
		// If there's another error, log and return it
		log.Printf("Error checking for existing order: %v", err)
		return 0, err
	}

	// If no existing order, create a new one
	query := `
	INSERT INTO orders (customer_id, updated_at, total)
	VALUES ($1, NOW(), 0)
	RETURNING id;`

	if err := r.db.QueryRow(query, customerID).Scan(&id); err != nil {
		log.Printf("Error creating order: %v", err)
		return 0, err
	}
	return id, nil
}

func (r *orderRepository) RecalculateTotalPrice(tx *sql.Tx, orderID int) error {
	_, err := tx.Exec(`
	WITH subtotal_sum AS (
		SELECT SUM(d.subtotal) as total_sum
		FROM order_details d
		WHERE d.order_id = $1)
		UPDATE orders o SET total = (SELECT total_sum FROM subtotal_sum), updated_at = NOW()
		WHERE o.id = $1`, orderID)
	return err
}
