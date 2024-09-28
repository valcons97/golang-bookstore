package handler

import (
	"bookstore/internal/model"
	"bookstore/internal/service"
	"bookstore/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	service service.OrderService
}

func (h *OrderHandler) AddBookToOrder(c *gin.Context) {
	var request model.AddOrderRequest

	// Get the customer ID from the JWT token
	customerID, err := utils.ExtractCustomerID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Bind the JSON input to request
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	orderId, err := h.service.GetOrdersByCustomerID()

	detail := model.OrderDetail{}

	// Optionally set the CustomerID in the orderDetail if needed
	// This is useful if you want to keep track of which customer added the book
	orderDetail.CustomerID = customerID

	// Here, you might want to calculate subtotal based on the book's price
	orderDetail.Subtotal = calculateSubtotal(orderDetail.BookID, orderDetail.Quantity)

	// Ensure that orderDetail.OrderID is set to the appropriate order ID
	if err := h.service.AddBookToOrder(&orderDetail); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orderDetail)
}

func calculateSubtotal(bookID int64, quantity int) int64 {
	// Placeholder logic for fetching book price and calculating subtotal
	// You should replace this with actual logic to fetch the book's price
	bookPrice := 100 // Assume a hardcoded price for now
	return int64(quantity) * int64(bookPrice)
}
