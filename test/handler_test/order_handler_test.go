package handler_test

import (
	"bookstore/internal/handler"
	"bookstore/internal/model"
	"bookstore/pkg/utils"
	"bookstore/test/mocks"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestOrderHandler_PayOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrderService := mocks.NewMockOrderService(ctrl)
	router := gin.Default()

	orderHandler := handler.NewOrderHandler(mockOrderService)
	router.POST("/pay", orderHandler.PayOrder)

	t.Run("success", func(t *testing.T) {
		customerID := int64(1)
		mockOrderService.EXPECT().PayOrder(int(customerID)).Return(nil)

		token, _ := utils.GenerateToken(customerID, "test@example.com")
		req, _ := http.NewRequest(http.MethodPost, "/pay", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var actualResponse gin.H
		err := json.Unmarshal(w.Body.Bytes(), &actualResponse)
		assert.NoError(t, err)
		assert.Equal(t, "Order paid successfully", actualResponse["message"])
	})

	t.Run("unauthorized", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "/pay", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestOrderHandler_GetCart(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrderService := mocks.NewMockOrderService(ctrl)
	router := gin.Default()

	orderHandler := handler.NewOrderHandler(mockOrderService)
	router.GET("/cart", orderHandler.GetCart)

	t.Run("success", func(t *testing.T) {
		customerID := int64(1)
		orderID := int64(1)
		mockOrderService.EXPECT().CreateOrderIfNotExists(int(customerID)).Return(int(orderID), nil)

		books := []model.Book{
			{ID: 1, Title: "Book 1", Author: "Author 1", Price: 10.0},
			{ID: 2, Title: "Book 2", Author: "Author 2", Price: 15.0},
		}

		orderDetails := []model.OrderDetailResponse{
			{
				ID:       1,
				Book:     books,
				Quantity: 2,
				Subtotal: 25.0,
			},
		}

		expectedResponse := model.OrderResponse{
			ID:          orderID,
			OrderDetail: orderDetails,
			Total:       25.0,
		}

		mockOrderService.EXPECT().
			GetCart(int(orderID)).
			Return(expectedResponse, nil)

		token, _ := utils.GenerateToken(customerID, "test@example.com")
		req, _ := http.NewRequest(http.MethodGet, "/cart", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var actualResponse model.OrderResponse
		err := json.Unmarshal(w.Body.Bytes(), &actualResponse)
		assert.NoError(t, err)

		assert.Equal(t, expectedResponse, actualResponse)
	})

	t.Run("unauthorized", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/cart", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestOrderHandler_GetOrderHistory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrderService := mocks.NewMockOrderService(ctrl)
	router := gin.Default()

	orderHandler := handler.NewOrderHandler(mockOrderService)
	router.POST("/history", orderHandler.GetOrderHistory)

	page, limit := 0, 10

	t.Run("success", func(t *testing.T) {
		customerID := int64(1)

		books := []model.Book{
			{ID: 1, Title: "Book 1", Author: "Author 1", Price: 10.0},
			{ID: 2, Title: "Book 2", Author: "Author 2", Price: 15.0},
		}

		request := handler.HistoryRequest{Page: page, Limit: limit}
		jsonReq, _ := json.Marshal(request)

		orderDetails := []model.OrderDetailResponse{
			{
				ID:       1,
				Book:     books,
				Quantity: 2,
				Subtotal: 25.0,
			},
		}

		expectedResponse := []model.OrderResponse{
			{
				ID:          1,
				OrderDetail: orderDetails,
				Total:       25.0,
			},
		}

		mockOrderService.EXPECT().
			GetOrderHistory(int(customerID), limit, page).
			Return(expectedResponse, nil)

		token, _ := utils.GenerateToken(customerID, "test@example.com")
		req, _ := http.NewRequest(http.MethodPost, "/history", bytes.NewBuffer(jsonReq))
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var actualResponse []model.OrderResponse
		err := json.Unmarshal(w.Body.Bytes(), &actualResponse)
		assert.NoError(t, err)

		assert.Equal(t, expectedResponse, actualResponse)
	})

	t.Run("unauthorized", func(t *testing.T) {
		request := handler.HistoryRequest{Page: page, Limit: limit}
		jsonReq, _ := json.Marshal(request)
		req, _ := http.NewRequest(http.MethodPost, "/history", bytes.NewBuffer(jsonReq))
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestOrderHandler_RemoveFromCart(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrderService := mocks.NewMockOrderService(ctrl)
	router := gin.Default()

	orderHandler := handler.NewOrderHandler(mockOrderService)
	router.POST("/remove", orderHandler.RemoveFromCart)

	t.Run("success", func(t *testing.T) {
		customerID := 1
		orderID := 1
		request := handler.RemoveItemFromCartRequest{BookId: 1}
		mockOrderService.EXPECT().CreateOrderIfNotExists(int(customerID)).Return(orderID, nil)
		mockOrderService.EXPECT().RemoveFromCart(orderID, int(request.BookId)).Return(nil)

		token, _ := utils.GenerateToken(
			int64(customerID),
			"test@example.com",
		)
		jsonReq, _ := json.Marshal(request)
		req, _ := http.NewRequest(http.MethodPost, "/remove", bytes.NewBuffer(jsonReq))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var actualResponse gin.H
		err := json.Unmarshal(w.Body.Bytes(), &actualResponse)
		assert.NoError(t, err)
		assert.Equal(t, "Book removed from cart", actualResponse["message"])
	})

	t.Run("unauthorized", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "/remove", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestOrderHandler_AddToCart(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrderService := mocks.NewMockOrderService(ctrl)
	router := gin.Default()

	orderHandler := handler.NewOrderHandler(mockOrderService)
	router.POST("/add", orderHandler.AddToCart)

	t.Run("success", func(t *testing.T) {
		customerID := 1
		orderID := 1
		request := handler.AddToCartRequest{BookId: 1, Quantity: 2, Price: 10}
		subtotal := float64(request.Quantity) * request.Price

		// since its converted to cent
		subtotalInt64 := int64(subtotal) * 100
		mockOrderService.EXPECT().CreateOrderIfNotExists(customerID).Return(orderID, nil)
		mockOrderService.EXPECT().
			AddToCart(orderID, int(request.BookId), int(request.Quantity), subtotalInt64).
			Return(nil)

		token, _ := utils.GenerateToken(
			int64(customerID),
			"test@example.com",
		)
		jsonReq, _ := json.Marshal(request)
		req, _ := http.NewRequest(http.MethodPost, "/add", bytes.NewBuffer(jsonReq))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var actualResponse gin.H
		err := json.Unmarshal(w.Body.Bytes(), &actualResponse)
		assert.NoError(t, err)
		assert.Equal(t, "Cart updated", actualResponse["message"])
	})

	t.Run("unauthorized", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "/add", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("bad request", func(t *testing.T) {
		token, _ := utils.GenerateToken(1, "test@example.com")
		req, _ := http.NewRequest(
			http.MethodPost,
			"/add",
			bytes.NewBuffer([]byte("{bad json}")),
		)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
