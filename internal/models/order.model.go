package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type OrderStatus string

const (
	OrderPending   OrderStatus = "pending"
	OrderPaid      OrderStatus = "paid"
	OrderShipped   OrderStatus = "shipped"
	OrderCancelled OrderStatus = "cancelled"
)

type Order struct {
	ID             uuid.UUID       `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID         uuid.UUID       `gorm:"type:uuid;not null;index" json:"user_id"`
	OrderNumber    string          `gorm:"type:varchar(30);not null;uniqueIndex" json:"order_number"`
	Status         OrderStatus     `gorm:"type:order_status;not null;default:'pending'" json:"status"`
	Subtotal       decimal.Decimal `gorm:"type:numeric(10,2);not null" json:"subtotal"`
	DiscountAmount decimal.Decimal `gorm:"type:numeric(10,2);default:0" json:"discount_amount"`
	TaxAmount      decimal.Decimal `gorm:"type:numeric(10,2);default:0" json:"tax_amount"`
	ShippingAmount decimal.Decimal `gorm:"type:numeric(10,2);default:0" json:"shipping_amount"`
	TotalAmount    decimal.Decimal `gorm:"type:numeric(10,2);not null" json:"total_amount"`
	CreatedAt      time.Time       `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt      time.Time       `gorm:"not null;default:now()" json:"updated_at"`
	DeletedAt      gorm.DeletedAt  `gorm:"index" json:"deleted_at" swaggerignore:"true"`

	OrderItems []OrderItem `gorm:"foreignKey:OrderID;references:ID;constraint:OnDelete:CASCADE" json:"order_items"`
	User       User        `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE" json:"user"`
}
