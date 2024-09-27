package book

import (
	"database/sql"
	"errors"
	"log"
)

type Repository interface {
	CreateBook(book *Book) error
	GetBooks() ([]Book, error)
	GetBookById(id int) (Book, error)
	UpdateBook(book *Book) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

// Add new book to the database.
func (r *repository) CreateBook(book *Book) error {
	query := "INSERT INTO books (title, author, price) VALUES ($1, $2, $3)"
	_, err := r.db.Exec(query, book.Title, book.Author, book.Price)
	return err
}

// Retrieves list of all books
func (r *repository) GetBooks() ([]Book, error) {
	query := "SELECT id, title, author, price FROM books"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		var book Book
		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Price)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	return books, nil
}

// Retrieve a single book define by its id
func (r *repository) GetBookById(id int) (Book, error) {
	var book Book
	query := "SELECT id, title, author, price FROM books WHERE id = $1"
	row := r.db.QueryRow(query, id)

	err := row.Scan(&book.ID, &book.Title, &book.Author, &book.Price)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("[GetBookById] book is not found with id : %d", id)
			return book, errors.New("book not found")
		}
		log.Printf("[GetBookById] got exception for id : %d , and e: %v", id, err)
		return book, err
	}

	return book, nil
}

// UpdateBook implements Repository.
func (r *repository) UpdateBook(book *Book) error {
	query := "UPDATE books SET title = $1, author = $2, price = $3 WHERE id = $4"
	_, err := r.db.Exec(query, book.Title, book.Author, book.Price, book.ID)

	if err != nil {
		return err
	}

	return nil
}
