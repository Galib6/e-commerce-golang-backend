package helper

import (
	"errors"

	"github.com/shopspring/decimal"
)

// CustomValidate validates business rules for Product DTO
func CustomValidate(ctx *CreateProductRequest) error {
	// Default discount percent set to 0
	discountPercent := ctx.DiscountPercent
	if discountPercent.IsZero() {
		discountPercent = decimal.Zero
	}

	discountPrice := ctx.BasePrice.
		Mul(discountPercent).
		Div(decimal.NewFromInt(100))

	if discountPrice.GreaterThanOrEqual(ctx.BasePrice) {
		return errors.New("discount price must be less than base price")
	}

	return nil
}
