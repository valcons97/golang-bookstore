package repository

import (
	"bookstore/internal/model"
	"bookstore/pkg/utils"
	"database/sql"
	"errors"
	"log"
)

type customerRepository struct {
	db *sql.DB
}

func NewCustomerRepository(db *sql.DB) CustomerRepository {
	return &customerRepository{db: db}
}

// Login implements CustomerRepository.
func (c *customerRepository) Login(email string, password string) (*model.Customer, error) {
	var customer model.Customer

	// Query the database for the customer with the given email
	query := `SELECT id, email, password, name, address FROM customers WHERE email = $1`
	err := c.db.QueryRow(query, email).
		Scan(&customer.ID, &customer.Email, &customer.Password, &customer.Name, &customer.Address)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("Invalid Email or Password")
		}
		return nil, err
	}

	// Here you should compare the provided password with the stored password
	// assuming you have a function `CheckPassword` to verify the password

	log.Printf("[Login] password e: %v", password)
	log.Printf("[Login] db pass e: %v", customer.Password)

	if !utils.CheckPassword(password, customer.Password) {
		return nil, errors.New("Invalid Email or Password")
	}

	return &customer, nil
}

// Register implements CustomerRepository.
func (c *customerRepository) Register(customer *model.Customer) error {
	var existingCustomer model.Customer

	queryCheck := "SELECT id FROM customers WHERE email = $1"

	err := c.db.QueryRow(queryCheck, customer.Email).Scan(existingCustomer.ID)

	if err != nil && err != sql.ErrNoRows {
		log.Printf("[Register] something went wrong when checking e: %v", err)
		return err
	}

	log.Printf("[Register] password e: %v", customer.Password)

	if existingCustomer.ID != 0 {
		return errors.New("Email already registered")
	} else {
		query := "INSERT INTO customers (email, password, name, address)  VALUES ($1, $2, $3, $4)"
		_, err := c.db.Exec(query, customer.Email, customer.Password, customer.Name, customer.Address)

		if err != nil {
			log.Printf("[Register] Could not register customer with e: %v", err)
			return err
		}

	}

	return nil
}
