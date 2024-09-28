# golang-bookstore

# Bookstore Application

This repository contains a simple bookstore application built in Go, utilizing the Gin web framework and PostgreSQL for database management. The project follows a clean architecture approach, organizing code into distinct packages for clarity and maintainability.

## Folder Structure

```plaintext
golang-bookstore/
├── cmd/                    # Entry point for the application
│   └── main.go             # Main application file
├── internal/               # Internal application code
│   ├── book/               # Book-related functionality
│   │   ├── handler/        # HTTP handlers for book routes
│   │   │   └── handler.go   # Book handler implementation
│   │   ├── model/          # Data models for book
│   │   │   └── model.go     # Book model definition
│   │   ├── repository/     # Data access layer for books
│   │   │   ├── book_repository_impl.go # Implementation of the book repository
│   │   │   └── book_repository.go      # Book repository interface
│   │   └── service/        # Business logic for book operations
│   │       ├── book_service_impl.go   # Implementation of the book service
│   │       └── book_service.go         # Book service interface
│   │   └── route.go        # Routes for book-related endpoints
│   ├── customer/           # Customer-related functionality
│   │   ├── auth/           # Authentication logic for customers
│   │   │   └── auth.go     # Auth functions for customer login/registration
│   │   ├── model/          # Data models for customers
│   │   │   └── model.go     # Customer model definition
│   │   ├── repository/     # Data access layer for customers
│   │   │   ├── customer_repository_impl.go # Implementation of the customer repository
│   │   │   └── customer_repository.go      # Customer repository interface
│   │   └── service/        # Business logic for customer operations
│   │       ├── customer_service_impl.go   # Implementation of the customer service
│   │       └── customer_service.go         # Customer service interface
│   ├── migration/          # Database migration scripts
│   │   └── migration.go     # Database migration logic
│   ├── order/              # Order-related functionality
│   ├── pkg/                # Utility packages
│   │   └── utils/          # General utility functions
│   │       ├── converter.go # Converter utility functions
│   │       └── password.go  # Password hashing utility functions
│   └── script/             # Scripts for database seeding and other tasks
│   │   └── seed.go         # Seed database with initial data
│   ├── test/                # For unit testing
```

mockgen -source=./pkg/book/service.go -destination=./pkg/book/mock_service.go -package=book
