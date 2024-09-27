package models

import (
	"math/rand"

	"gorm.io/gorm"
)

type Book struct {
    Id int `json:"id" gorm:"primaryKey"`
    Title string `json:"title"`
    Author string `json:"author"`
    Description string `json:"description"`
}

func Migrate(db *gorm.DB) {
    db.AutoMigrate(&Book{})
}

// SeedBooks generates random book data
func SeedBooks(db *gorm.DB, count int) {
	titles := []string{"1984", "To Kill a Mockingbird", "The Great Gatsby", "Moby Dick", "War and Peace"}
	authors := []string{"George Orwell", "Harper Lee", "F. Scott Fitzgerald", "Herman Melville", "Leo Tolstoy"}

	for i := 0; i < count; i++ {
		book := Book{
			Title:  titles[rand.Intn(len(titles))],
			Author: authors[rand.Intn(len(authors))],
		}
		db.Create(&book)
	}
}
