package repository_test

import (
	"bookstore/internal/model"
	"bookstore/internal/repository"
	"bookstore/pkg/utils"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCustomerRepository_Register(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	customerRepo := repository.NewCustomerRepository(db)
	query := "SELECT id FROM customers WHERE email = ?"

	customer := &model.Customer{
		Email:    "test@example.com",
		Password: "hashedpassword",
		Name:     "John Doe",
		Address:  "123 Street",
	}
	t.Run("successful registration", func(t *testing.T) {

		mock.ExpectQuery(query).
			WithArgs(customer.Email).
			WillReturnRows(sqlmock.NewRows([]string{"id"}))

		mock.ExpectExec("INSERT INTO customers").
			WithArgs(customer.Email, customer.Password, customer.Name, customer.Address).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := customerRepo.Register(customer)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("email already registered", func(t *testing.T) {

		mock.ExpectQuery(query).
			WithArgs(customer.Email).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		err := customerRepo.Register(customer)

		assert.Error(t, err)
		assert.EqualError(t, err, "email registered")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestCustomerRepository_Login(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	customerRepo := repository.NewCustomerRepository(db)
	query := "SELECT id, email, password, name, address FROM customers WHERE email = ?"

	t.Run("successful login", func(t *testing.T) {
		email := "test@example.com"
		password := "correctpassword"
		hashedPassword, _ := utils.HashPassword(password)

		customer := &model.Customer{
			ID:       1,
			Email:    email,
			Password: hashedPassword,
			Name:     "John Doe",
			Address:  "123 Street",
		}

		mock.ExpectQuery(query).
			WithArgs(email).
			WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password", "name", "address"}).
				AddRow(customer.ID, customer.Email, customer.Password, customer.Name, customer.Address),
			)

		result, err := customerRepo.Login(email, password)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, customer.ID, result.ID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("login invalid email", func(t *testing.T) {
		email := "wrong@example.com"
		password := "password"

		mock.ExpectQuery(query).
			WithArgs(email).
			WillReturnError(sql.ErrNoRows)

		result, err := customerRepo.Login(email, password)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("login incorrect password", func(t *testing.T) {
		email := "test@example.com"
		password := "wrongpassword"
		hashedPassword, _ := utils.HashPassword("correctpassword")

		customer := &model.Customer{
			ID:       1,
			Email:    email,
			Password: hashedPassword,
			Name:     "John Doe",
			Address:  "123 Street",
		}

		mock.ExpectQuery(query).
			WithArgs(email).
			WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password", "name", "address"}).
				AddRow(customer.ID, customer.Email, customer.Password, customer.Name, customer.Address),
			)

		result, err := customerRepo.Login(email, password)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
