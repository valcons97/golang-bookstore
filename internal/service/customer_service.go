package service

import (
	"bookstore/internal/model"
	"bookstore/internal/repository"
	"bookstore/pkg/utils"
	"log"
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
		log.Printf("[Login] error in login service for email: %s, e: %v", email, err)
		return "", err
	}

	token, err := utils.GenerateToken(
		customer.ID,
		customer.Email,
	)

	if err != nil {
		log.Printf("[Login] error generating token for email: %s, e: %v", email, err)
		return "", err
	}

	return token, nil
}

// Register implements CustomerService.
func (s *customerService) Register(customer *model.Customer) error {
	return s.repository.Register(customer)
}
