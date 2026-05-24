package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type PortfolioItem struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	Title         string         `gorm:"size:255;not null" json:"title"`
	Category      string         `gorm:"size:100;not null;index" json:"category"` // website, app, dashboard, dll
	CategoryLabel string         `gorm:"size:255" json:"category_label"`          // mis. "Web + CMS + Membership"
	Description   string         `gorm:"type:text" json:"description"`
	ThumbnailPath string         `gorm:"size:500" json:"thumbnail_path"`
	Tags          datatypes.JSON `gorm:"type:jsonb" json:"tags"` // ["ISI", "Membership System", ...]
	DisplayOrder  int            `gorm:"not null;default:0;index" json:"display_order"`
	IsPublished   bool           `gorm:"not null;default:true;index" json:"is_published"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}
