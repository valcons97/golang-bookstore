package service

import (
	"bookstore/internal/model"
	"bookstore/internal/repository"
)

type customerService struct {
	repository repository.CustomerRepository
}

func NewCustomerService(repository repository.CustomerRepository) CustomerService {
	return &customerService{repository: repository}
}

// Login implements CustomerService.
func (s *customerService) Login(email string, password string) (*model.Customer, error) {
	return s.repository.Login(email, password)
}

// Register implements CustomerService.
func (s *customerService) Register(customer *model.Customer) error {
	return s.repository.Register(customer)
}
