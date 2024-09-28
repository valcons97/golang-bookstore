package handler_test

// import (
// 	"bookstore/pkg/book"
// 	"bytes"
// 	"encoding/json"
// 	"errors"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/gin-gonic/gin"
// 	gomock "github.com/golang/mock/gomock"
// 	"github.com/stretchr/testify/assert"
// )

// func TestCreateBook(t *testing.T) {
// 	// Create a controller for GoMock
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	// Create a mock service using the generated mock
// 	mockService := book.NewMockService(ctrl)

// 	// Create the handler that will use the mock service
// 	bookHandler := handler.NewBookHandler(mockService)

// 	// Create a new Gin router and register the handler
// 	r := gin.Default()
// 	r.POST("/books", bookHandler.CreateBook)

// 	// Create a test case with a valid request
// 	t.Run("success", func(t *testing.T) {
// 		// Prepare the DTO and marshal it into JSON
// 		bookDTO := book.CreateBookDTO{
// 			Title:  "1984",
// 			Author: "George Orwell",
// 			Price:  9.99,
// 		}
// 		requestBody, _ := json.Marshal(bookDTO)

// 		// Prepare the expected book to be passed to the service
// 		newBook := book.Book{
// 			Title:  bookDTO.Title,
// 			Author: bookDTO.Author,
// 			Price:  bookDTO.Price,
// 		}

// 		// Expect that the mock service's CreateBook method will be called with the correct book
// 		mockService.EXPECT().CreateBook(&newBook).Return(nil)

// 		// Create a test HTTP request
// 		req, _ := http.NewRequest(http.MethodPost, "/books", bytes.NewBuffer(requestBody))
// 		req.Header.Set("Content-Type", "application/json")

// 		// Create a response recorder to capture the response
// 		rr := httptest.NewRecorder()

// 		// Perform the request
// 		r.ServeHTTP(rr, req)

// 		// Check the status code and response body
// 		assert.Equal(t, http.StatusCreated, rr.Code)
// 		assert.Contains(t, rr.Body.String(), "Book created successfully")
// 	})

// 	// Create a test case for validation errors
// 	t.Run("validation error", func(t *testing.T) {
// 		// Send an invalid request (missing price)
// 		invalidBookDTO := `{"title": "1984", "author": "George Orwell"}`
// 		req, _ := http.NewRequest(
// 			http.MethodPost,
// 			"/books",
// 			bytes.NewBuffer([]byte(invalidBookDTO)),
// 		)
// 		req.Header.Set("Content-Type", "application/json")

// 		// Create a response recorder
// 		rr := httptest.NewRecorder()

// 		// Perform the request
// 		r.ServeHTTP(rr, req)

// 		// Check that the status code is 400 and contains a validation error
// 		assert.Equal(t, http.StatusBadRequest, rr.Code)
// 		assert.Contains(t, rr.Body.String(), "Field validation for 'Price' failed")
// 	})

// 	// Create a test case for service errors
// 	t.Run("service error", func(t *testing.T) {
// 		// Prepare the DTO and marshal it into JSON
// 		bookDTO := book.CreateBookDTO{
// 			Title:  "1984",
// 			Author: "George Orwell",
// 			Price:  9.99,
// 		}
// 		requestBody, _ := json.Marshal(bookDTO)

// 		// Prepare the expected book to be passed to the service
// 		newBook := book.Book{
// 			Title:  bookDTO.Title,
// 			Author: bookDTO.Author,
// 			Price:  bookDTO.Price,
// 		}

// 		// Mock the service to return an error
// 		mockService.EXPECT().CreateBook(&newBook).Return(errors.New("failed to create book"))

// 		// Create a test HTTP request
// 		req, _ := http.NewRequest(http.MethodPost, "/books", bytes.NewBuffer(requestBody))
// 		req.Header.Set("Content-Type", "application/json")

// 		// Create a response recorder to capture the response
// 		rr := httptest.NewRecorder()

// 		// Perform the request
// 		r.ServeHTTP(rr, req)

// 		// Check the status code and response body
// 		assert.Equal(t, http.StatusInternalServerError, rr.Code)
// 		assert.Contains(t, rr.Body.String(), "Failed to create book")
// 	})
// }
