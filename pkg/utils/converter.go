package utils

import (
	"bookstore/internal/model"
	"database/sql"
	"log"
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
		var orderID, detailID, quantity int
		var bookID int64
		var total, subtotal, price int64
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
			log.Printf("[ConvertToDetailResponse] could not scan order row: %v", err)
			return nil, err
		}

		// If the order is not in the map, create a new OrderResponse entry
		if _, ok := orderMap[orderID]; !ok {
			orderMap[orderID] = &model.OrderResponse{
				ID:          (int64(orderID)),
				OrderDetail: []model.OrderDetailResponse{},
				Total:       *ConvertToDisplayPrice(&total),
			}
		}

		// Create a new Book entry
		book := model.Book{
			ID:     bookID,
			Title:  title,
			Author: author,
			Price:  *ConvertToDisplayPrice(&price),
		}

		// Create a new OrderDetailResponse entry
		orderDetail := model.OrderDetailResponse{
			ID:       (int64(detailID)),
			Book:     []model.Book{book},
			Quantity: (int64(quantity)),
			Subtotal: *ConvertToDisplayPrice(&subtotal),
		}

		orderMap[orderID].OrderDetail = append(orderMap[orderID].OrderDetail, orderDetail)
	}

	for _, order := range orderMap {
		orders = append(orders, *order)
	}

	if err := rows.Err(); err != nil {
		log.Printf("[ConvertToDetailResponse] could not retrieve orders: %v", err)
		return nil, err
	}

	return orders, nil
}
