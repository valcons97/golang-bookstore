package handler

import (
	"bookstore/internal/handler/request"
	"bookstore/internal/model"
	"bookstore/internal/service"
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
	var request request.LoginRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		ErrorHandler(c, http.StatusBadRequest, "")
		return
	}

	// Check if email or password is empty
	if request.Email == "" || request.Password == "" {
		ErrorHandler(c, http.StatusBadRequest, "Email or password cannot be empty")
		return
	}

	token, err := h.Service.Login(request.Email, request.Password)
	if err != nil {
		ErrorHandler(c, http.StatusUnauthorized, "Failed to login")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
	})
}

func (h *CustomerHandler) Register(c *gin.Context) {
	var request model.Customer

	if err := c.ShouldBindJSON(&request); err != nil {
		ErrorHandler(c, http.StatusBadRequest, "")
		return
	}

	hashedPassword, err := utils.HashPassword(request.Password)
	if err != nil {
		ErrorHandler(c, http.StatusInternalServerError, "Failed")
		return
	}
	request.Password = hashedPassword

	if err := h.Service.Register(&request); err != nil {
		ErrorHandler(c, http.StatusInternalServerError, "Failed to register customer")
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Customer registered successfully",
	})
}
