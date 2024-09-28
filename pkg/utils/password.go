package utils

import "golang.org/x/crypto/bcrypt"

func CheckPassword(providedPassword, storedPassword string) bool {
	// Implement your password hashing and comparison logic here
	// For example, if you're using bcrypt:
	// err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(providedPassword))
	// return err == nil

	return providedPassword == storedPassword // Temporary implementation for example
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
