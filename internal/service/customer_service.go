package service

import (
	"bookstore/internal/model"
	"bookstore/internal/repository"
	"bookstore/pkg/utils"
)

type CustomerService interface {
	Register(customer *model.Customer) error
	Login(email, password string) (string, error)
}

type customerService struct {
	repository repository.CustomerRepository
}

func NewCustomerService(repository repository.CustomerRepository) CustomerService {
	return &customerService{repository: repository}
}

// Login implements CustomerService.
func (s *customerService) Login(email string, password string) (string, error) {
	customer, err := s.repository.Login(email, password)

	if err != nil {
		return "", err
	}

	token, err := utils.GenerateToken(
		customer.ID,
		customer.Email,
	)

	if err != nil {
		return "", err
	}

	return token, nil
}

// Register implements CustomerService.
func (s *customerService) Register(customer *model.Customer) error {
	return s.repository.Register(customer)
}
