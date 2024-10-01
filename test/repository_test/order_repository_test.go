package repository_test

import (
	"bookstore/internal/model"
	"bookstore/internal/repository"
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

var recalculationQuery = regexp.QuoteMeta(
	`WITH subtotal_sum AS (
        SELECT SUM(d.subtotal) as total_sum
        FROM order_details d
        WHERE d.order_id = $1)
        UPDATE orders o SET total = (SELECT total_sum FROM subtotal_sum), updated_at = NOW()
        WHERE o.id = $1`,
)

func TestOrderRepository_GetCart(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	orderRepo := repository.NewOrderRepository(db)

	orderID := 1

	query := `SELECT o.id, o.total,
    d.id AS detail_id, d.book_id, d.quantity, d.subtotal,
    b.title, b.author, b.price
    FROM orders o
    JOIN order_details d ON o.id = d.order_id
    JOIN books b ON d.book_id = b.id
    WHERE o.id = \$1`

	t.Run("successful retrieval of cart", func(t *testing.T) {
		mock.ExpectQuery(query).
			WithArgs(orderID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "total", "detail_id", "book_id", "quantity", "subtotal", "title", "author", "price"}).
				AddRow(1, 400, 1, 1, 2, 400, "Book Title", "Author Name", 200))

		result, err := orderRepo.GetCart(orderID)

		expected := model.OrderResponse{
			ID: int64(orderID),
			OrderDetail: []model.OrderDetailResponse{
				{
					ID: 1,
					Book: []model.Book{
						{ID: 1, Title: "Book Title", Author: "Author Name", Price: 2.00},
					},
					Quantity: 2,
					Subtotal: 4.00,
				},
			},
			Total: 4.00,
		}

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("no cart found", func(t *testing.T) {
		mock.ExpectQuery(query).
			WithArgs(orderID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "total", "detail_id", "book_id", "quantity", "subtotal", "title", "author", "price"}))

		result, err := orderRepo.GetCart(orderID)

		expected := model.OrderResponse{
			ID:          int64(orderID),
			OrderDetail: []model.OrderDetailResponse{},
			Total:       0,
		}

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error retrieving cart", func(t *testing.T) {
		mock.ExpectQuery(query).
			WithArgs(orderID).
			WillReturnError(sql.ErrNoRows)

		result, err := orderRepo.GetCart(orderID)

		assert.Error(t, err)
		assert.Equal(t, model.OrderResponse{}, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestOrderRepository_GetOrderHistory(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	orderRepo := repository.NewOrderRepository(db)

	customerID := 1
	limit := 10
	page := 0

	query := `SELECT o.id, o.total, d.id AS detail_id, d.book_id, d.quantity, d.subtotal, b.title, b.author, b.price`

	t.Run("successful retrieval of order history", func(t *testing.T) {
		mock.ExpectQuery(query).
			WithArgs(customerID, limit, page*limit).
			WillReturnRows(sqlmock.NewRows([]string{"id", "total", "detail_id", "book_id", "quantity", "subtotal", "title", "author", "price"}).
				AddRow(1, 3197, 1, 1, 2, 2398, "1984", "George Orwell", 999).
				AddRow(1, 3197, 2, 2, 1, 799, "To Kill a Mockingbird", "Harper Lee", 799))

		orders, err := orderRepo.GetOrderHistory(customerID, limit, page)
		assert.NoError(t, err)
		assert.Len(t, orders, 1)

		expectedOrder := model.OrderResponse{
			ID:    1,
			Total: 31.97,
			OrderDetail: []model.OrderDetailResponse{
				{
					ID:       1,
					Quantity: 2,
					Subtotal: 23.98,
					Book: []model.Book{
						{
							ID:     1,
							Title:  "1984",
							Author: "George Orwell",
							Price:  9.99,
						},
					},
				},
				{
					ID:       2,
					Quantity: 1,
					Subtotal: 7.99,
					Book: []model.Book{
						{
							ID:     2,
							Title:  "To Kill a Mockingbird",
							Author: "Harper Lee",
							Price:  7.99,
						},
					},
				},
			},
		}

		assert.Equal(t, expectedOrder, orders[0])
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error when retrieving order history", func(t *testing.T) {
		mock.ExpectQuery(query).
			WithArgs(customerID, limit, page*limit).
			WillReturnError(sql.ErrNoRows)

		orders, err := orderRepo.GetOrderHistory(customerID, limit, page)
		assert.Error(t, err)
		assert.Nil(t, orders)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestOrderRepository_PayOrder(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	orderRepo := repository.NewOrderRepository(db)

	customerID := 1

	query := regexp.QuoteMeta(
		`UPDATE orders SET order_state = 2, updated_at = NOW() WHERE customer_id = $1 AND order_state = 1 RETURNING id;`,
	)

	t.Run("successful payment of order", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(query).
			WithArgs(customerID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		err := orderRepo.PayOrder(customerID)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error when starting transaction", func(t *testing.T) {
		mock.ExpectBegin().WillReturnError(sql.ErrConnDone)

		err := orderRepo.PayOrder(customerID)
		assert.Error(t, err)
		assert.EqualError(t, err, "sql: connection is already closed")
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error when executing update", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(query).
			WithArgs(customerID).
			WillReturnError(sql.ErrNoRows)

		mock.ExpectRollback()

		err := orderRepo.PayOrder(customerID)
		assert.Error(t, err)
		assert.EqualError(t, err, "sql: no rows in result set")
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error when committing transaction", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(query).
			WithArgs(customerID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit().WillReturnError(sql.ErrConnDone)

		err := orderRepo.PayOrder(customerID)
		assert.Error(t, err)
		assert.EqualError(t, err, "sql: connection is already closed")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestOrderRepository_CreateOrderIfNotExists(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	orderRepo := repository.NewOrderRepository(db)

	customerID := 1

	checkQuery := `SELECT id FROM orders WHERE customer_id = \$1 AND order_state = 1`

	insertQuery := `INSERT INTO orders \(customer_id, updated_at, total\) VALUES \(\$1, NOW\(\), 0\) RETURNING id`

	t.Run("existing order found", func(t *testing.T) {
		mock.ExpectQuery(checkQuery).
			WithArgs(customerID).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		id, err := orderRepo.CreateOrderIfNotExists(customerID)
		assert.NoError(t, err)
		assert.Equal(t, 1, id)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("no existing order found, create new order", func(t *testing.T) {
		mock.ExpectQuery(checkQuery).
			WithArgs(customerID).
			WillReturnError(sql.ErrNoRows)

		mock.ExpectQuery(insertQuery).
			WithArgs(customerID).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(2))

		id, err := orderRepo.CreateOrderIfNotExists(customerID)
		assert.NoError(t, err)
		assert.Equal(t, 2, id)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error while checking existing orders", func(t *testing.T) {
		mock.ExpectQuery(checkQuery).
			WithArgs(customerID).
			WillReturnError(sql.ErrConnDone)

		id, err := orderRepo.CreateOrderIfNotExists(customerID)
		assert.Error(t, err)
		assert.Equal(t, 0, id)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error while creating new order", func(t *testing.T) {
		mock.ExpectQuery(checkQuery).
			WithArgs(customerID).
			WillReturnError(sql.ErrNoRows)

		mock.ExpectQuery(insertQuery).
			WithArgs(customerID).
			WillReturnError(sql.ErrConnDone)

		id, err := orderRepo.CreateOrderIfNotExists(customerID)
		assert.Error(t, err)
		assert.Equal(t, 0, id)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestOrderRepository_RemoveFromCart(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	orderRepo := repository.NewOrderRepository(db)

	orderID := 1
	bookID := 2

	deleteQuery := regexp.QuoteMeta(
		`DELETE FROM order_details WHERE order_id = $1 AND book_id = $2`,
	)

	t.Run("successful removal from cart", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(deleteQuery).
			WithArgs(orderID, bookID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectExec(recalculationQuery).
			WithArgs(orderID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectCommit()

		err := orderRepo.RemoveFromCart(orderID, bookID)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error starting transaction", func(t *testing.T) {
		mock.ExpectBegin().WillReturnError(sql.ErrConnDone)

		err := orderRepo.RemoveFromCart(orderID, bookID)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error deleting from cart", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(deleteQuery).
			WithArgs(orderID, bookID).
			WillReturnError(sql.ErrConnDone)

		mock.ExpectRollback()

		err := orderRepo.RemoveFromCart(orderID, bookID)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error recalculating total", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(deleteQuery).
			WithArgs(orderID, bookID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectExec(recalculationQuery).
			WithArgs(orderID).
			WillReturnError(sql.ErrConnDone)

		mock.ExpectRollback()

		err := orderRepo.RemoveFromCart(orderID, bookID)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error committing transaction", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(deleteQuery).
			WithArgs(orderID, bookID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectExec(recalculationQuery).
			WithArgs(orderID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectCommit().WillReturnError(sql.ErrConnDone)

		err := orderRepo.RemoveFromCart(orderID, bookID)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestOrderRepository_AddOrUpdateCart(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	orderRepo := repository.NewOrderRepository(db)

	orderID := 1
	bookID := 2
	quantity := 3
	subtotal := int64(600)

	insertQuery := regexp.QuoteMeta(
		`INSERT INTO order_details (order_id, book_id, quantity, subtotal) VALUES 
    ($1, $2, $3, $4) ON CONFLICT (order_id, book_id) DO UPDATE SET quantity = $3, subtotal = $4;`,
	)

	t.Run("successful add or update cart", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(insertQuery).
			WithArgs(orderID, bookID, quantity, subtotal).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectExec(recalculationQuery).
			WithArgs(orderID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectCommit()

		err := orderRepo.AddOrUpdateCart(orderID, bookID, quantity, subtotal)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error starting transaction", func(t *testing.T) {
		mock.ExpectBegin().WillReturnError(sql.ErrConnDone)

		err := orderRepo.AddOrUpdateCart(orderID, bookID, quantity, subtotal)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error during insert/update", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(insertQuery).
			WithArgs(orderID, bookID, quantity, subtotal).
			WillReturnError(sql.ErrConnDone)

		mock.ExpectRollback()

		err := orderRepo.AddOrUpdateCart(orderID, bookID, quantity, subtotal)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error recalculating total", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(insertQuery).
			WithArgs(orderID, bookID, quantity, subtotal).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectExec(recalculationQuery).
			WithArgs(orderID).
			WillReturnError(sql.ErrConnDone)

		mock.ExpectRollback()

		err := orderRepo.AddOrUpdateCart(orderID, bookID, quantity, subtotal)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error committing transaction", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(insertQuery).
			WithArgs(orderID, bookID, quantity, subtotal).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectExec(recalculationQuery).
			WithArgs(orderID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectCommit().WillReturnError(sql.ErrConnDone)

		err := orderRepo.AddOrUpdateCart(orderID, bookID, quantity, subtotal)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
