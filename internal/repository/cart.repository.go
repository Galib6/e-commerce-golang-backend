package repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/config"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/models"
)

func GetCartByUserId(id uuid.UUID) (*models.Cart, error) {
	var cart models.Cart
	err := config.DB.Preload("CartItems").Where("user_id = ?", id).First(&cart).Error
	return &cart, err
}

func CreateCart(userId uuid.UUID) error {
	var cart models.Cart
	cart.UserID = userId
	err := config.DB.Create(&cart).Error
	return err
}

func AddOrUpdateCartItem(userId uuid.UUID, productId uuid.UUID, qty int) {

}

func CreateCartItem(cartItem *models.CartItems) error {
	return config.DB.Create(cartItem).Error
}

func UpdateCartItem(cartItem *models.CartItems) error {
	return config.DB.
		Model(&models.CartItems{}).
		Where("id = ?", cartItem.ID).
		Updates(map[string]interface{}{
			"quantity":   cartItem.Quantity,
			"updated_at": time.Now(),
		}).Error
}

func GetCartItem(cartId uuid.UUID, productId uuid.UUID) (*models.CartItems, error) {
	var cartItem models.CartItems

	err := config.DB.
		Where("cart_id = ? AND product_id = ?", cartId, productId).
		First(&cartItem).
		Error

	return &cartItem, err
}

func GetCartItems(cartId uuid.UUID) ([]models.CartItems, error) {
	var items []models.CartItems
	err := config.DB.
		Where("cart_id = ?", cartId).
		Preload("Product").
		Find(&items).
		Error
	return items, err
}

func RemoveCartItemFrom(cartId uuid.UUID, productId uuid.UUID) error {
	return config.DB.
		Unscoped().
		Where("cart_id = ? AND product_id = ?", cartId, productId).
		Delete(&models.CartItems{}).
		Error
}
