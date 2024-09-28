package repository

import "bookstore/internal/model"

type BookRepository interface {
	CreateBook(book *model.Book) (int, error)
	GetBooks() ([]model.Book, error)
	GetBookById(id int) (model.Book, error)
	UpdateBook(book *model.Book) error
}
