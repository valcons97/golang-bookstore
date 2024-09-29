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

func NewOrderHandler(service service.OrderService) *OrderHandler {
	return &OrderHandler{service: service}
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

	subTotal := request.Price * float64(request.Quantity)

	orderId, err := h.service.CreateOrderIfNotExists(customerID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Ensure that orderDetail.OrderID is set to the appropriate order ID
	if err := h.service.AddToCart(
		orderId,
		int(request.BookId),
		int(request.Quantity),
		*utils.ConvertStorePrice(&subTotal)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Book added to cart",
	})
}
