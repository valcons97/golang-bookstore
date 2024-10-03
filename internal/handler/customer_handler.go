package handler

import (
	"bookstore/internal/handler/request"
	"bookstore/internal/model"
	"bookstore/internal/service"
	"bookstore/pkg/utils"
	"errors"

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

	token, err := h.Service.Login(request.Email, request.Password)
	if err != nil {
		if errors.Is(err, utils.ErrEmptyEmailOrPassword) {
			ErrorHandler(c, http.StatusBadRequest, "Email or password cannot be empty")
		} else if errors.Is(err, utils.ErrEmailNotFound) {
			ErrorHandler(c, http.StatusUnauthorized, "Email not registered")
		} else if errors.Is(err, utils.ErrWrongPassword) {
			ErrorHandler(c, http.StatusUnauthorized, "Wrong Password")
		} else {
			ErrorHandler(c, http.StatusInternalServerError, err.Error())
		}

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

	err := h.Service.Register(&request)

	if err != nil {
		if errors.Is(err, utils.ErrDuplicateEmail) {
			ErrorHandler(c, http.StatusBadRequest, "Email already registered")
		} else {
			ErrorHandler(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Customer registered successfully",
	})
}
