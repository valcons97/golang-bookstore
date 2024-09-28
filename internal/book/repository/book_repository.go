package repository

import (
	"bookstore/internal/book/model"
	converter "bookstore/pkg/utils"
	"database/sql"
	"errors"
	"log"
)

type BookRepository interface {
	CreateBook(book *model.Book) (int, error)
	GetBooks() ([]model.Book, error)
	GetBookById(id int) (model.Book, error)
	UpdateBook(book *model.Book) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) BookRepository {
	return &repository{db: db}
}

// Add new book to the database.
func (r *repository) CreateBook(book *model.Book) (int, error) {
	var id int

	query := "INSERT INTO books (title, author, price) VALUES ($1, $2, $3) RETURNING ID"
	err := r.db.QueryRow(query, book.Title, book.Author, converter.ConvertStorePrice(&book.Price)).
		Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// Retrieves list of all books
func (r *repository) GetBooks() ([]model.Book, error) {
	query := "SELECT id, title, author, price FROM books"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []model.Book
	for rows.Next() {
		var book model.Book
		var price int64
		err := rows.Scan(&book.ID, &book.Title, &book.Author, &price)
		if err != nil {
			return nil, err
		}
		book.Price = *converter.ConvertToDisplayPrice(&price)
		books = append(books, book)
	}
	return books, nil
}

// Retrieve a single book define by its id
func (r *repository) GetBookById(id int) (model.Book, error) {
	var book model.Book
	var price int64
	query := "SELECT id, title, author, price FROM books WHERE id = $1"
	row := r.db.QueryRow(query, id)

	err := row.Scan(&book.ID, &book.Title, &book.Author, &price)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("[GetBookById] book is not found with id : %d", id)
			return book, errors.New("book not found")
		}
		log.Printf("[GetBookById] got exception for id : %d , and e: %v", id, err)
		return book, err
	}

	book.Price = *converter.ConvertToDisplayPrice(&price)
	return book, nil
}

// UpdateBook implements Repository.
func (r *repository) UpdateBook(book *model.Book) error {
	var updateId int
	query := "UPDATE books SET title = $1, author = $2, price = $3 WHERE id = $4 RETURNING id"
	err := r.db.QueryRow(query, book.Title, book.Author, converter.ConvertStorePrice(&book.Price), book.ID).
		Scan(&updateId)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("[UpdateBook] book is not found with id : %d", book.ID)
			return errors.New("book not found")
		}
		return err
	}

	return nil
}
