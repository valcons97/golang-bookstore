package service_test

import (
	"bookstore/internal/model"
	"bookstore/test/mocks"

	"bookstore/internal/service"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCreateBook(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockBookRepository(ctrl)
	bookService := service.NewBookService(mockRepo)

	book := &model.Book{Title: "Test Book", Author: "Test Author"}

	t.Run("Success", func(t *testing.T) {
		mockRepo.EXPECT().CreateBook(book).Return(int64(1), nil)

		createdBook, err := bookService.CreateBook(book)

		assert.NoError(t, err)
		assert.Equal(t, int64(1), createdBook.ID)
		assert.Equal(t, "Test Book", createdBook.Title)
	})

	t.Run("Error", func(t *testing.T) {
		mockRepo.EXPECT().CreateBook(book).Return(int64(0), errors.New("creation error"))

		createdBook, err := bookService.CreateBook(book)

		assert.Error(t, err)
		assert.Nil(t, createdBook)
	})
}

func TestGetBookById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockBookRepository(ctrl)
	bookService := service.NewBookService(mockRepo)

	t.Run("Success", func(t *testing.T) {
		book := &model.Book{ID: 1, Title: "Test Book"}
		mockRepo.EXPECT().GetBookById(1).Return(book, nil)

		result, err := bookService.GetBookById(1)

		assert.NoError(t, err)
		assert.Equal(t, book, result)
	})

	t.Run("Error", func(t *testing.T) {
		mockRepo.EXPECT().GetBookById(1).Return(nil, errors.New("not found"))

		result, err := bookService.GetBookById(1)

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestGetBooks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockBookRepository(ctrl)
	bookService := service.NewBookService(mockRepo)

	t.Run("Success", func(t *testing.T) {
		books := []model.Book{{ID: 1, Title: "Book 1"}, {ID: 2, Title: "Book 2"}}
		mockRepo.EXPECT().GetBooks().Return(books, nil)

		result, err := bookService.GetBooks()

		assert.NoError(t, err)
		assert.Equal(t, books, result)
	})

	t.Run("Error", func(t *testing.T) {
		mockRepo.EXPECT().GetBooks().Return(nil, errors.New("error"))

		result, err := bookService.GetBooks()

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestUpdateBook(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockBookRepository(ctrl)
	bookService := service.NewBookService(mockRepo)

	book := &model.Book{ID: 1, Title: "Updated Book"}

	t.Run("Success", func(t *testing.T) {
		mockRepo.EXPECT().UpdateBook(book).Return(nil)

		err := bookService.UpdateBook(book)

		assert.NoError(t, err)
	})

	t.Run("Error", func(t *testing.T) {
		mockRepo.EXPECT().UpdateBook(book).Return(errors.New("update error"))

		err := bookService.UpdateBook(book)

		assert.Error(t, err)
	})
}
