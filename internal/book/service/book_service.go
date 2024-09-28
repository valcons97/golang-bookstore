package service

import (
	"bookstore/internal/book/model"
	"bookstore/internal/book/repository"
)

type BookService interface {
	CreateBook(book *model.Book) (int, error)
	GetBooks() ([]model.Book, error)
	GetBookById(id int) (model.Book, error)
	UpdateBook(book *model.Book) error
}

type service struct {
	repository repository.BookRepository
}

func NewService(repository repository.BookRepository) BookService {
	return &service{repository: repository}
}

// CreateBook implements Service.
func (s *service) CreateBook(book *model.Book) (int, error) {
	return s.repository.CreateBook(book)
}

// GetBookById implements Service.
func (s *service) GetBookById(id int) (model.Book, error) {
	return s.repository.GetBookById(id)
}

// GetBooks implements Service.
func (s *service) GetBooks() ([]model.Book, error) {
	return s.repository.GetBooks()
}

// UpdateBook implements Service.
func (s *service) UpdateBook(book *model.Book) error {
	return s.repository.UpdateBook(book)
}
