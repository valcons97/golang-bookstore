package order

import "time"

type Order struct {
	ID         int
	CustomerID int
	BookID     int
	CreatedAt  time.Time
}
