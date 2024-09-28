package model

type Book struct {
	ID     int     `json:"id"`     // Unique identifier for the book
	Title  string  `json:"title"`  // Title of the book
	Author string  `json:"author"` // Author of the book
	Price  float64 `json:"price"`  // Price of the book
}
