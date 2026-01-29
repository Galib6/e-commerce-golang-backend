package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID       uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Fullname string    `gorm:"size:50;not null"`
	Username string    `gorm:"size:50;unique;not null"`
	Email    string    `gorm:"size:100;unique;not null"`
	Password string    `gorm:"size:255;not null" json:"-"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at" swaggerignore:"true"`

	Orders []Order `gorm:"foreignKey:UserID"`
}
