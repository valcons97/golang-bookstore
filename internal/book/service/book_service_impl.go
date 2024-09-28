package service

import (
	"bookstore/internal/book/model"
	"bookstore/internal/repository"
)

type bookService struct {
	repository repository.BookRepository
}

func NewBookService(repository repository.BookRepository) BookService {
	return &bookService{repository: repository}
}

// CreateBook implements Service.
func (s *bookService) CreateBook(book *model.Book) (int, error) {
	return s.repository.CreateBook(book)
}

// GetBookById implements Service.
func (s *bookService) GetBookById(id int) (model.Book, error) {
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
