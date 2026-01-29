package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductImages struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ProductId uuid.UUID `gorm:"type:uuid;not null;index"`
	ImageUrl  string    `gorm:"type:text;not null"`
	IsPrimary bool      `gorm:"type:boolean;default:false"`
	SortOrder int       `gorm:"type:integer;default:0"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at" swaggerignore:"true"`
	Product   Product        `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE"`
}
