package service_test

import (
	"bookstore/internal/handler/request"
	"bookstore/internal/model"

	"bookstore/internal/service"
	"bookstore/test/mocks"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestAddToCart(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockOrderRepository(ctrl)
	orderService := service.NewOrderService(mockRepo)

	customerID := 1
	request := request.AddToCartRequest{
		BookId:   1,
		Quantity: 2,
		Price:    10.0,
	}

	t.Run("Success", func(t *testing.T) {
		orderID := 1
		mockRepo.EXPECT().CreateOrderIfNotExists(customerID).Return(orderID, nil)
		mockRepo.EXPECT().
			AddOrUpdateCart(orderID, int(request.BookId), int(request.Quantity), gomock.Any()).
			Return(nil)

		err := orderService.AddToCart(customerID, request)

		assert.NoError(t, err)
	})

	t.Run("Error creating order", func(t *testing.T) {
		mockRepo.EXPECT().CreateOrderIfNotExists(customerID).Return(0, errors.New("creation error"))

		err := orderService.AddToCart(customerID, request)

		assert.Error(t, err)
		assert.EqualError(t, err, "creation error")
	})

	t.Run("Error adding to cart", func(t *testing.T) {
		orderID := 1
		mockRepo.EXPECT().CreateOrderIfNotExists(customerID).Return(orderID, nil)
		mockRepo.EXPECT().
			AddOrUpdateCart(orderID, int(request.BookId), int(request.Quantity), gomock.Any()).
			Return(errors.New("add error"))

		err := orderService.AddToCart(customerID, request)

		assert.Error(t, err)
		assert.EqualError(t, err, "add error")
	})
}

func TestGetCart(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockOrderRepository(ctrl)
	orderService := service.NewOrderService(mockRepo)

	customerID := 1

	t.Run("Success", func(t *testing.T) {
		orderID := 1
		detailID := int64(1)
		quantity := int64(2)
		subtotal := 20.0

		expectedBook := model.Book{
			ID:     1,
			Title:  "Example Book",
			Author: "Author Name",
			Price:  10.0,
		}

		expectedOrderDetailResponse := model.OrderDetailResponse{
			ID:       detailID,
			Book:     []model.Book{expectedBook},
			Quantity: quantity,
			Subtotal: subtotal,
		}

		expectedResponse := &model.OrderResponse{
			ID: int64(orderID),
			OrderDetail: []model.OrderDetailResponse{
				expectedOrderDetailResponse,
			},
			Total: subtotal,
		}

		mockRepo.EXPECT().CreateOrderIfNotExists(customerID).Return(orderID, nil)
		mockRepo.EXPECT().GetCart(orderID).Return(expectedResponse, nil)

		response, err := orderService.GetCart(customerID)

		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, response)
	})

	t.Run("Cart Empty", func(t *testing.T) {
		orderID := 1

		mockRepo.EXPECT().CreateOrderIfNotExists(customerID).Return(orderID, nil)
		mockRepo.EXPECT().GetCart(orderID).Return(&model.OrderResponse{
			ID:          int64(orderID),
			OrderDetail: []model.OrderDetailResponse{},
			Total:       0.0,
		}, nil)

		response, err := orderService.GetCart(customerID)

		assert.Error(t, err)
		assert.EqualError(t, err, "cart empty")
		assert.Nil(t, response)
	})

	t.Run("Error creating order", func(t *testing.T) {
		mockRepo.EXPECT().CreateOrderIfNotExists(customerID).Return(0, errors.New("creation error"))

		response, err := orderService.GetCart(customerID)

		assert.Error(t, err)
		assert.Nil(t, response)
	})

	t.Run("Error getting cart", func(t *testing.T) {
		orderID := 1
		mockRepo.EXPECT().CreateOrderIfNotExists(customerID).Return(orderID, nil)
		mockRepo.EXPECT().GetCart(orderID).Return(nil, errors.New("get cart error"))

		response, err := orderService.GetCart(customerID)

		assert.Error(t, err)
		assert.Nil(t, response)
	})
}

func TestGetOrderHistory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockOrderRepository(ctrl)
	orderService := service.NewOrderService(mockRepo)

	customerID := 1
	request := request.HistoryRequest{Limit: 10, Page: 1}

	t.Run("Success", func(t *testing.T) {
		expectedResponse := []model.OrderResponse{
			{ID: 1},
			{ID: 2},
		}

		mockRepo.EXPECT().
			GetOrderHistory(customerID, request.Limit, request.Page).
			Return(expectedResponse, nil)

		response, err := orderService.GetOrderHistory(customerID, request)

		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, response)
	})

	t.Run("Error", func(t *testing.T) {
		mockRepo.EXPECT().
			GetOrderHistory(customerID, request.Limit, request.Page).
			Return(nil, errors.New("history error"))

		response, err := orderService.GetOrderHistory(customerID, request)

		assert.Error(t, err)
		assert.Nil(t, response)
	})
}

func TestRemoveFromCart(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockOrderRepository(ctrl)
	orderService := service.NewOrderService(mockRepo)

	customerID := 1
	bookID := 1

	t.Run("Success", func(t *testing.T) {
		orderID := 1

		mockRepo.EXPECT().CreateOrderIfNotExists(customerID).Return(orderID, nil)
		mockRepo.EXPECT().RemoveFromCart(orderID, bookID).Return(nil)

		err := orderService.RemoveFromCart(customerID, bookID)

		assert.NoError(t, err)
	})

	t.Run("Error creating order", func(t *testing.T) {
		mockRepo.EXPECT().CreateOrderIfNotExists(customerID).Return(0, errors.New("creation error"))

		err := orderService.RemoveFromCart(customerID, bookID)

		assert.Error(t, err)
		assert.EqualError(t, err, "creation error")
	})

	t.Run("Error removing from cart", func(t *testing.T) {
		orderID := 1

		mockRepo.EXPECT().CreateOrderIfNotExists(customerID).Return(orderID, nil)
		mockRepo.EXPECT().RemoveFromCart(orderID, bookID).Return(errors.New("remove error"))

		err := orderService.RemoveFromCart(customerID, bookID)

		assert.Error(t, err)
		assert.EqualError(t, err, "remove error")
	})
}

func TestPayOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockOrderRepository(ctrl)
	orderService := service.NewOrderService(mockRepo)

	customerID := 1

	t.Run("Success", func(t *testing.T) {
		mockRepo.EXPECT().PayOrder(customerID).Return(nil)

		err := orderService.PayOrder(customerID)

		assert.NoError(t, err)
	})

	t.Run("Error", func(t *testing.T) {
		mockRepo.EXPECT().PayOrder(customerID).Return(errors.New("payment error"))

		err := orderService.PayOrder(customerID)

		assert.Error(t, err)
		assert.EqualError(t, err, "payment error")
	})
}
