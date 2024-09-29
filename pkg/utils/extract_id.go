package utils

import (
	"fmt"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func ExtractCustomerID(c *gin.Context) (int, error) {
	// Get the token from the Authorization header
	tokenString := c.Request.Header.Get("Authorization")
	if tokenString == "" {
		return 0, fmt.Errorf("no token provided")
	}

	// Split the token string to get the actual token
	tokenString = strings.Split(tokenString, " ")[1]

	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid token")
		}
		return []byte(os.Getenv("SECRET_KEY")), nil
	})

	if err != nil || !token.Valid {
		return 0, fmt.Errorf("invalid token")
	}

	// Extract the claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		customerID := int(claims["id"].(float64))
		return customerID, nil
	}

	return 0, fmt.Errorf("invalid token claims")
}
