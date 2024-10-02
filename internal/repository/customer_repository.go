package repository

import (
	"bookstore/internal/model"
	"bookstore/pkg/utils"
	"database/sql"
	"fmt"
	"log"
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

// Login implements CustomerRepository.
func (c *customerRepository) Login(email string, password string) (*model.Customer, error) {
	var customer model.Customer

	query := `SELECT id, email, password FROM customers WHERE email = $1`
	err := c.db.QueryRow(query, email).
		Scan(&customer.ID, &customer.Email, &customer.Password)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("[Login] Invalid email or password for email: %s", email)
			return nil, fmt.Errorf("invalid email or password")
		}
		log.Printf("[Login] Error querying customer: %v", err)
		return nil, err
	}

	if !utils.CheckPassword(password, customer.Password) {
		log.Printf("[Login] Password mismatch for email: %s", email)
		return nil, fmt.Errorf("invalid email or password")
	}

	return &customer, nil
}

// Register implements CustomerRepository.
func (c *customerRepository) Register(customer *model.Customer) error {
	var existingCustomer model.Customer

	queryCheck := "SELECT id FROM customers WHERE email = $1"

	err := c.db.QueryRow(queryCheck, customer.Email).Scan(&existingCustomer.ID)

	if err != nil && err != sql.ErrNoRows {
		log.Printf("[Register] Error checking existing customer: %v", err)
		return err
	}

	if existingCustomer.ID != 0 {
		log.Printf("[Register] Email already registered: %s", customer.Email)
		return fmt.Errorf("email registered")
	} else {
		query := "INSERT INTO customers (email, password, name, address)  VALUES ($1, $2, $3, $4)"
		_, err := c.db.Exec(query, customer.Email, customer.Password, customer.Name, customer.Address)

		if err != nil {
			log.Printf("[Register] Could not register customer: %v", err)
			return err
		}
	}

	log.Printf("[Register] Customer registered successfully: %s", customer.Email)
	return nil
}
