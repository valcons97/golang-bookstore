package repository_test

import (
	"bookstore/internal/model"
	"bookstore/internal/repository"
	"bookstore/pkg/utils"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateBook(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewBookRepository(db)

	book := &model.Book{
		Title:  "Test Book",
		Author: "Test Author",
		Price:  9.99,
	}

	query := "INSERT INTO books \\(title, author, price\\) VALUES \\(\\$1, \\$2, \\$3\\) RETURNING ID"
	mock.ExpectQuery(query).
		WithArgs(book.Title, book.Author, utils.ConvertStorePrice(&book.Price)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	id, err := repo.CreateBook(book)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), id)
}

func TestGetBooks(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewBookRepository(db)

	rows := sqlmock.NewRows([]string{"id", "title", "author", "price"}).
		AddRow(1, "Test Book", "Test Author", 999).
		AddRow(2, "Another Book", "Another Author", 1299)

	query := "SELECT id, title, author, price FROM books"
	mock.ExpectQuery(query).WillReturnRows(rows)

	books, err := repo.GetBooks()
	assert.NoError(t, err)
	assert.Len(t, books, 2)

	assert.Equal(t, "Test Book", books[0].Title)
	assert.Equal(t, "Another Book", books[1].Title)
}

func TestGetBookById(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewBookRepository(db)

	bookID := 1
	row := sqlmock.NewRows([]string{"id", "title", "author", "price"}).
		AddRow(1, "Test Book", "Test Author", 999)

	query := "SELECT id, title, author, price FROM books WHERE id = \\$1"
	mock.ExpectQuery(query).WithArgs(bookID).WillReturnRows(row)

	book, err := repo.GetBookById(bookID)
	assert.NoError(t, err)
	assert.NotNil(t, book)
	assert.Equal(t, "Test Book", book.Title)
}

func TestUpdateBook(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewBookRepository(db)

	book := &model.Book{
		ID:     1,
		Title:  "Updated Title",
		Author: "Updated Author",
		Price:  15.99,
	}

	query := "UPDATE books SET title = \\$1, author = \\$2, price = \\$3 WHERE id = \\$4 RETURNING id"
	mock.ExpectQuery(query).
		WithArgs(book.Title, book.Author, utils.ConvertStorePrice(&book.Price), book.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	err = repo.UpdateBook(book)
	assert.NoError(t, err)
}
