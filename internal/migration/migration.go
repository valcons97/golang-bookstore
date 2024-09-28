package migration

import (
	"bookstore/internal/model"
	converter "bookstore/pkg/utils"
	"database/sql"
	"fmt"
	"log"
)

func Migrate(db *sql.DB) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS customers (
            id SERIAL PRIMARY KEY,
            email VARCHAR(255) UNIQUE NOT NULL,
            password VARCHAR(255) NOT NULL,
            name VARCHAR(255) NOT NULL,
            address VARCHAR(255) NOT NULL
        )`,
		`CREATE TABLE IF NOT EXISTS books (
            id SERIAL PRIMARY KEY,
            title VARCHAR(255) NOT NULL,
            author VARCHAR(255) NOT NULL,
            price bigint
        )`,
		`CREATE TABLE IF NOT EXISTS orders (
            id SERIAL PRIMARY KEY,
            customer_id INT NOT NULL,
            updated_at TIMESTAMP DEFAULT NOW(),
            order_state INT DEFAULT 1,
            total bigint,
            CONSTRAINT fk_customer
                FOREIGN KEY(customer_id) 
                REFERENCES customers(id)
        )`,
		`CREATE TABLE IF NOT EXISTS order_details (
            id SERIAL PRIMARY KEY,
            order_id INT NOT NULL,
            book_id INT NOT NULL,
            quantity INT NOT NULL,
            subtotal bigint NOT NULL,
            CONSTRAINT fk_order
                FOREIGN KEY(order_id) 
                REFERENCES orders(id) ON DELETE CASCADE,
            CONSTRAINT fk_book
                FOREIGN KEY(book_id) 
                REFERENCES books(id)
        )`,
	}

	for _, query := range queries {
		_, err := db.Exec(query)
		if err != nil {
			return fmt.Errorf("failed to execute this query: %v, error: %v", query, err)
		}
	}
	return nil
}

// SeedBooks when first building docker image
func SeedBooks(db *sql.DB) {
	// Check if the books table is empty
	var count int
	query := "SELECT COUNT(*) FROM books"
	err := db.QueryRow(query).Scan(&count)
	if err != nil {
		log.Fatalf("Error checking books table: %v", err)
	}

	// If the table is empty, insert seed data
	if count == 0 {
		seedBooks := []model.Book{
			{Title: "1984", Author: "George Orwell", Price: 9.99},
			{Title: "To Kill a Mockingbird", Author: "Harper Lee", Price: 7.99},
			{Title: "The Great Gatsby", Author: "F. Scott Fitzgerald", Price: 10.99},
			{Title: "Moby Dick", Author: "Herman Melville", Price: 11.99},
			{Title: "War and Peace", Author: "Leo Tolstoy", Price: 12.99},
		}

		// Insert each book into the database
		for _, book := range seedBooks {
			_, err := db.Exec(
				"INSERT INTO books (title, author, price) VALUES ($1, $2, $3)",
				book.Title,
				book.Author,
				converter.ConvertStorePrice(&book.Price),
			)
			if err != nil {
				log.Fatalf("Error inserting book %s: %v", book.Title, err)
			}
		}
	}
}
