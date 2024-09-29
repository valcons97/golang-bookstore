package service

import (
	"bookstore/internal/model"
	"bookstore/internal/repository"
)

type OrderService interface {
	AddToCart(orderID, bookID, quantity int, subtotal int64) error
	GetCart(customerID int) (model.OrderResponse, error)
	GetPaidOrder(customerID int) ([]model.OrderResponse, error)
	CreateOrderIfNotExists(customerID int) (int, error)
}

type orderService struct {
	repository repository.OrderRepository
}

func NewOrderService(repository repository.OrderRepository) OrderService {
	return &orderService{repository: repository}
}

func (s *orderService) AddToCart(orderID, bookID, quantity int, subtotal int64) error {
	return s.repository.AddOrUpdateCart(orderID, bookID, quantity, subtotal)
}

// CreateOrderIfNotExists implements OrderService.
func (s *orderService) CreateOrderIfNotExists(customerID int) (int, error) {
	return s.repository.CreateOrderIfNotExists(customerID)
}

// GetCart implements OrderService.
func (s *orderService) GetCart(customerID int) (model.OrderResponse, error) {
	return s.repository.GetCart(customerID)
}

// GetPaidOrder implements OrderService.
func (s *orderService) GetPaidOrder(customerID int) ([]model.OrderResponse, error) {
	return s.repository.GetPaidOrder(customerID)
}
