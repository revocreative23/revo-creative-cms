package models

import (
	"time"

	"gorm.io/gorm"
)

const (
	LogoTypeLight   = "light"
	LogoTypeDark    = "dark"
	LogoTypeFavicon = "favicon"
)

type Logo struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Type      string         `gorm:"size:20;not null;index" json:"type"` // light, dark, favicon
	FilePath  string         `gorm:"size:500;not null" json:"file_path"`
	IsActive  bool           `gorm:"not null;default:false;index" json:"is_active"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
