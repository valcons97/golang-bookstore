package service

import (
	"bookstore/internal/handler/request"
	"bookstore/internal/model"
	"bookstore/internal/repository"
	"bookstore/pkg/utils"
	"errors"
)

type OrderService interface {
	AddToCart(customerID int, request request.AddToCartRequest) error
	GetCart(customerID int) (*model.OrderResponse, error)
	GetOrderHistory(customerID int, request request.HistoryRequest) ([]model.OrderResponse, error)
	CreateOrderIfNotExists(customerID int) (int, error)
	RemoveFromCart(customerID int, bookId int) error
	PayOrder(customerID int) error
}

type orderService struct {
	repository repository.OrderRepository
}

func NewOrderService(repository repository.OrderRepository) OrderService {
	return &orderService{repository: repository}
}

func (s *orderService) AddToCart(customerID int, request request.AddToCartRequest) error {

	orderId, err := s.CreateOrderIfNotExists(customerID)

	if err != nil {
		return err
	}

	subTotal := request.Price * float64(request.Quantity)

	return s.repository.AddOrUpdateCart(
		orderId,
		int(request.BookId),
		int(request.Quantity),
		*utils.ConvertStorePrice(&subTotal),
	)
}

func (s *orderService) CreateOrderIfNotExists(customerID int) (int, error) {
	return s.repository.CreateOrderIfNotExists(customerID)
}

func (s *orderService) GetCart(customerID int) (*model.OrderResponse, error) {
	orderId, err := s.CreateOrderIfNotExists(customerID)

	if err != nil {
		return nil, err
	}

	cart, err := s.repository.GetCart(orderId)

	if err != nil {
		return nil, err
	}
	if cart == nil || len(cart.OrderDetail) == 0 {
		return nil, errors.New("cart empty")
	}

	return cart, nil
}

func (s *orderService) GetOrderHistory(
	customerID int,
	request request.HistoryRequest,
) ([]model.OrderResponse, error) {

	if request.Limit == 0 {
		request.Limit = 10
	}

	orders, err := s.repository.GetOrderHistory(customerID, request.Limit, request.Page)

	if err != nil {
		return nil, err
	}

	if len(orders) == 0 {
		return nil, utils.WarnCartEmpty
	}

	return orders, nil
}

func (s *orderService) RemoveFromCart(customerID int, bookId int) error {
	orderId, err := s.CreateOrderIfNotExists(customerID)
	if err != nil {
		return err
	}

	return s.repository.RemoveFromCart(orderId, bookId)
}

func (s *orderService) PayOrder(customerId int) error {
	return s.repository.PayOrder(customerId)
}
