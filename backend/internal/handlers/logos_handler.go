package handlers

import (
	"net/http"
	"strconv"

	"revocreative-cms/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type LogosHandler struct {
	db *gorm.DB
}

func NewLogosHandler(db *gorm.DB) *LogosHandler {
	return &LogosHandler{db: db}
}

// GetActive — GET /api/logos/active
// Return logo aktif per type ({ "light": {...}, "dark": {...}, "favicon": {...} })
func (h *LogosHandler) GetActive(c *gin.Context) {
	var logos []models.Logo
	if err := h.db.Where("is_active = ?", true).Find(&logos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	result := make(map[string]models.Logo)
	for _, l := range logos {
		result[l.Type] = l
	}
	c.JSON(http.StatusOK, result)
}

// ListAdmin — GET /api/admin/logos
// Return semua logo (semua type, semua history) untuk panel admin.
func (h *LogosHandler) ListAdmin(c *gin.Context) {
	var logos []models.Logo
	q := h.db.Order("type asc, created_at desc")
	if t := c.Query("type"); t != "" {
		q = q.Where("type = ?", t)
	}
	if err := q.Find(&logos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, logos)
}

type createLogoRequest struct {
	Type     string `json:"type" binding:"required,oneof=light dark favicon"`
	FilePath string `json:"file_path" binding:"required"`
	Activate bool   `json:"activate"` // kalau true, langsung jadi aktif (non-aktif-kan yang lain dengan type sama)
}

// Create — POST /api/admin/logos
// (Untuk Tahap 5: ini akan dipanggil setelah file upload sukses.
// Untuk sekarang, frontend kirim path file yang sudah ada di /uploads)
func (h *LogosHandler) Create(c *gin.Context) {
	var req createLogoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logo := models.Logo{
		Type:     req.Type,
		FilePath: req.FilePath,
		IsActive: false,
	}

	if err := h.db.Create(&logo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if req.Activate {
		if err := h.activate(&logo); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusCreated, logo)
}

// Activate — PUT /api/admin/logos/:id/activate
// Aktifkan logo ini, non-aktifkan semua logo lain dengan type yang sama.
func (h *LogosHandler) Activate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id tidak valid"})
		return
	}

	var logo models.Logo
	if err := h.db.First(&logo, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "logo tidak ditemukan"})
		return
	}

	if err := h.activate(&logo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, logo)
}

func (h *LogosHandler) activate(logo *models.Logo) error {
	return h.db.Transaction(func(tx *gorm.DB) error {
		// non-aktifkan semua logo lain dengan type sama
		if err := tx.Model(&models.Logo{}).
			Where("type = ? AND id <> ?", logo.Type, logo.ID).
			Update("is_active", false).Error; err != nil {
			return err
		}
		// aktifkan logo ini
		logo.IsActive = true
		return tx.Save(logo).Error
	})
}

// Delete — DELETE /api/admin/logos/:id (soft delete)
func (h *LogosHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id tidak valid"})
		return
	}

	if err := h.db.Delete(&models.Logo{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "logo dihapus"})
}
