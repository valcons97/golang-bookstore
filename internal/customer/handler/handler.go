package handler

import (
	"bookstore/internal/customer/model"
	"bookstore/internal/customer/service"
	"bookstore/pkg/utils"
	"log"
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
	var request model.LoginRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
		return
	}

	// Check if email or password is empty
	if request.Email == "" || request.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email or password cannot be empty"})
		return
	}

	customer, err := h.Service.Login(request.Email, request.Password)
	if err != nil {
		log.Printf("[Login]: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err})
		return
	}

	token, err := utils.GenerateToken(
		customer.ID,
		customer.Email,
	)
	if err != nil {
		log.Printf("[Login]: %v", err)
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
