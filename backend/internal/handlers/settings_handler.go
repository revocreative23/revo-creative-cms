package handlers

import (
	"net/http"

	"revocreative-cms/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SettingsHandler struct {
	db *gorm.DB
}

func NewSettingsHandler(db *gorm.DB) *SettingsHandler {
	return &SettingsHandler{db: db}
}

// ListPublic — GET /api/settings
// Return semua settings sebagai map { key: value } untuk dikonsumsi frontend publik.
func (h *SettingsHandler) ListPublic(c *gin.Context) {
	var settings []models.SiteSetting
	if err := h.db.Find(&settings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	result := make(map[string]string, len(settings))
	for _, s := range settings {
		result[s.Key] = s.Value
	}
	c.JSON(http.StatusOK, result)
}

// ListAdmin — GET /api/admin/settings
// Return semua settings lengkap (dengan description, updated_at) untuk panel admin.
func (h *SettingsHandler) ListAdmin(c *gin.Context) {
	var settings []models.SiteSetting
	if err := h.db.Order("key asc").Find(&settings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, settings)
}

type updateSettingRequest struct {
	Value string `json:"value"`
}

// Update — PUT /api/admin/settings/:key
func (h *SettingsHandler) Update(c *gin.Context) {
	key := c.Param("key")
	var req updateSettingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var setting models.SiteSetting
	if err := h.db.Where("key = ?", key).First(&setting).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "setting tidak ditemukan"})
		return
	}

	setting.Value = req.Value
	if err := h.db.Save(&setting).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, setting)
}
