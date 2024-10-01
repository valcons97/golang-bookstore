package service

import (
	"bookstore/internal/model"
	"bookstore/internal/repository"
	"log"
)

type OrderService interface {
	AddToCart(orderID, bookID, quantity int, subtotal int64) error
	GetCart(orderId int) (model.OrderResponse, error)
	GetOrderHistory(customerID, limit, page int) ([]model.OrderResponse, error)
	CreateOrderIfNotExists(customerID int) (int, error)
	RemoveFromCart(orderId int, bookId int) error
	PayOrder(customerID int) error
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
func (s *orderService) GetCart(orderId int) (model.OrderResponse, error) {
	return s.repository.GetCart(orderId)
}

// GetOrderHistory implements OrderService.
func (s *orderService) GetOrderHistory(customerID, limit, page int) ([]model.OrderResponse, error) {
	log.Print(customerID, limit, page)
	return s.repository.GetOrderHistory(customerID, limit, page)
}

// GetCart implements OrderService.
func (s *orderService) RemoveFromCart(orderId int, bookId int) error {
	return s.repository.RemoveFromCart(orderId, bookId)
}

// GetCart implements OrderService.
func (s *orderService) PayOrder(customerId int) error {
	return s.repository.PayOrder(customerId)
}
