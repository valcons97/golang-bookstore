package repository_test

// import (
// 	"bookstore/internal/model"
// 	"bookstore/internal/repository"
// 	"bookstore/pkg/utils"
// 	"database/sql"
// 	"testing"

// 	"github.com/DATA-DOG/go-sqlmock"
// 	"github.com/stretchr/testify/assert"
// )

// func TestOrderRepository_AddOrUpdateCart(t *testing.T) {
// 	db, mock, err := sqlmock.New()
// 	assert.NoError(t, err)
// 	defer db.Close()

// 	orderRepo := repository.NewOrderRepository(db)

// 	t.Run("successful add to cart", func(t *testing.T) {
// 		orderID := 1
// 		bookID := 1
// 		quantity := 2
// 		subtotal := int64(5000)

// 		mock.ExpectBegin()

// 		mock.ExpectExec("INSERT INTO order_details").
// 			WithArgs(orderID, bookID, quantity, subtotal).
// 			WillReturnResult(sqlmock.NewResult(1, 1))

// 		mock.ExpectExec("WITH subtotal_sum AS").WithArgs(orderID).
// 			WillReturnResult(sqlmock.NewResult(1, 1))

// 		mock.ExpectCommit()

// 		err := orderRepo.AddOrUpdateCart(orderID, bookID, quantity, subtotal)

// 		assert.NoError(t, err)
// 		assert.NoError(t, mock.ExpectationsWereMet())
// 	})

// 	t.Run("error when adding to cart", func(t *testing.T) {
// 		orderID := 1
// 		bookID := 1
// 		quantity := 2
// 		subtotal := int64(5000)

// 		mock.ExpectBegin()

// 		mock.ExpectExec("INSERT INTO order_details").
// 			WithArgs(orderID, bookID, quantity, subtotal).
// 			WillReturnError(sql.ErrConnDone)

// 		mock.ExpectRollback()

// 		err := orderRepo.AddOrUpdateCart(orderID, bookID, quantity, subtotal)

// 		assert.Error(t, err)
// 		assert.EqualError(t, err, sql.ErrConnDone.Error())
// 		assert.NoError(t, mock.ExpectationsWereMet())
// 	})
// }

// func TestOrderRepository_GetCart(t *testing.T) {
// 	db, mock, err := sqlmock.New()
// 	assert.NoError(t, err)
// 	defer db.Close()

// 	orderRepo := repository.NewOrderRepository(db)

// 	t.Run("successful get cart", func(t *testing.T) {
// 		orderID := 1

// 		mock.ExpectQuery("SELECT o.id, o.total").
// 			WithArgs(orderID).
// 			WillReturnRows(sqlmock.NewRows([]string{"id", "total", "detail_id", "book_id", "quantity", "subtotal", "title", "author", "price"}).
// 				AddRow(orderID, 5000, 1, 1, 2, 2500, "Book Title", "Author Name", 2500))

// 		utilsMock := utils.ConvertToDetailResponse

// 		utils.ConvertToDetailResponse = func(rows *sql.Rows) ([]model.OrderResponse, error) {
// 			return []model.OrderResponse{
// 				{
// 					ID:    int64(orderID),
// 					Total: 5000,
// 					OrderDetail: []model.OrderDetailResponse{
// 						{
// 							DetailID: 1,
// 							BookID:   1,
// 							Quantity: 2,
// 							Subtotal: 2500,
// 							Title:    "Book Title",
// 							Author:   "Author Name",
// 							Price:    2500,
// 						},
// 					},
// 				},
// 			}, nil
// 		}

// 		cart, err := orderRepo.GetCart(orderID)

// 		assert.NoError(t, err)
// 		assert.Equal(t, orderID, int(cart.ID))
// 		assert.NoError(t, mock.ExpectationsWereMet())

// 		utils.ConvertToDetailResponse = utilsMock
// 	})

// 	t.Run("cart not found", func(t *testing.T) {
// 		orderID := 1

// 		mock.ExpectQuery("SELECT o.id, o.total").
// 			WithArgs(orderID).
// 			WillReturnRows(sqlmock.NewRows([]string{"id", "total", "detail_id", "book_id", "quantity", "subtotal", "title", "author", "price"}))

// 		cart, err := orderRepo.GetCart(orderID)

// 		assert.NoError(t, err)
// 		assert.Equal(t, int64(orderID), cart.ID)
// 		assert.Empty(t, cart.OrderDetail)
// 		assert.NoError(t, mock.ExpectationsWereMet())
// 	})
// }

// func TestOrderRepository_RemoveFromCart(t *testing.T) {
// 	db, mock, err := sqlmock.New()
// 	assert.NoError(t, err)
// 	defer db.Close()

// 	orderRepo := repository.NewOrderRepository(db)

// 	t.Run("successful remove from cart", func(t *testing.T) {
// 		orderID := 1
// 		bookID := 1

// 		mock.ExpectBegin()

// 		mock.ExpectExec("DELETE FROM order_details").
// 			WithArgs(orderID, bookID).
// 			WillReturnResult(sqlmock.NewResult(1, 1))

// 		mock.ExpectExec("WITH subtotal_sum AS").
// 			WithArgs(orderID).
// 			WillReturnResult(sqlmock.NewResult(1, 1))

// 		mock.ExpectCommit()

// 		err := orderRepo.RemoveFromCart(orderID, bookID)

// 		assert.NoError(t, err)
// 		assert.NoError(t, mock.ExpectationsWereMet())
// 	})

// 	t.Run("error removing from cart", func(t *testing.T) {
// 		orderID := 1
// 		bookID := 1

// 		mock.ExpectBegin()

// 		mock.ExpectExec("DELETE FROM order_details").
// 			WithArgs(orderID, bookID).
// 			WillReturnError(sql.ErrConnDone)

// 		mock.ExpectRollback()

// 		err := orderRepo.RemoveFromCart(orderID, bookID)

// 		assert.Error(t, err)
// 		assert.EqualError(t, err, sql.ErrConnDone.Error())
// 		assert.NoError(t, mock.ExpectationsWereMet())
// 	})
// }

// func TestOrderRepository_GetOrderHistory(t *testing.T) {
// 	db, mock, err := sqlmock.New()
// 	assert.NoError(t, err)
// 	defer db.Close()

// 	orderRepo := repository.NewOrderRepository(db)

// 	t.Run("successful get order history", func(t *testing.T) {
// 		customerID := 1

// 		mock.ExpectQuery("SELECT o.id, o.total").
// 			WithArgs(customerID).
// 			WillReturnRows(sqlmock.NewRows([]string{"id", "total", "detail_id", "book_id", "quantity", "subtotal", "title", "author", "price"}).
// 				AddRow(1, 5000, 1, 1, 2, 2500, "Book Title", "Author Name", 2500))

// 		utilsMock := utils.ConvertToDetailResponse

// 		utils.ConvertToDetailResponse = func(rows *sql.Rows) ([]model.OrderResponse, error) {
// 			return []model.OrderResponse{
// 				{
// 					ID:    1,
// 					Total: 5000,
// 					OrderDetail: []model.OrderDetailResponse{
// 						{
// 							DetailID: 1,
// 							BookID:   1,
// 							Quantity: 2,
// 							Subtotal: 2500,
// 							Title:    "Book Title",
// 							Author:   "Author Name",
// 							Price:    2500,
// 						},
// 					},
// 				},
// 			}, nil
// 		}

// 		orderHistory, err := orderRepo.GetOrderHistory(customerID)

// 		assert.NoError(t, err)
// 		assert.NotNil(t, orderHistory)
// 		assert.NoError(t, mock.ExpectationsWereMet())

// 		utils.ConvertToDetailResponse = utilsMock
// 	})
// }

// func TestOrderRepository_PayOrder(t *testing.T) {
// 	db, mock, err := sqlmock.New()
// 	assert.NoError(t, err)
// 	defer db.Close()

// 	orderRepo := repository.NewOrderRepository(db)

// 	t.Run("successful pay order", func(t *testing.T) {
// 		customerID := 1

// 		mock.ExpectBegin()

// 		mock.ExpectExec("UPDATE orders SET order_state = 2").
// 			WithArgs(customerID).
// 			WillReturnResult(sqlmock.NewResult(1, 1))

// 		mock.ExpectCommit()

// 		err := orderRepo.PayOrder(customerID)

// 		assert.NoError(t, err)
// 		assert.NoError(t, mock.ExpectationsWereMet())
// 	})

// 	t.Run("error paying order", func(t *testing.T) {
// 		customerID := 1

// 		mock.ExpectBegin()

// 		mock.ExpectExec("UPDATE orders SET order_state = 2").
// 			WithArgs(customerID).
// 			WillReturnError(sql.ErrConnDone)

// 		mock.ExpectRollback()

// 		err := orderRepo.PayOrder(customerID)

// 		assert.Error(t, err)
// 		assert.EqualError(t, err, sql.ErrConnDone.Error())
// 		assert.NoError(t, mock.ExpectationsWereMet())
// 	})
// }
