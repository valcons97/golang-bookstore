package handler

import (
	"bookstore/internal/customer/model"
	"bookstore/internal/customer/service"
	"bookstore/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CustomerHandler struct {
	Service service.CustomerService
}

func NewCustomerHandler(service service.CustomerService) *CustomerHandler {
	return &CustomerHandler{Service: service}
}

func (h *CustomerHandler) Login(c *gin.Context) {
	var request model.Customer

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	hashedPassword, err := utils.HashPassword(request.Password)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	customer, err := h.Service.Login(request.Email, hashedPassword)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err})
		return
	}

	token, err := utils.GenerateToken(
		customer.ID,
		customer.Email,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
		"customer": gin.H{
			"id":      customer.ID,
			"email":   customer.Email,
			"name":    customer.Name,
			"address": customer.Address,
		},
	})
}

func (h *CustomerHandler) Register(c *gin.Context) {
	var request model.Customer

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
		return
	}

	hashedPassword, err := utils.HashPassword(request.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	request.Password = hashedPassword

	if err := h.Service.Register(&request); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register customer"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Customer registered successfully",
	})
}
