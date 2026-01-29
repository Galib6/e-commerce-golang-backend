package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/config"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/models"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/repository"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/utils"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type AddCartItemRequest struct {
	ProductID uuid.UUID `json:"product_id" binding:"required"`
	Quantity  int       `json:"quantity" binding:"required,min=1"`
}

// CreateCart godoc
// @Summary     Create a cart
// @Description Create a new cart for authenticated user
// @Tags        Cart
// @Accept      json
// @Produce     json
// @Security    ApiKeyAuth
// @Success     200  {object}  map[string]interface{}
// @Success     201  {object}  map[string]interface{}
// @Failure     401  {object}  map[string]interface{}
// @Failure     500  {object}  map[string]interface{}
// @Router      /cart [post]
func CreateCart(c *gin.Context) {
	val, ok := c.Get("userId")
	if !ok {
		utils.ResponseError(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	userId, err := uuid.Parse(val.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"invalid userid": err})
		return
	}
	cart, err := repository.GetCartByUserId(userId)

	// cart already exists
	if err == nil && cart != nil {
		utils.ResponseSuccess(c, http.StatusOK, "Cart already exists", cart)
		return
	}

	// real DB error (not record not found)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		utils.ResponseError(c, http.StatusInternalServerError, "Failed to fetch cart", nil)
		return
	}

	// create cart
	err = repository.CreateCart(userId)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, "Failed to create cart", nil)
		return
	}
	utils.ResponseSuccess(c, http.StatusCreated, "Cart created", nil)

}

// AddOrUpdateCartItem godoc
// @Summary     Add or update cart item
// @Description Add a product to cart or update quantity
// @Tags        Cart
// @Accept      json
// @Produce     json
// @Security    ApiKeyAuth
// @Param       item  body      AddCartItemRequest  true  "Cart item data"
// @Success     200   {object}  map[string]interface{}
// @Failure     400   {object}  map[string]interface{}
// @Failure     401   {object}  map[string]interface{}
// @Failure     500   {object}  map[string]interface{}
// @Router      /cart/item [post]
func AddOrUpdateCartItem(c *gin.Context) {
	var ctx = context.Background()
	var req AddCartItemRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	val, ok := c.Get("userId")
	if !ok {
		utils.ResponseError(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	userId, err := uuid.Parse(val.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	cart, err := repository.GetCartByUserId(userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = repository.CreateCart(userId)
			if err != nil {
				utils.ResponseError(c, http.StatusInternalServerError, "Failed to create cart", err)
				return
			}
			// fetch newly created cart
			cart, err = repository.GetCartByUserId(userId)
			if err != nil {
				utils.ResponseError(c, http.StatusInternalServerError, "Failed to fetch cart", err)
				return
			}
		} else {
			// real internal error
			utils.ResponseError(c, http.StatusInternalServerError, "Internal server error", nil)
			return
		}

	}
	product, err := repository.GetProductByUUID(req.ProductID)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Product does not exist", nil)
		return
	}
	cartItem, err := repository.GetCartItem(cart.ID, req.ProductID)
	currentQtyInCart := 0

	if err == nil {
		// Update quantity
		currentQtyInCart = cartItem.Quantity
		cartItem.Quantity += req.Quantity
		if product.NumberOfStock < (req.Quantity + currentQtyInCart) {
			utils.ResponseError(c, http.StatusBadRequest, "Product out of stock", nil)
			return
		}
		if err := repository.UpdateCartItem(cartItem); err != nil {
			utils.ResponseError(c, http.StatusInternalServerError, "Failed to update cart item", err)
			return
		}
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		// Add new cart item
		if (currentQtyInCart + req.Quantity) > product.NumberOfStock {
			utils.ResponseError(c, http.StatusBadRequest, "Product out of stock", nil)
			return
		}
		err := repository.CreateCartItem(&models.CartItems{
			CartID:    cart.ID,
			ProductID: req.ProductID,
			Quantity:  req.Quantity,
		})
		if err != nil {
			utils.ResponseError(c, http.StatusInternalServerError, "Failed to add cart item", err)
			return
		}
	} else {
		utils.ResponseError(c, http.StatusInternalServerError, "Failed to fetch cart item", err)
		return
	}
	cacheKey := fmt.Sprintf("cart:%s", cart.ID.String())

	// It's okay if this fails (Redis might be down), so just log it.
	if err := config.RDB.Del(ctx, cacheKey).Err(); err != nil {
		fmt.Println("Failed to clear cache:", err)
	}
}

// GetAllCartItems godoc
// @Summary     Get all cart items
// @Description Retrieve all items in user's cart (uses cache)
// @Tags        Cart
// @Accept      json
// @Produce     json
// @Security    ApiKeyAuth
// @Success     200  {object}  map[string]interface{}
// @Failure     401  {object}  map[string]interface{}
// @Failure     500  {object}  map[string]interface{}
// @Router      /cart/items [get]
func GetAllCartItems(c *gin.Context) {
	var ctx = context.Background()
	val, ok := c.Get("userId")
	if !ok {
		utils.ResponseError(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	userId, err := uuid.Parse(val.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	cart, err := repository.GetCartByUserId(userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.ResponseError(c, http.StatusInternalServerError, "Cart does not exist", err)
			return
		} else {
			utils.ResponseError(c, http.StatusInternalServerError, "Failed to fetch cart", err)
		}
	}
	cacheKey := fmt.Sprintf("cart:%s", cart.ID.String())
	cachedData, err := config.RDB.Get(ctx, cacheKey).Result()

	if err == redis.Nil {
		cartItems, err := repository.GetCartItems(cart.ID)
		if err != nil {
			utils.ResponseError(c, http.StatusUnauthorized, "Failed to fetch", err)
			return
		}
		// Serialize data to JSON to store in Redis
		jsonData, _ := json.Marshal(cartItems)

		err = config.RDB.Set(ctx, cacheKey, jsonData, 24*time.Hour).Err() //24 hour
		if err != nil {
			fmt.Println("Failed to set cache:", err) // Log error but don't fail request
		}
		utils.ResponseSuccess(c, http.StatusOK, "data fetch successfully", cartItems)
		return
	} else if err != nil {
		fmt.Println("Redis Error:", err)
		cartItems, err := repository.GetCartItems(cart.ID) // Fallback
		if err != nil {
			utils.ResponseError(c, http.StatusUnauthorized, "Failed to fetch", err)
			return
		}
		utils.ResponseSuccess(c, http.StatusOK, "data fetched from DB (Redis down)", cartItems)
		return
	}

	var cartItems []models.CartItems
	if err := json.Unmarshal([]byte(cachedData), &cartItems); err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, "Failed to parse cache", err)
		return
	}
	utils.ResponseSuccess(c, http.StatusOK, "data fetched from Redis", cartItems)
}

// GetCart godoc
// @Summary     Get cart by user ID (Admin)
// @Description Retrieve a user's cart (admin only)
// @Tags        Cart
// @Accept      json
// @Produce     json
// @Security    ApiKeyAuth
// @Param       userId  path      string  true  "User UUID"
// @Success     200     {object}  map[string]interface{}
// @Failure     401     {object}  map[string]interface{}
// @Failure     403     {object}  map[string]interface{}
// @Failure     500     {object}  map[string]interface{}
// @Router      /cart/{userId} [get]
func GetCart(c *gin.Context) {

}

// DeleteCart godoc
// @Summary     Delete cart by user ID (Admin)
// @Description Delete a user's cart (admin only)
// @Tags        Cart
// @Accept      json
// @Produce     json
// @Security    ApiKeyAuth
// @Param       userId  path      string  true  "User UUID"
// @Success     200     {object}  map[string]interface{}
// @Failure     401     {object}  map[string]interface{}
// @Failure     403     {object}  map[string]interface{}
// @Failure     500     {object}  map[string]interface{}
// @Router      /cart/{userId} [delete]
func DeleteCart(c *gin.Context) {

}

// RemoveCartItemFromCart godoc
// @Summary     Remove item from cart
// @Description Remove a specific product from cart
// @Tags        Cart
// @Accept      json
// @Produce     json
// @Security    ApiKeyAuth
// @Param       product_id  path      string  true  "Product UUID"
// @Success     200         {object}  map[string]interface{}
// @Failure     400         {object}  map[string]interface{}
// @Failure     401         {object}  map[string]interface{}
// @Failure     500         {object}  map[string]interface{}
// @Router      /cart/item/{product_id} [delete]
func RemoveCartItemFromCart(c *gin.Context) {
	ctx := context.Background()
	val, ok := c.Get("userId")
	if !ok {
		utils.ResponseError(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	userId, err := uuid.Parse(val.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	cart, err := repository.GetCartByUserId(userId)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, "Failed to fetch cart", err)
		return
	}
	productIDStr := c.Param("product_id")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid product ID", err)
		return
	}
	err = repository.RemoveCartItemFrom(cart.ID, productID)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, "Failed to remove cart item", err)
		return
	}

	cacheKey := fmt.Sprintf("cart:%s", userId.String())

	// We ignore errors here because if the key doesn't exist, it's fine.
	config.RDB.Del(ctx, cacheKey)
}
