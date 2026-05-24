package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Product struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	Title         string         `gorm:"size:255;not null" json:"title"`
	Slug          string         `gorm:"uniqueIndex;size:255" json:"slug"`
	Description   string         `gorm:"type:text" json:"description"`
	ThumbnailPath string         `gorm:"size:500" json:"thumbnail_path"`
	Features      datatypes.JSON `gorm:"type:jsonb" json:"features"` // ["Custom Domain", "CMS", ...]
	Price         string         `gorm:"size:100" json:"price"`      // string supaya bisa "Mulai Rp 3jt"
	DisplayOrder  int            `gorm:"not null;default:0;index" json:"display_order"`
	IsPublished   bool           `gorm:"not null;default:true;index" json:"is_published"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}
