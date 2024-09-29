package utils

import (
	"bookstore/internal/model"
	"database/sql"
	"fmt"
)

func ConvertStorePrice(value *float64) *int64 {

	if value == nil {
		return nil
	}

	result := int64(*value * 100)

	return &result
}

func ConvertToDisplayPrice(value *int64) *float64 {
	if value == nil {
		return nil
	}

	result := float64(*value) / 100

	return &result
}

func ConvertToDetailResponse(rows *sql.Rows) ([]model.OrderResponse, error) {
	defer rows.Close()

	var orders []model.OrderResponse
	orderMap := make(map[int]*model.OrderResponse)

	for rows.Next() {
		var orderID, detailID, bookID, quantity int
		var total, subtotal, price float64
		var title, author string

		err := rows.Scan(
			&orderID,
			&total,
			&detailID,
			&bookID,
			&quantity,
			&subtotal,
			&title,
			&author,
			&price,
		)
		if err != nil {
			return nil, fmt.Errorf("could not scan order row: %w", err)
		}

		// If the order is not in the map, create a new OrderResponse entry
		if _, ok := orderMap[orderID]; !ok {
			orderMap[orderID] = &model.OrderResponse{
				ID:          (int64(orderID)),
				OrderDetail: []model.OrderDetailResponse{},
				Total:       total,
			}
		}

		// Create a new Book entry
		book := model.Book{
			ID:     bookID,
			Title:  title,
			Author: author,
			Price:  price,
		}

		// Create a new OrderDetailResponse entry
		orderDetail := model.OrderDetailResponse{
			ID:       (int64(detailID)),
			Book:     []model.Book{book}, // Assuming multiple books could be part of order details
			Quantity: (int64(quantity)),
			Subtotal: subtotal,
		}

		// Append the order details to the corresponding order
		orderMap[orderID].OrderDetail = append(orderMap[orderID].OrderDetail, orderDetail)
	}

	// Convert the map to a slice
	for _, order := range orderMap {
		orders = append(orders, *order)
	}

	// Check for errors during row iteration
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("could not retrieve orders: %w", err)
	}

	return orders, nil
}
