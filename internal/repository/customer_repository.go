package repository

import (
	"bookstore/internal/model"
	"bookstore/pkg/utils"
	"database/sql"
	"log"
	"strings"
)

type CustomerRepository interface {
	Register(customer *model.Customer) error
	Login(email, password string) (*model.Customer, error)
}

type customerRepository struct {
	db *sql.DB
}

func NewCustomerRepository(db *sql.DB) CustomerRepository {
	return &customerRepository{db: db}
}

// Register implements CustomerRepository.
func (c *customerRepository) Register(customer *model.Customer) error {

	query := "INSERT INTO customers (email, password, name, address)  VALUES ($1, $2, $3, $4)"
	_, err := c.db.Exec(query, customer.Email, customer.Password, customer.Name, customer.Address)

	if err != nil {

		// checking the error for duplicate email since email is unique
		if strings.Contains(err.Error(), "23505") &&
			strings.Contains(err.Error(), "customers_email_key") {

			return utils.ErrDuplicateEmail
		}

		log.Printf("[Register] Could not register customer: %v", err)
		return err
	}

	return nil
}

// Login implements CustomerRepository.
// Returning models of customer to be checked in customer service
func (c *customerRepository) Login(email string, password string) (*model.Customer, error) {
	var customer model.Customer

	query := `SELECT id, email, password FROM customers WHERE email = $1`
	err := c.db.QueryRow(query, email).
		Scan(&customer.ID, &customer.Email, &customer.Password)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, utils.ErrEmailNotFound
		}

		log.Printf("[Login] Error getting customer from database: %v", err)
		return nil, err
	}

	return &customer, nil
}
