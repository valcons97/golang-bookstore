package book

import (
	"database/sql"
)

type Repository interface {
	CreateBook(book *Book) error
	FindAll() ([]Book, error)
	FindById(id int) (Book, error)
	Update(book *Book) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

// Create adds a new book to the database.
func (r *repository) CreateBook(book *Book) error {
	query := "INSERT INTO books (title, author, price) VALUES ($1, $2, $3)"
	_, err := r.db.Exec(query, book.Title, book.Author, book.Price)
	return err
}

// FindAll implements Repository.
func (r *repository) FindAll() ([]Book, error) {
	panic("unimplemented")
}

// FindById implements Repository.
func (r *repository) FindById(id int) (Book, error) {
	panic("unimplemented")
}

// Update implements Repository.
func (r *repository) Update(book *Book) error {
	panic("unimplemented")
}
