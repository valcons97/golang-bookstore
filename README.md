# Golang Bookstore

## Bookstore Application

This repository contains a simple bookstore application built in Go, utilizing the Gin web framework and PostgreSQL for database management. The project follows a clean architecture approach, organizing code into distinct packages for clarity and maintainability.

## Folder Structure

```
golang-bookstore
├── cmd # Main entry point for the application (Gin server setup, routes initialization).
│
├── internal # Contains core application logic.
│ ├── handler # HTTP handlers that process requests and generate responses.
│ ├── middleware # Custom middleware functions (e.g., JWT authentication).
│ ├── migration # Database migrations to set up schema.
│ ├── model # Structs representing database entities (Book, Order, Customer, etc.).
│ ├── repository # Database access logic for handling CRUD operations.
│ ├── router # Route definition and grouping.
│ └── service # Business logic and service layer for handling core functionalities.
│
├── pkg # Contains utility packages.
│ └── utils # Utility functions (e.g., password hashing, token generation, price conversions).
│
├── script # Helpful scripts (e.g., DB seeding, testing utilities).
│
├── test # Unit and integration tests.
│
├── .env # Environment variables for configuration.
├── .gitignore # Specifies files to be ignored by Git.
├── docker-compose.yml # Configuration for setting up Dockerized services.
├── Dockerfile # Docker setup for the application.
├── go.mod # Dependency and module management file.
├── go.sum # Version checksum of the dependencies.
└── README.md # Project overview and setup instructions.
```

## Database Diagram

<img width="640" alt="image" src="https://github.com/user-attachments/assets/4002361a-587a-4331-bb4a-197805d1e98f">


## Project Setup

### Generating Mocks

You can generate mocks for your services using the following commands:

```bash
mockgen -source=./pkg/book/service.go -destination=./pkg/book/mock_service.go -package=book
mockgen -source=internal/service/book_service.go -destination=test/mocks/mock_book_service.go -package=mocks
mockgen -source=internal/service/customer_service.go -destination=test/mocks/mock_customer_service.go -package=mocks
mockgen -source=internal/service/order_service.go -destination=test/mocks/mock_order_service.go -package=mocks
```

To run all tests in the project, use the following command:

```
go test ./...
```

To start the application along with its dependencies, run:

```
docker-compose up --build
```

You can copy the following code into main.go to expose endpoints for checking data:

```
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
