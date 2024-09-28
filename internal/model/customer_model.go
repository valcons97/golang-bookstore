package model

type Customer struct {
	ID       int    `json:"id"`
	Email    string `json:"email"    binding:"required,email"` // Email field with validation
	Password string `json:"password" binding:"required"`       // Password field with validation
	Name     string `json:"name"     binding:"required"`       // Name field with validation
	Address  string `json:"address"  binding:"required"`       // Address field with validation
}
