package migration

import (
	"bookstore/internal/book"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&book.Book{})
}

// SeedBooks when first building docker image
func SeedBooks(db *gorm.DB, count int) {
	var books []book.Book
	if result := db.Find(&books); result.Error == nil && len(books) == 0 {

		seedBooks := []book.Book{
			{Title: "1984", Author: "George Orwell", Price: 9.99},
			{Title: "To Kill a Mockingbird", Author: "Harper Lee", Price: 7.99},
			{Title: "The Great Gatsby", Author: "F. Scott Fitzgerald", Price: 10.99},
			{Title: "Moby Dick", Author: "Herman Melville", Price: 11.99},
			{Title: "War and Peace", Author: "Leo Tolstoy", Price: 12.99},
		}

		db.Create(&seedBooks)
	}
}
