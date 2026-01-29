package repository

import (
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/config"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/models"
)

func Checkout() (*models.Order, error) {
	var order models.Order

	return &order, nil
}

func CreateOrder(order *models.Order) (*models.Order, error) {
	if err := config.DB.Create(order).Error; err != nil {
		return nil, err
	}
	return order, nil
}
