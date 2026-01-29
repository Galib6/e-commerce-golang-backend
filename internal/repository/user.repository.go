package repository

import (
	"github.com/google/uuid"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/config"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/helper"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/models"
)

func Login() {}

func Register(user *models.User) (*models.User, error) {
	if err := config.DB.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func GetUserByUUID(id uuid.UUID) (*models.User, error) {
	var user models.User
	if err := config.DB.First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := config.DB.Select("id", "username", "email", "created_at").First(&user, "email = ?", email).Error; err != nil {
		return nil, err // GORM returns gorm.ErrRecordNotFound if no match
	}
	return &user, nil
}

func GetAllUsers() ([]models.User, error) {
	var users []models.User

	result := config.DB.
		Select("id", "username", "email", "created_at").
		Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}

	return users, nil
}

func FilterAndSearchUsers(params helper.UserFilterParams) (*[]models.User, int64, error) {
	var users []models.User
	var total int64
	offset := (params.Page - 1) * params.Limit
	query := config.DB.Model(&models.User{})
	if params.FullName != "" {
		query = query.Where("fullname ILIKE ?", "%"+params.FullName+"%")
	}

	// 3. Filter by Product Name (Joined via Orders -> OrderItems -> Products)
	if params.ProductName != "" {
		// Path: Users -> Orders -> OrderItems -> Products
		query = query.Joins("JOIN orders ON orders.user_id = users.id").
			Joins("JOIN order_items ON order_items.order_id = orders.id").
			Joins("JOIN products ON order_items.product_id = products.id").
			Where("products.product_name ILIKE ?", "%"+params.ProductName+"%").
			Group("users.id")
	}

	if err := config.DB.Table("(?) as sub", query).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := query.Limit(params.Limit).
		Offset(offset).
		Preload("Orders.OrderItems.Product").
		Find(&users).Error

	return &users, total, err
}
