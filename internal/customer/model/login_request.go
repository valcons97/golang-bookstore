package model

type LoginRequest struct {
	Email    string `json:"email"    binding:"required,email"` // Email field with validation
	Password string `json:"password" binding:"required"`       // Password field with validation
}
