package repository

import (
	"bookstore/internal/model"
	"bookstore/pkg/utils"
	"log"

	"database/sql"
)

type OrderRepository interface {
	AddOrUpdateCart(orderID, bookID, quantity int, subtotal int64) error
	RemoveFromCart(orderId int, bookId int) error
	GetCart(orderId int) (model.OrderResponse, error)
	GetOrderHistory(customerID, limit, page int) ([]model.OrderResponse, error)
	CreateOrderIfNotExists(customerID int) (int, error)
	PayOrder(customerID int) error
}

type orderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) GetCart(orderId int) (model.OrderResponse, error) {
	query := `SELECT o.id, o.total,
			  d.id AS detail_id, d.book_id, d.quantity, d.subtotal,
			  b.title, b.author, b.price
			  FROM orders o
			  JOIN order_details d ON o.id = d.order_id
			  JOIN books b ON d.book_id = b.id
			  WHERE o.id = $1`

	rows, err := r.db.Query(query, orderId)
	if err != nil {
		log.Printf("[GetCart] Error retrieving cart: %v", err)
		return model.OrderResponse{}, err
	}

	defer rows.Close()

	cart, err := utils.ConvertToDetailResponse(rows)
	if err != nil {
		log.Printf("[GetCart] Error converting rows to detail response: %v", err)
		return model.OrderResponse{}, err
	}

	// Check if the cart slice is empty
	if len(cart) == 0 {
		return model.OrderResponse{
			ID:          int64(orderId),
			OrderDetail: []model.OrderDetailResponse{},
			Total:       0,
		}, nil
	}

	// Return the first OrderResponse
	return cart[0], nil
}

func (r *orderRepository) AddOrUpdateCart(orderID, bookID, quantity int, subtotal int64) error {
	// Begin a transaction
	tx, err := r.db.Begin()
	if err != nil {
		log.Printf(
			"[AddOrUpdateCart] Could not start transaction for order ID %d: %v",
			orderID,
			err,
		)
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			log.Println("Recovered from panic, rolling back transaction in AddOrUpdateCart")
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
		log.Printf(
			"[AddOrUpdateCart] Error updating order_details for order ID %d: %v",
			orderID,
			err,
		)
		return err
	}

	// Recalculate the total for the order
	if err := r.RecalculateTotalPrice(tx, orderID); err != nil {
		tx.Rollback()
		log.Printf("[AddOrUpdateCart] Error updating order total for order ID %d: %v", orderID, err)
		return err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		log.Printf(
			"[AddOrUpdateCart] Could not commit transaction for order ID %d: %v",
			orderID,
			err,
		)
		return err
	}

	return nil
}

func (r *orderRepository) RemoveFromCart(orderID, bookID int) error {
	// Begin a transaction to handle potential rollback in case of errors
	tx, err := r.db.Begin()
	if err != nil {
		log.Printf("[RemoveFromCart] Could not start transaction for order ID %d: %v", orderID, err)
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			log.Println("Recovered from panic, rolling back transaction in RemoveFromCart")
			tx.Rollback()
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
		log.Printf(
			"[RemoveFromCart] Error deleting from order_details for order ID %d, book ID %d: %v",
			orderID,
			bookID,
			err,
		)
		return err
	}

	// Recalculate the total for the order
	if err := r.RecalculateTotalPrice(tx, orderID); err != nil {
		tx.Rollback()
		log.Printf("[RemoveFromCart] Error updating order total for order ID %d: %v", orderID, err)
		return err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		log.Printf(
			"[RemoveFromCart] Could not commit transaction for order ID %d: %v",
			orderID,
			err,
		)
		return err
	}

	return nil
}

// GetOrderHistory retrieves all paid orders for a specific customer
func (r *orderRepository) GetOrderHistory(
	customerID, limit, page int,
) ([]model.OrderResponse, error) {
	query := `SELECT o.id, o.total,
			  d.id AS detail_id, d.book_id, d.quantity, d.subtotal,
			  b.title, b.author, b.price
			  FROM (
				SELECT * FROM orders
				WHERE customer_id = $1 AND order_state != 1
				ORDER BY updated_at ASC
				LIMIT $2 OFFSET $3
			  ) o
			  JOIN order_details d ON o.id = d.order_id
			  JOIN books b ON d.book_id = b.id`

	rows, err := r.db.Query(query, customerID, limit, page*limit)
	if err != nil {
		log.Printf(
			"[GetOrderHistory] Error retrieving paid orders for customer ID %d: %v",
			customerID,
			err,
		)
		return nil, err
	}

	defer rows.Close()

	log.Println(rows)
	return utils.ConvertToDetailResponse(rows)
}

func (r *orderRepository) CreateOrderIfNotExists(customerID int) (int, error) {
	var id int

	err := r.db.QueryRow(`
		SELECT id FROM orders
		WHERE customer_id = $1 AND order_state = 1
	`, customerID).Scan(&id)

	if err == nil {

		return id, nil
	} else if err != sql.ErrNoRows {
		// If there's another error, log and return it
		log.Printf("[CreateOrderIfNotExists] Error checking for existing order for customer ID %d: %v", customerID, err)
		return 0, err
	}

	// If no existing order, create a new one
	query := `
	INSERT INTO orders (customer_id, updated_at, total)
	VALUES ($1, NOW(), 0)
	RETURNING id;`

	if err := r.db.QueryRow(query, customerID).Scan(&id); err != nil {
		log.Printf(
			"[CreateOrderIfNotExists] Error creating order for customer ID %d: %v",
			customerID,
			err,
		)
		return 0, err
	}
	return id, nil
}

func (r *orderRepository) PayOrder(customerID int) error {
	// Begin a transaction
	tx, err := r.db.Begin()
	if err != nil {
		log.Printf("[PayOrder] Could not start transaction for customer ID %d: %v", customerID, err)
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			log.Println("Recovered from panic, rolling back transaction in PayOrder")
			tx.Rollback()
		}
	}()

	// Update the order state to indicate it has been paid for the first order that matches the customerID
	_, err = tx.Exec(`
    UPDATE orders
    SET order_state = 2, updated_at = NOW()
    WHERE customer_id = $1 AND order_state = 1
    RETURNING id;`, customerID) // Only update if the current state is 1
	if err != nil {
		tx.Rollback()
		log.Printf("[PayOrder] Error updating order state for customer ID %d: %v", customerID, err)
		return err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		log.Printf(
			"[PayOrder] Could not commit transaction for customer ID %d: %v",
			customerID,
			err,
		)
		return err
	}

	return nil
}

func (r *orderRepository) RecalculateTotalPrice(tx *sql.Tx, orderID int) error {
	_, err := tx.Exec(`
	WITH subtotal_sum AS (
		SELECT SUM(d.subtotal) as total_sum
		FROM order_details d
		WHERE d.order_id = $1)
		UPDATE orders o SET total = (SELECT total_sum FROM subtotal_sum), updated_at = NOW()
		WHERE o.id = $1`, orderID)

	if err != nil {
		log.Printf(
			"[RecalculateTotalPrice] Error recalculating total price for order ID %d: %v",
			orderID,
			err,
		)
	}
	return err
}
