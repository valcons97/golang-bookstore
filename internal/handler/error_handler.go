package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

func ErrorHandler(c *gin.Context, statusCode int, errMsg string) {
	// Unauthorized error handling
	if statusCode == http.StatusUnauthorized && errMsg == "" {
		errMsg = "Unauthorized access. Please log in."
	}

	// Error related to invalid data
	if statusCode == http.StatusBadRequest && errMsg == "" {
		errMsg = "Invalid input data"
	}

	c.JSON(statusCode, ErrorResponse{
		Error:   http.StatusText(statusCode),
		Message: errMsg,
	})
	c.Abort()
}
