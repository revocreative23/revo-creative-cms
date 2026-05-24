package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"revocreative-cms/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type PortfolioHandler struct {
	db *gorm.DB
}

func NewPortfolioHandler(db *gorm.DB) *PortfolioHandler {
	return &PortfolioHandler{db: db}
}

type portfolioRequest struct {
	Title         string   `json:"title" binding:"required,min=1,max=255"`
	Category      string   `json:"category" binding:"required,max=100"`
	CategoryLabel string   `json:"category_label"`
	Description   string   `json:"description"`
	ThumbnailPath string   `json:"thumbnail_path"`
	Tags          []string `json:"tags"`
	DisplayOrder  int      `json:"display_order"`
	IsPublished   *bool    `json:"is_published"`
}

func (r *portfolioRequest) applyTo(item *models.PortfolioItem) error {
	item.Title = r.Title
	item.Category = r.Category
	item.CategoryLabel = r.CategoryLabel
	item.Description = r.Description
	item.ThumbnailPath = r.ThumbnailPath
	item.DisplayOrder = r.DisplayOrder
	if r.IsPublished != nil {
		item.IsPublished = *r.IsPublished
	}
	if r.Tags != nil {
		b, err := json.Marshal(r.Tags)
		if err != nil {
			return err
		}
		item.Tags = datatypes.JSON(b)
	}
	return nil
}

// ListPublic — GET /api/portfolio
// Hanya yang published, urut display_order.
func (h *PortfolioHandler) ListPublic(c *gin.Context) {
	var items []models.PortfolioItem
	q := h.db.Where("is_published = ?", true).Order("display_order asc, id asc")
	if cat := c.Query("category"); cat != "" {
		q = q.Where("category = ?", cat)
	}
	if err := q.Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

// ListAdmin — GET /api/admin/portfolio (semua, termasuk unpublished)
func (h *PortfolioHandler) ListAdmin(c *gin.Context) {
	var items []models.PortfolioItem
	if err := h.db.Order("display_order asc, id asc").Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

// GetOne — GET /api/admin/portfolio/:id
func (h *PortfolioHandler) GetOne(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id tidak valid"})
		return
	}
	var item models.PortfolioItem
	if err := h.db.First(&item, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, item)
}

// Create — POST /api/admin/portfolio
func (h *PortfolioHandler) Create(c *gin.Context) {
	var req portfolioRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item := models.PortfolioItem{IsPublished: true}
	if err := req.applyTo(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.Create(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, item)
}

// Update — PUT /api/admin/portfolio/:id
func (h *PortfolioHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id tidak valid"})
		return
	}

	var item models.PortfolioItem
	if err := h.db.First(&item, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "tidak ditemukan"})
		return
	}

	var req portfolioRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := req.applyTo(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.Save(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

// Delete — DELETE /api/admin/portfolio/:id (soft delete)
func (h *PortfolioHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id tidak valid"})
		return
	}
	if err := h.db.Delete(&models.PortfolioItem{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "portfolio dihapus"})
}
