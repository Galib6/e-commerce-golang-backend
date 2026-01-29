package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type OrderItem struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	OrderID   uuid.UUID `gorm:"type:uuid;not null;index"`
	ProductID uuid.UUID `gorm:"type:uuid;not null;index"`

	ProductName string `gorm:"type:varchar(150);not null"`

	// snapshot at purchase time
	ProductPrice    decimal.Decimal `gorm:"type:numeric(10,2);not null"`
	DiscountPercent decimal.Decimal `gorm:"type:numeric(10,2);default:0"`
	Quantity        int             `gorm:"not null;check:quantity > 0"`
	TotalPrice      decimal.Decimal `gorm:"type:numeric(10,2);not null"`

	CreatedAt time.Time `gorm:"type:timestamptz;not null;default:now()"`

	// Relations
	Order   Order   `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE"`
	Product Product `gorm:"foreignKey:ProductID"`

	// Unique constraint
	// uq_order_product (order_id, product_id)
}
