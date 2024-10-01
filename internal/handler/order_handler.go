package handler

import (
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
	customerID, err := utils.ExtractCustomerID(c)
	if err != nil {
		ErrorHandler(c, http.StatusUnauthorized, "")
		return
	}

	err = h.service.PayOrder(customerID)

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
	// Get the customer ID from the JWT token
	customerID, err := utils.ExtractCustomerID(c)
	if err != nil {
		ErrorHandler(c, http.StatusUnauthorized, "")
		return
	}

	orderId, err := h.service.CreateOrderIfNotExists(customerID)

	if err != nil {
		ErrorHandler(
			c,
			http.StatusInternalServerError,
			"Unable to access cart. Please try again later.",
		)
		return
	}

	response, err := h.service.GetCart(orderId)

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
	customerID, err := utils.ExtractCustomerID(c)
	var request HistoryRequest

	if err != nil {
		ErrorHandler(c, http.StatusUnauthorized, "")
		return
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		ErrorHandler(c, http.StatusBadRequest, "")
		return
	}

	if request.Limit == 0 {
		request.Limit = 10
	}

	response, err := h.service.GetOrderHistory(customerID, request.Limit, request.Page)

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
	var request RemoveItemFromCartRequest
	customerID, err := utils.ExtractCustomerID(c)
	if err != nil {
		ErrorHandler(c, http.StatusUnauthorized, "")
		return
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		ErrorHandler(c, http.StatusBadRequest, "")
		return
	}

	orderId, err := h.service.CreateOrderIfNotExists(customerID)

	if err != nil {
		ErrorHandler(
			c,
			http.StatusInternalServerError,
			"Unable to access cart. Please try again later.",
		)
		return
	}

	err = h.service.RemoveFromCart(orderId, int(request.BookId))

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
	var request AddToCartRequest

	customerID, err := utils.ExtractCustomerID(c)
	if err != nil {
		ErrorHandler(c, http.StatusUnauthorized, "")
		return
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		ErrorHandler(c, http.StatusBadRequest, "")
		return
	}

	orderId, err := h.service.CreateOrderIfNotExists(customerID)

	if err != nil {
		ErrorHandler(
			c,
			http.StatusInternalServerError,
			"Unable to access cart. Please try again later.",
		)
		return
	}

	subTotal := request.Price * float64(request.Quantity)

	if err := h.service.AddToCart(
		orderId,
		int(request.BookId),
		int(request.Quantity),
		*utils.ConvertStorePrice(&subTotal)); err != nil {

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
