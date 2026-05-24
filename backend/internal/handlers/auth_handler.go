package handlers

import (
	"net/http"

	"revocreative-cms/internal/middleware"
	"revocreative-cms/internal/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthHandler struct {
	db *gorm.DB
}

func NewAuthHandler(db *gorm.DB) *AuthHandler {
	return &AuthHandler{db: db}
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type userResponse struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Role  string `json:"role"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := h.db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		// generic error supaya tidak bocor info user ada/tidak
		c.JSON(http.StatusUnauthorized, gin.H{"error": "email atau password salah"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "email atau password salah"})
		return
	}

	s := sessions.Default(c)
	s.Set(middleware.SessionKeyUserID, user.ID)
	s.Set(middleware.SessionKeyEmail, user.Email)
	s.Set(middleware.SessionKeyRole, user.Role)
	if err := s.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal simpan session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": userResponse{ID: user.ID, Email: user.Email, Name: user.Name, Role: user.Role},
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	s := sessions.Default(c)
	s.Clear()
	s.Options(sessions.Options{Path: "/", MaxAge: -1})
	if err := s.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal hapus session"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "logout berhasil"})
}

func (h *AuthHandler) Me(c *gin.Context) {
	uid := c.GetUint("user_id")
	var user models.User
	if err := h.db.First(&user, uid).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"user": userResponse{ID: user.ID, Email: user.Email, Name: user.Name, Role: user.Role},
	})
}
