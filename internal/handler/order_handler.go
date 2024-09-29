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

func (h *OrderHandler) PayOrder(c *gin.Context) {
	// Get the customer ID from the JWT token
	customerID, err := utils.ExtractCustomerID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	err = h.service.PayOrder(customerID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Order paid",
	})
}

func (h *OrderHandler) GetCart(c *gin.Context) {
	// Get the customer ID from the JWT token
	customerID, err := utils.ExtractCustomerID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	orderId, err := h.service.CreateOrderIfNotExists(customerID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	response, err := h.service.GetCart(orderId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, response)
}

func (h *OrderHandler) RemoveFromCart(c *gin.Context) {
	// Get the customer ID from the JWT token
	var request model.RemoveItemFromCartRequest
	customerID, err := utils.ExtractCustomerID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	orderId, err := h.service.CreateOrderIfNotExists(customerID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = h.service.RemoveFromCart(orderId, int(request.BookId))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Book removed from cart",
	})
}

func (h *OrderHandler) AddToCart(c *gin.Context) {
	var request model.AddToCartRequest

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

	orderId, err := h.service.CreateOrderIfNotExists(customerID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	subTotal := request.Price * float64(request.Quantity)

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
		"message": "Cart updated",
	})
}
