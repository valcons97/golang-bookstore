package service

import "bookstore/internal/customer/model"

type CustomerService interface {
	Register(customer *model.Customer) error
	Login(email, password string) (*model.Customer, error)
}
