# golang-bookstore

# Bookstore Application

This repository contains a simple bookstore application built in Go, utilizing the Gin web framework and PostgreSQL for database management. The project follows a clean architecture approach, organizing code into distinct packages for clarity and maintainability.

## Folder Structure

## Overview

-   **cmd**: Contains the main entry point for the application, setting up the Gin server and initializing routes.
-   **internal**: Holds the core components of the application:

    -   **handler**: Includes HTTP handlers that process requests and send responses.
        -   **middleware**: Contains custom middleware functions, such as JWT authentication.
        -   **migration**: Handles database migrations for schema setup.
        -   **model**: Defines structs representing database entities, including Book, Order, and Customer.
        -   **repository**: Implements the database access logic for CRUD operations.
        -   **router**: Manages route definitions and groupings for the application.
        -   **service**: Contains the business logic and service layer that handles core functionalities.

-   **pkg**: Provides utility functions that support various operations, such as password hashing and token generation.

-   **script**: Contains helpful scripts, like database seeding and testing utilities.

-   **test**: Includes unit and integration tests to ensure the application works correctly.

-   **.env**: Stores environment variables for configuration.

-   **.gitignore**: Lists files and directories that should be ignored by Git.

-   **docker-compose.yml**: Configuration file for setting up Dockerized services.

-   **Dockerfile**: Defines the Docker setup for the application.

-   **go.mod**: Manages dependencies and module information for the project.

-   **go.sum**: Contains checksums for the project dependencies to ensure integrity.

-   **README.md**: Provides an overview of the project and instructions for setup and usage.

## Getting Started

To get started with this project, please refer to the setup instructions below.

`````bash
# Example commands for generating mocks and running tests
mockgen -source=./pkg/book/service.go -destination=./pkg/book/mock_service.go -package=book
mockgen -source=internal/service/book_service.go -destination=test/mocks/mock_book_service.go -package=mocks
mockgen -source=internal/service/customer_service.go -destination=test/mocks/mock_customer_service.go -package=mocks
mockgen -source=internal/service/order_service.go -destination=test/mocks/mock_order_service.go -package=mocks

go test ./...

## DB Diagram

<img width="672" alt="image" src="https://github.com/user-attachments/assets/297ffc1e-72c3-4156-9fe1-9b33ceb5dae0">

Here’s the content formatted as a single, cohesive `README.md`:

````markdown
# Project Setup

## Generating Mocks

You can generate mocks for your services using the following commands:

```bash
mockgen -source=./pkg/book/service.go -destination=./pkg/book/mock_service.go -package=book
mockgen -source=internal/service/book_service.go -destination=test/mocks/mock_book_service.go -package=mocks
mockgen -source=internal/service/customer_service.go -destination=test/mocks/mock_customer_service.go -package=mocks
mockgen -source=internal/service/order_service.go -destination=test/mocks/mock_order_service.go -package=mocks
`````

````

## Running Tests

To run all tests in the project, use the following command:

```bash
go test ./...
```

## Example Routes

You can copy the following code into `main.go` to expose endpoints for checking data:

```go
r.GET("/tables", func(c *gin.Context) {
	// Get the table names from the database schema using GORM's Migrator
	tables := []string{}

	migrator := db.Migrator()

	tables, err := migrator.GetTables()
	if err != nil {
		return
	}

	// Return the table names as JSON response
	c.JSON(http.StatusOK, tables)
})

r.GET("/customers", func(c *gin.Context) {
	// Declare a slice to hold customers
	var customers []model.Customer

	// Use GORM to find all customers
	err := db.Find(&customers).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the customers as JSON response
	c.JSON(http.StatusOK, customers)
})

r.GET("/orders", func(c *gin.Context) {
	// Declare a slice to hold orders
	var orders []model.Order

	// Use GORM to find all orders
	err := db.Find(&orders).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the orders as JSON response
	c.JSON(http.StatusOK, orders)
})

r.GET("/orderDetails", func(c *gin.Context) {
	// Declare a slice to hold order details
	var orderDetails []model.OrderDetail

	// Use GORM to find all order details
	err := db.Find(&orderDetails).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the order details as JSON response
	c.JSON(http.StatusOK, orderDetails)
})
```

## Starting the Application with Docker Compose

To start the application along with its dependencies, run:

```bash
docker-compose up
```
````
