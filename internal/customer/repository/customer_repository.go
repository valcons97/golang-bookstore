package repository

import "bookstore/internal/customer/model"

type CustomerRepository interface {
	Register(customer *model.Customer) error
	Login(email, password string) (*model.Customer, error)
}
