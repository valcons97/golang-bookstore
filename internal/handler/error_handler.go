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

	if statusCode == http.StatusUnauthorized && errMsg == "" {
		// Unauthorized error handling
		errMsg = "Unauthorized access. Please log in."
	} else if statusCode == http.StatusBadRequest && errMsg == "" {
		// Error related to invalid data
		errMsg = "Invalid input data"
	}

	c.JSON(statusCode, ErrorResponse{
		Error:   http.StatusText(statusCode),
		Message: errMsg,
	})
	c.Abort()
}
