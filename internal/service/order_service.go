package service

import (
	"bookstore/internal/model"
	"bookstore/internal/repository"
)

type OrderService interface {
	AddBookToOrder(orderDetail *model.OrderDetail) error
	GetCart(customerID int) (model.OrderResponse, error)
	GetPaidOrder(customerID int) ([]model.OrderResponse, error)
	CreateOrderIfNotExists(customerID int, total float64) (int, error)
}

type orderService struct {
	repository repository.OrderRepository
}

func NewOrderService(repository repository.OrderRepository) OrderService {
	return &orderService{repository: repository}
}

func (s *orderService) AddBookToOrder(orderDetail *model.OrderDetail) error {
	return s.repository.AddBookToOrder(orderDetail)
}

// CreateOrderIfNotExists implements OrderService.
func (s *orderService) CreateOrderIfNotExists(customerID int, total float64) (int, error) {
	return s.repository.CreateOrderIfNotExists(customerID, total)
}

// GetCart implements OrderService.
func (s *orderService) GetCart(customerID int) (model.OrderResponse, error) {
	return s.repository.GetCart(customerID)
}

// GetPaidOrder implements OrderService.
func (s *orderService) GetPaidOrder(customerID int) ([]model.OrderResponse, error) {
	return s.repository.GetPaidOrder(customerID)
}
