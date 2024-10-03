package service

import (
	"bookstore/internal/model"
	"bookstore/internal/repository"
	"bookstore/pkg/utils"
	"log"
	"strings"
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

	if email == "" || password == "" {
		return "", utils.ErrEmptyEmailOrPassword
	}

	customer, err := s.repository.Login(strings.ToLower(email), password)

	if err != nil {
		return "", err
	}

	if !utils.CheckPassword(password, customer.Password) {
		return "", utils.ErrWrongPassword
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

func (s *customerService) Register(customer *model.Customer) error {
	hashedPassword, err := utils.HashPassword(customer.Password)
	if err != nil {
		log.Printf("[Register] failed to hash for email: %s e: %v", customer.Email, err)
		return err
	}

	customer.Email = strings.ToLower(customer.Email)
	customer.Password = hashedPassword

	return s.repository.Register(customer)
}
