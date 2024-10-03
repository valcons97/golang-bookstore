package repository

import (
	"bookstore/internal/model"
	"bookstore/pkg/utils"

	"database/sql"
	"log"
)

type BookRepository interface {
	CreateBook(book *model.Book) error
	GetBooks() ([]model.Book, error)
	GetBookById(id int) (*model.Book, error)
	UpdateBook(book *model.Book) error
}

type bookRepository struct {
	db *sql.DB
}

func NewBookRepository(db *sql.DB) BookRepository {
	return &bookRepository{db: db}
}

// Add new book to the database.
func (r *bookRepository) CreateBook(book *model.Book) error {

	query := "INSERT INTO books (title, author, price) VALUES ($1, $2, $3)"
	_, err := r.db.Query(query, book.Title, book.Author, utils.ConvertStorePrice(&book.Price))
	if err != nil {
		log.Printf("[CreateBook] Error inserting book: %v", err)
		return err
	}
	return nil
}

// Retrieves list of all books
func (r *bookRepository) GetBooks() ([]model.Book, error) {
	query := "SELECT id, title, author, price FROM books"
	rows, err := r.db.Query(query)
	if err != nil {
		// db error
		log.Printf("[GetBooks] Error retrieving list of books from database: %v", err)
		return nil, err
	}
	defer rows.Close()

	var books []model.Book
	for rows.Next() {
		var book model.Book
		var price int64
		err := rows.Scan(&book.ID, &book.Title, &book.Author, &price)
		if err != nil {
			log.Printf("[GetBooks] Error getting book: %v", err)
			return nil, err
		}
		book.Price = *utils.ConvertToDisplayPrice(&price)
		books = append(books, book)
	}
	return books, nil
}

// Retrieve a single book define by its id
func (r *bookRepository) GetBookById(id int) (*model.Book, error) {
	var book model.Book
	var price int64
	query := "SELECT id, title, author, price FROM books WHERE id = $1"
	row := r.db.QueryRow(query, id)

	err := row.Scan(&book.ID, &book.Title, &book.Author, &price)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("[GetBookById] Book not found with id: %d", id)
			return &book, utils.ErrBookNotFound
		}
		log.Printf("[GetBookById] Error retrieving book with id: %d, error: %v", id, err)
		return &book, err
	}

	book.Price = *utils.ConvertToDisplayPrice(&price)
	return &book, nil
}

// UpdateBook implements Repository.
func (r *bookRepository) UpdateBook(book *model.Book) error {
	var updateId int
	query := "UPDATE books SET title = $1, author = $2, price = $3 WHERE id = $4 RETURNING id"
	err := r.db.QueryRow(query, book.Title, book.Author, utils.ConvertStorePrice(&book.Price), book.ID).
		Scan(&updateId)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("[UpdateBook] Book not found with id: %d", book.ID)
			return utils.ErrBookNotFound
		}
		log.Printf("[UpdateBook] Error updating book with id: %d, error: %v", book.ID, err)
		return err
	}

	return nil
}
