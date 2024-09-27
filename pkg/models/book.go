package models

import (
	"gorm.io/gorm"
)

type Book struct {
	Id     int    `json:"id"     gorm:"primaryKey"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&Book{})
}

// SeedBooks when first building docker image
func SeedBooks(db *gorm.DB, count int) {
	var books []Book
	if result := db.Find(&books); result.Error == nil && len(books) == 0 {

		seedBooks := []Book{
			{Title: "1984", Author: "George Orwell"},
			{Title: "To Kill a Mockingbird", Author: "Harper Lee"},
			{Title: "The Great Gatsby", Author: "F. Scott Fitzgerald"},
			{Title: "Moby Dick", Author: "Herman Melville"},
			{Title: "War and Peace", Author: "Leo Tolstoy"},
		}

		db.Create(&seedBooks)
	}
}
