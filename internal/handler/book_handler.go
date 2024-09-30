package handler

import (
	"bookstore/internal/model"
	"bookstore/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BookHandler struct {
	Service service.BookService
}

func NewBookHandler(service service.BookService) *BookHandler {
	return &BookHandler{Service: service}
}

func (h *BookHandler) GetBooks(c *gin.Context) {
	books, err := h.Service.GetBooks()
	if err != nil {
		ErrorHandler(c, http.StatusInternalServerError, "Failed to retrieve books")
		return
	}

	c.JSON(http.StatusOK, books)
}

func (h *BookHandler) GetBookById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		ErrorHandler(c, http.StatusBadRequest, "")
		return
	}

	book, err := h.Service.GetBookById(id)
	if err != nil {
		ErrorHandler(c, http.StatusNotFound, "Book not found")
		return
	}

	c.JSON(http.StatusOK, book)
}

func (h *BookHandler) CreateBook(c *gin.Context) {
	var book model.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		ErrorHandler(c, http.StatusBadRequest, "")
		return
	}
	id, err := h.Service.CreateBook(&book)
	if err != nil {
		ErrorHandler(c, http.StatusInternalServerError, "Failed to create book")
		return
	}

	book.ID = id

	c.JSON(http.StatusCreated, book)
}

func (h *BookHandler) UpdateBook(c *gin.Context) {
	var book model.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		ErrorHandler(c, http.StatusBadRequest, "")
		return
	}

	if err := h.Service.UpdateBook(&book); err != nil {
		ErrorHandler(c, http.StatusNotFound, "Book not found")
		return
	}

	c.JSON(http.StatusOK, book)
}
