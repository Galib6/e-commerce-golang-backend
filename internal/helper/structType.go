package helper

import (
	"mime/multipart"

	"github.com/shopspring/decimal"
)

type ProductImageRequest struct {
	ImageURL  string `json:"imageUrl" validate:"required,url"`
	IsPrimary bool   `json:"isPrimary"`
	SortOrder int    `json:"sortOrder"`
}

// type CreateProductRequest struct {
// 	Name            string                  `json:"name" binding:"required,min=3,max=100"`
// 	Description     string                  `json:"description" binding:"omitempty,max=2000"`
// 	BasePrice       decimal.Decimal         `json:"price" binding:"required,gt=0"`
// 	DiscountPercent decimal.Decimal         `json:"discount_percent"`
// }

// if you are using form-data for file upload
type CreateProductRequest struct {
	Name            string                  `form:"name" validate:"required"`
	Description     string                  `form:"description" validate:"max=2000"`
	BasePrice       decimal.Decimal         `form:"base_price" validate:"required"`
	DiscountPercent decimal.Decimal         `form:"discount_percent" validate:"gte=0,lte=100"`
	ImageFiles      []*multipart.FileHeader `form:"images"`
}
type UserFilterParams struct {
	ProductName string
	FullName    string
	Page        int
	Limit       int
}
