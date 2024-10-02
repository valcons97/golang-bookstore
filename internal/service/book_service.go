package service

import (
	"bookstore/internal/model"
	"bookstore/internal/repository"
)

type BookService interface {
	CreateBook(book *model.Book) (*model.Book, error)
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
func (s *bookService) CreateBook(book *model.Book) (*model.Book, error) {
	id, err := s.repository.CreateBook(book)

	if err != nil {
		return nil, err
	}

	book.ID = id

	return book, nil
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
