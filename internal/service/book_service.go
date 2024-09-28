package service

import "bookstore/internal/model"

type BookService interface {
	CreateBook(book *model.Book) (int, error)
	GetBooks() ([]model.Book, error)
	GetBookById(id int) (model.Book, error)
	UpdateBook(book *model.Book) error
}
