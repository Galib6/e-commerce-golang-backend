package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/config"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/helper"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/models"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/repository"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/utils"
)

// use goroutines and channels to handle multiple tasks concurrently
// Scenario: "Load Product Page"
// When a user views an iPhone page, you need:
// Product Details (Name, Price).
// Reviews (4.5 Stars).
// "People also bought" recommendations.
// Estimated Delivery Date.

// GetAllProducts godoc
// @Summary     Get all products
// @Description Retrieve list of all products
// @Tags        Products
// @Accept      json
// @Produce     json
// @Success     200  {object}  map[string]interface{}
// @Failure     500  {object}  map[string]interface{}
// @Router      /products/all [get]
func GetAllProducts(c *gin.Context) {
	products, err := repository.GetAllProduct()
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, "Something went wrong", nil)
		return
	}
	utils.ResponseSuccess(c, http.StatusOK, "data fetched successfully", products)
}

// GetProductById godoc
// @Summary     Get product by ID
// @Description Retrieve a single product by UUID (uses cache)
// @Tags        Products
// @Accept      json
// @Produce     json
// @Security    ApiKeyAuth
// @Param       id   path      string  true  "Product UUID"
// @Success     200  {object}  map[string]interface{}
// @Failure     400  {object}  map[string]interface{}
// @Failure     500  {object}  map[string]interface{}
// @Router      /products/{id} [get]
func GetProductById(c *gin.Context) {
	id := c.Param("id")
	productID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	// product, err := repository.GetProductByUUID(productID)
	product, err := GetProductWithCache(productID)

	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, "Something went wrong", err)
		return
	}

	utils.ResponseSuccess(c, http.StatusOK, "data fetched successfully", product)
}

// CreateNewProduct godoc
// @Summary     Create a new product
// @Description Create a product (requires auth)
// @Tags        Products
// @Accept      json
// @Produce     json
// @Security    ApiKeyAuth
// @Param       product  body      helper.CreateProductRequest  true  "Product data"
// @Success     200      {object}  map[string]interface{}
// @Failure     400      {object}  map[string]interface{}
// @Failure     401      {object}  map[string]interface{}
// @Failure     500      {object}  map[string]interface{}
// @Router      /products/ [post]
func CreateNewProduct(c *gin.Context) {
	var req helper.CreateProductRequest
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

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if err := config.Validate.Struct(req); err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Validation failed", err)
		return
	}

	if err := helper.CustomValidate(&req); err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Validation failed", err)
		return
	}
	var modelImages []models.ProductImages
	uploadedURLs := []string{}
	for i, fileHeader := range req.ImageFiles {

		// Call our helper to upload to S3/Cloudinary
		imgURL, err := utils.UploadFile(fileHeader)
		if err != nil {
			for _, url := range uploadedURLs {
				if delErr := utils.DeleteFile(url); delErr != nil {
					log.Printf("failed to cleanup uploaded file %s: %v", url, delErr)
				}
			}
			utils.ResponseError(c, http.StatusInternalServerError, "Image upload failed", err)
			return
		}
		uploadedURLs = append(uploadedURLs, imgURL)

		modelImages = append(modelImages, models.ProductImages{
			ImageUrl:  imgURL,
			IsPrimary: (i == 0),
			SortOrder: i,
		})
	}

	product := models.Product{
		Name:             req.Name,
		ShortDescription: req.Description,
		BasePrice:        req.BasePrice,
		// for now static, later get from auth middleware
		CreatedBy:       userId, // uuid.MustParse(userId.String()),
		DiscountPercent: req.DiscountPercent,
		ProductImages:   modelImages,
	}
	createdProduct, err := repository.CreateProduct(&product)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, "Something went wrong", err)
		return
	}
	utils.ResponseSuccess(c, http.StatusOK, "product created successfully", createdProduct)
}

// UpdateProduct godoc
// @Summary     Update a product
// @Description Update product by ID (requires auth)
// @Tags        Products
// @Accept      json
// @Produce     json
// @Security    ApiKeyAuth
// @Param       id       path      string          true  "Product UUID"
// @Param       product  body      models.Product  true  "Product data"
// @Success     200      {object}  map[string]interface{}
// @Failure     400      {object}  map[string]interface{}
// @Failure     500      {object}  map[string]interface{}
// @Router      /products/{id} [put]
func UpdateProduct(c *gin.Context) {
	var ctx = context.Background()
	var product models.Product
	id := c.Param("id")
	productID, err := uuid.Parse(id)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid Id", err)
		return
	}
	if err := c.ShouldBindJSON(&product); err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	existingProduct, err := repository.GetProductByUUID(productID)

	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, "Something went wrong", err)
		return
	}

	if existingProduct.CreatedBy != uuid.MustParse("f064c2fc-3523-4e1b-b166-e2cc57064fd8") {
		utils.ResponseError(c, http.StatusBadRequest, "You are not authorized", err)
		return
	}

	if err := repository.UpdateProduct(&product); err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, "Update failed", err)
		return
	}
	// 2. Invalidate Cache (Delete the specific key)
	cacheKey := fmt.Sprintf("product:details:%s", product.ID.String())
	if err := config.RDB.Del(ctx, cacheKey).Err(); err != nil {
		fmt.Println("Failed to clear cache:", err)
	}
	utils.ResponseSuccess(c, http.StatusOK, "product updated successfully", product)
}

// DeleteProduct godoc
// @Summary     Delete a product
// @Description Delete product by ID (requires auth)
// @Tags        Products
// @Accept      json
// @Produce     json
// @Security    ApiKeyAuth
// @Param       id   path      string  true  "Product UUID"
// @Success     200  {object}  map[string]interface{}
// @Failure     400  {object}  map[string]interface{}
// @Failure     500  {object}  map[string]interface{}
// @Router      /products/{id} [delete]
func DeleteProduct(c *gin.Context) {

}

// some helper function to get product with caching
func GetProductWithCache(productID uuid.UUID) (*models.Product, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("product:details:%s", productID.String())

	val, err := config.RDB.Get(ctx, cacheKey).Result()
	if err == nil {
		// Cache Hit! Unmarshal JSON to Struct
		var product models.Product
		if jsonErr := json.Unmarshal([]byte(val), &product); jsonErr == nil {
			return &product, nil
		}
	}

	if err != nil {
		fmt.Println("Redis error:", err)
	}

	product, err := repository.GetProductByUUID(productID)
	if err != nil {
		return nil, err
	}

	// here use goroutine to set cache asynchronouslys
	go func() {
		data, _ := json.Marshal(product)
		config.RDB.Set(ctx, cacheKey, data, 1*time.Hour)
	}()

	return product, nil
}
