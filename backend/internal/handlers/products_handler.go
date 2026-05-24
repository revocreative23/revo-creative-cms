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

type ProductsHandler struct {
	db *gorm.DB
}

func NewProductsHandler(db *gorm.DB) *ProductsHandler {
	return &ProductsHandler{db: db}
}

type productRequest struct {
	Title         string   `json:"title" binding:"required,min=1,max=255"`
	Slug          string   `json:"slug" binding:"required,max=255"`
	Description   string   `json:"description"`
	ThumbnailPath string   `json:"thumbnail_path"`
	Features      []string `json:"features"`
	Price         string   `json:"price"`
	DisplayOrder  int      `json:"display_order"`
	IsPublished   *bool    `json:"is_published"`
}

func (r *productRequest) applyTo(p *models.Product) error {
	p.Title = r.Title
	p.Slug = r.Slug
	p.Description = r.Description
	p.ThumbnailPath = r.ThumbnailPath
	p.Price = r.Price
	p.DisplayOrder = r.DisplayOrder
	if r.IsPublished != nil {
		p.IsPublished = *r.IsPublished
	}
	if r.Features != nil {
		b, err := json.Marshal(r.Features)
		if err != nil {
			return err
		}
		p.Features = datatypes.JSON(b)
	}
	return nil
}

// ListPublic — GET /api/products
func (h *ProductsHandler) ListPublic(c *gin.Context) {
	var items []models.Product
	if err := h.db.Where("is_published = ?", true).
		Order("display_order asc, id asc").
		Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

// ListAdmin — GET /api/admin/products
func (h *ProductsHandler) ListAdmin(c *gin.Context) {
	var items []models.Product
	if err := h.db.Order("display_order asc, id asc").Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

// GetOne — GET /api/admin/products/:id
func (h *ProductsHandler) GetOne(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id tidak valid"})
		return
	}
	var p models.Product
	if err := h.db.First(&p, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, p)
}

// Create — POST /api/admin/products
func (h *ProductsHandler) Create(c *gin.Context) {
	var req productRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	p := models.Product{IsPublished: true}
	if err := req.applyTo(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.db.Create(&p).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, p)
}

// Update — PUT /api/admin/products/:id
func (h *ProductsHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id tidak valid"})
		return
	}
	var p models.Product
	if err := h.db.First(&p, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "tidak ditemukan"})
		return
	}
	var req productRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := req.applyTo(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.db.Save(&p).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, p)
}

// Delete — DELETE /api/admin/products/:id (soft delete)
func (h *ProductsHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id tidak valid"})
		return
	}
	if err := h.db.Delete(&models.Product{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "product dihapus"})
}
