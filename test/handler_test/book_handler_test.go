package handler_test

import (
	"bookstore/internal/handler"
	"bookstore/internal/model"
	"bookstore/pkg/utils"
	"bookstore/test/mocks"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestBookHandler_GetBooks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookService := mocks.NewMockBookService(ctrl)
	h := handler.NewBookHandler(mockBookService)

	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.GET("/books", h.GetBooks)

	t.Run("success", func(t *testing.T) {
		mockBooks := []model.Book{
			{ID: 1, Title: "Book 1", Author: "Author 1"},
			{ID: 2, Title: "Book 2", Author: "Author 2"},
		}

		mockBookService.EXPECT().GetBooks().Return(mockBooks, nil)

		req, _ := http.NewRequest(http.MethodGet, "/books", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var books []model.Book
		err := json.Unmarshal(w.Body.Bytes(), &books)
		assert.NoError(t, err)
		assert.Equal(t, mockBooks, books)
	})

	t.Run("error", func(t *testing.T) {
		mockBookService.EXPECT().GetBooks().Return(nil, errors.New("failed to retrieve books"))

		req, _ := http.NewRequest(http.MethodGet, "/books", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestBookHandler_GetBookById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookService := mocks.NewMockBookService(ctrl)
	h := handler.NewBookHandler(mockBookService)

	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.GET("/books/:id", h.GetBookById)

	t.Run("success", func(t *testing.T) {
		mockBook := model.Book{ID: 1, Title: "Book 1", Author: "Author 1", Price: 1.19}
		mockBookService.EXPECT().GetBookById(1).Return(&mockBook, nil)

		req, _ := http.NewRequest(http.MethodGet, "/books/1", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var book model.Book
		err := json.Unmarshal(w.Body.Bytes(), &book)
		assert.NoError(t, err)
		assert.Equal(t, mockBook, book)
	})

	t.Run("invalid id", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/books/abc", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("book not found", func(t *testing.T) {
		mockBookService.EXPECT().
			GetBookById(1).
			Return((*model.Book)(nil), utils.ErrBookNotFound)

		req, _ := http.NewRequest(http.MethodGet, "/books/1", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestBookHandler_CreateBook(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookService := mocks.NewMockBookService(ctrl)
	h := handler.NewBookHandler(mockBookService)

	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.POST("/books", h.CreateBook)

	t.Run("success", func(t *testing.T) {
		newBook := model.Book{ID: int64(1), Title: "New Book", Author: "New Author", Price: 1.23}
		mockBookService.EXPECT().CreateBook(&newBook).Return(nil)

		jsonBook, _ := json.Marshal(newBook)
		req, _ := http.NewRequest(http.MethodPost, "/books", bytes.NewBuffer(jsonBook))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("bad request", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "/books", bytes.NewBuffer([]byte("{bad json}")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("service error", func(t *testing.T) {
		newBook := model.Book{Title: "New Book", Author: "New Author"}
		mockBookService.EXPECT().
			CreateBook(&newBook).
			Return(utils.ErrBookNotFound)

		jsonBook, _ := json.Marshal(newBook)
		req, _ := http.NewRequest(http.MethodPost, "/books", bytes.NewBuffer(jsonBook))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestBookHandler_UpdateBook(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookService := mocks.NewMockBookService(ctrl)
	h := handler.NewBookHandler(mockBookService)

	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.PUT("/books", h.UpdateBook)

	t.Run("success", func(t *testing.T) {
		updatedBook := model.Book{ID: 1, Title: "Updated Book", Author: "Updated Author"}
		mockBookService.EXPECT().UpdateBook(&updatedBook).Return(nil)

		jsonBook, _ := json.Marshal(updatedBook)
		req, _ := http.NewRequest(http.MethodPut, "/books", bytes.NewBuffer(jsonBook))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("bad request", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPut, "/books", bytes.NewBuffer([]byte("{bad json}")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("book not found", func(t *testing.T) {
		updatedBook := model.Book{ID: 1, Title: "Updated Book", Author: "Updated Author"}
		mockBookService.EXPECT().UpdateBook(&updatedBook).Return(utils.ErrBookNotFound)

		jsonBook, _ := json.Marshal(updatedBook)
		req, _ := http.NewRequest(http.MethodPut, "/books", bytes.NewBuffer(jsonBook))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}
