package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type ProductStatus string

const (
	ProductDraft    ProductStatus = "draft"
	ProductActive   ProductStatus = "active"
	ProductInactive ProductStatus = "inactive"
	ProductArchived ProductStatus = "archived"
)

type Product struct {
	ID               uuid.UUID       `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name             string          `gorm:"size:100;not null"`
	ShortDescription string          `gorm:"type:text"`
	BasePrice        decimal.Decimal `gorm:"type:numeric(10,2);not null"`
	DiscountPercent  decimal.Decimal `gorm:"type:numeric(10,2);default:0;check:discount_percent >= 0 AND discount_percent <= 100"`
	Currency         string          `gorm:"type:char(3);default:'INR'"`
	IsReturnable     bool            `gorm:"type:boolean;default:true"`
	IsCodAvailable   bool            `gorm:"type:boolean;default:true"`
	NumberOfStock    int             `gorm:"type:integer;default:0;check:number_of_stock >= 0"`
	Status           ProductStatus   `gorm:"type:product_status;default:'draft'"`
	CreatedBy        uuid.UUID       `gorm:"type:uuid;not null" json:"created_by"`
	User             User            `gorm:"foreignKey:CreatedBy"`

	CreatedAt time.Time      `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt time.Time      `gorm:"not null;default:now()" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at" swaggerignore:"true"`

	ProductImages []ProductImages `gorm:"foreignKey:ProductID"`
}
