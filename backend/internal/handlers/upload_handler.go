package handlers

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"

	"revocreative-cms/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UploadHandler struct {
	cfg *config.Config
}

func NewUploadHandler(cfg *config.Config) *UploadHandler {
	return &UploadHandler{cfg: cfg}
}

// daftar ekstensi yang diizinkan (lowercase, dengan titik)
var allowedExt = map[string]string{
	".png":  "image/png",
	".jpg":  "image/jpeg",
	".jpeg": "image/jpeg",
	".webp": "image/webp",
	".svg":  "image/svg+xml",
	".gif":  "image/gif",
	".ico":  "image/x-icon",
}

// Upload — POST /api/admin/upload (multipart)
// Form field: "file"
// Optional form field: "subdir" (mis. "logos", "portfolio") untuk organisasi
func (h *UploadHandler) Upload(c *gin.Context) {
	// batasi total request size
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, h.cfg.UploadMaxBytes)

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file tidak ditemukan: " + err.Error()})
		return
	}

	if file.Size > h.cfg.UploadMaxBytes {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("file terlalu besar (max %d bytes)", h.cfg.UploadMaxBytes),
		})
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if _, ok := allowedExt[ext]; !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "tipe file tidak diizinkan. Hanya: png, jpg, jpeg, webp, svg, gif, ico",
		})
		return
	}

	if err := validateMimeByContent(file, ext); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// subdir opsional (mis. "logos", "portfolio") supaya rapi
	subdir := strings.TrimSpace(c.PostForm("subdir"))
	subdir = sanitizeSubdir(subdir)

	// generate nama unik
	newName := uuid.New().String() + ext

	// path absolut untuk simpan
	saveDir := h.cfg.UploadDir
	if subdir != "" {
		saveDir = filepath.Join(saveDir, subdir)
	}
	if err := ensureDir(saveDir); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal buat folder: " + err.Error()})
		return
	}

	savePath := filepath.Join(saveDir, newName)
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal simpan file: " + err.Error()})
		return
	}

	// URL publik (relatif) — yang nanti disimpan di DB
	urlPath := "/uploads/"
	if subdir != "" {
		urlPath += subdir + "/"
	}
	urlPath += newName

	c.JSON(http.StatusCreated, gin.H{
		"file_path":     urlPath,
		"original_name": file.Filename,
		"size":          file.Size,
		"mime":          allowedExt[ext],
	})
}

// validateMimeByContent membaca magic bytes file (512 byte pertama) dan
// memastikan content type-nya konsisten dengan ekstensi. Mencegah upload
// file .exe yang di-rename jadi .png.
func validateMimeByContent(fh *multipart.FileHeader, ext string) error {
	f, err := fh.Open()
	if err != nil {
		return fmt.Errorf("gagal buka file: %w", err)
	}
	defer f.Close()

	buf := make([]byte, 512)
	n, _ := f.Read(buf)
	detected := http.DetectContentType(buf[:n])

	expected := allowedExt[ext]
	// SVG kadang ke-detect sebagai text/xml atau text/plain
	if ext == ".svg" && (strings.HasPrefix(detected, "text/") || strings.Contains(detected, "xml")) {
		return nil
	}
	if !strings.HasPrefix(detected, expected) {
		return fmt.Errorf("isi file (%s) tidak cocok dengan ekstensi (%s)", detected, ext)
	}
	return nil
}
