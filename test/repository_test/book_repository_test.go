package repository_test

import (
	"bookstore/internal/model"
	"bookstore/internal/repository"
	"bookstore/pkg/utils"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateBook(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	bookRepo := repository.NewBookRepository(db)

	t.Run("success", func(t *testing.T) {
		book := &model.Book{Title: "Test Book", Author: "Author", Price: 10.5}
		mock.ExpectQuery("INSERT INTO books").
			WithArgs(book.Title, book.Author, utils.ConvertStorePrice(&book.Price)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		id, err := bookRepo.CreateBook(book)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), id)
	})

	t.Run("error on insert", func(t *testing.T) {
		book := &model.Book{Title: "Test Book", Author: "Author", Price: 10.5}
		mock.ExpectQuery("INSERT INTO books").
			WithArgs(book.Title, book.Author, utils.ConvertStorePrice(&book.Price)).
			WillReturnError(errors.New("insert error"))

		id, err := bookRepo.CreateBook(book)
		assert.Error(t, err)
		assert.Equal(t, int64(0), id)
	})
}

func TestGetBooks(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	bookRepo := repository.NewBookRepository(db)

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "title", "author", "price"}).
			AddRow(1, "Book 1", "Author 1", 1000).
			AddRow(2, "Book 2", "Author 2", 2000)
		mock.ExpectQuery("SELECT id, title, author, price FROM books").
			WillReturnRows(rows)

		books, err := bookRepo.GetBooks()
		assert.NoError(t, err)
		assert.Len(t, books, 2)
	})

	t.Run("error on query", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, title, author, price FROM books").
			WillReturnError(errors.New("query error"))

		books, err := bookRepo.GetBooks()
		assert.Error(t, err)
		assert.Nil(t, books)
	})
}

func TestGetBookById(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	bookRepo := repository.NewBookRepository(db)

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "title", "author", "price"}).
			AddRow(1, "Test Book", "Author", 1000)
		mock.ExpectQuery("SELECT id, title, author, price FROM books WHERE id =").
			WithArgs(1).
			WillReturnRows(rows)

		book, err := bookRepo.GetBookById(1)
		assert.NoError(t, err)
		assert.Equal(t, "Test Book", book.Title)
	})

	t.Run("book not found", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, title, author, price FROM books WHERE id =").
			WithArgs(1).
			WillReturnError(sql.ErrNoRows)

		book, err := bookRepo.GetBookById(1)
		assert.Error(t, err)
		assert.Equal(t, "book not found", err.Error())
		assert.Equal(t, int64(0), book.ID)
	})

	t.Run("error on query", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, title, author, price FROM books WHERE id =").
			WithArgs(1).
			WillReturnError(errors.New("query error"))

		book, err := bookRepo.GetBookById(1)
		assert.Error(t, err)
		assert.Equal(t, int64(0), book.ID)
	})
}

func TestUpdateBook(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	bookRepo := repository.NewBookRepository(db)

	t.Run("success", func(t *testing.T) {
		book := &model.Book{ID: 1, Title: "Updated Book", Author: "Updated Author", Price: 20.0}
		mock.ExpectQuery("UPDATE books SET").
			WithArgs(book.Title, book.Author, utils.ConvertStorePrice(&book.Price), book.ID).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		err := bookRepo.UpdateBook(book)
		assert.NoError(t, err)
	})

	t.Run("book not found", func(t *testing.T) {
		book := &model.Book{ID: 1, Title: "Updated Book", Author: "Updated Author", Price: 20.0}
		mock.ExpectQuery("UPDATE books SET").
			WithArgs(book.Title, book.Author, utils.ConvertStorePrice(&book.Price), book.ID).
			WillReturnError(sql.ErrNoRows)

		err := bookRepo.UpdateBook(book)
		assert.Error(t, err)
		assert.Equal(t, "book not found", err.Error())
	})

	t.Run("error on update", func(t *testing.T) {
		book := &model.Book{ID: 1, Title: "Updated Book", Author: "Updated Author", Price: 20.0}
		mock.ExpectQuery("UPDATE books SET").
			WithArgs(book.Title, book.Author, utils.ConvertStorePrice(&book.Price), book.ID).
			WillReturnError(errors.New("update error"))

		err := bookRepo.UpdateBook(book)
		assert.Error(t, err)
	})
}
