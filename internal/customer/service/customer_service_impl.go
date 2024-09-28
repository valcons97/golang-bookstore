package service

import (
	"bookstore/internal/customer/model"
	"bookstore/internal/customer/repository"
)

type service struct {
	repository repository.CustomerRepository
}

func NewCustomerRepository(repository repository.CustomerRepository) CustomerService {
	return &service{repository: repository}
}

// Login implements CustomerService.
func (s *service) Login(email string, password string) (*model.Customer, error) {
	panic("unimplemented")
}

// Register implements CustomerService.
func (s *service) Register(customer *model.Customer) error {
	panic("unimplemented")
}
