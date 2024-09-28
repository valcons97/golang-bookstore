package utils

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func CheckPassword(providedPassword, storedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(providedPassword))
	if err != nil {
		log.Printf("Password mismatch: %v", err) // Log the error for debugging
		return false
	}
	return true
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
