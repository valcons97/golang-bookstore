package service

import (
	"bookstore/internal/model"
	"bookstore/internal/repository"
)

type BookService interface {
	CreateBook(book *model.Book) (int64, error)
	GetBooks() ([]model.Book, error)
	GetBookById(id int) (*model.Book, error)
	UpdateBook(book *model.Book) error
}

type bookService struct {
	repository repository.BookRepository
}

func NewBookService(repository repository.BookRepository) BookService {
	return &bookService{repository: repository}
}

// CreateBook implements Service.
func (s *bookService) CreateBook(book *model.Book) (int64, error) {
	return s.repository.CreateBook(book)
}

// GetBookById implements Service.
func (s *bookService) GetBookById(id int) (*model.Book, error) {
	return s.repository.GetBookById(id)
}

// GetBooks implements Service.
func (s *bookService) GetBooks() ([]model.Book, error) {
	return s.repository.GetBooks()
}

// UpdateBook implements Service.
func (s *bookService) UpdateBook(book *model.Book) error {
	return s.repository.UpdateBook(book)
}
