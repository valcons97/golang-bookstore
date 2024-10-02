package handler

import (
	"bookstore/internal/handler/request"
	"bookstore/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	service service.OrderService
}

func NewOrderHandler(service service.OrderService) *OrderHandler {
	return &OrderHandler{service: service}
}

func (h *OrderHandler) PayOrder(c *gin.Context) {
	id, exists := c.Get("customerID")
	if !exists {
		ErrorHandler(c, http.StatusUnauthorized, "")
		return
	}

	err := h.service.PayOrder(id.(int))

	if err != nil {
		ErrorHandler(
			c,
			http.StatusInternalServerError,
			"Failed to process payment. Please try again later.",
		)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order paid successfully"})
}

func (h *OrderHandler) GetCart(c *gin.Context) {
	id, exists := c.Get("customerID")
	if !exists {
		ErrorHandler(c, http.StatusUnauthorized, "")
		return
	}

	response, err := h.service.GetCart(id.(int))

	if err != nil {
		ErrorHandler(
			c,
			http.StatusInternalServerError,
			"Unable to retrieve cart. Please try again later.",
		)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *OrderHandler) GetOrderHistory(c *gin.Context) {
	var request request.HistoryRequest

	id, exists := c.Get("customerID")
	if !exists {
		ErrorHandler(c, http.StatusUnauthorized, "")
		return
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		ErrorHandler(c, http.StatusBadRequest, "")
		return
	}

	response, err := h.service.GetOrderHistory(id.(int), request)

	if err != nil {
		ErrorHandler(
			c,
			http.StatusInternalServerError,
			"Unable to retrieve order history. Please try again later.",
		)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *OrderHandler) RemoveFromCart(c *gin.Context) {
	// Get the customer ID from the JWT token
	var request request.RemoveItemFromCartRequest
	id, exists := c.Get("customerID")
	if !exists {
		ErrorHandler(c, http.StatusUnauthorized, "")
		return
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		ErrorHandler(c, http.StatusBadRequest, "")
		return
	}

	err := h.service.RemoveFromCart(id.(int), int(request.BookId))

	if err != nil {
		ErrorHandler(
			c,
			http.StatusInternalServerError,
			"Unable to remove book from cart. Please try again later.",
		)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Book removed from cart",
	})
}

func (h *OrderHandler) AddToCart(c *gin.Context) {
	var request request.AddToCartRequest

	id, exists := c.Get("customerID")
	if !exists {
		ErrorHandler(c, http.StatusUnauthorized, "")
		return
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		ErrorHandler(c, http.StatusBadRequest, "")
		return
	}

	if err := h.service.AddToCart(id.(int), request); err != nil {
		ErrorHandler(
			c,
			http.StatusInternalServerError,
			"Unable to add book to cart. Please try again later.",
		)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Cart updated",
	})
}
