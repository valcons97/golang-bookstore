package service

import "bookstore/internal/model"

type CustomerService interface {
	Register(customer *model.Customer) error
	Login(email, password string) (*model.Customer, error)
}
