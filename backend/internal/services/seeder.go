package services

import (
	"log"
	"os"

	"revocreative-cms/internal/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// SeedAll menjalankan semua seeder. Idempotent — aman dipanggil berkali-kali.
func SeedAll(db *gorm.DB) error {
	if err := seedAdminUser(db); err != nil {
		return err
	}
	if err := seedSiteSettings(db); err != nil {
		return err
	}
	return nil
}

func seedAdminUser(db *gorm.DB) error {
	email := os.Getenv("SEED_ADMIN_EMAIL")
	password := os.Getenv("SEED_ADMIN_PASSWORD")
	if email == "" || password == "" {
		log.Println("⚠ SEED_ADMIN_EMAIL / SEED_ADMIN_PASSWORD tidak di-set, skip seed admin")
		return nil
	}

	var count int64
	if err := db.Model(&models.User{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		log.Printf("✓ admin user %s sudah ada, skip", email)
		return nil
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	admin := models.User{
		Email:        email,
		PasswordHash: string(hash),
		Name:         "Administrator",
		Role:         "admin",
	}
	if err := db.Create(&admin).Error; err != nil {
		return err
	}
	log.Printf("✓ admin user %s berhasil dibuat", email)
	return nil
}

func seedSiteSettings(db *gorm.DB) error {
	defaults := []models.SiteSetting{
		{Key: "company_name", Value: "Revo Creative", Description: "Nama brand yang ditampilkan di header/footer"},
		{Key: "company_legal_name", Value: "PT Rajawali Cakra Digdaya", Description: "Nama legal perusahaan"},
		{Key: "company_address", Value: "Jl. Raya Pd. Gede No.14A, Pinang Ranti, Kec. Makasar, Jakarta Timur 13560", Description: "Alamat perusahaan (footer & contact)"},
		{Key: "company_phone", Value: "+62 856-7990-037", Description: "Nomor telepon utama"},
		{Key: "company_email", Value: "rajawalicakradigdaya@gmail.com", Description: "Email perusahaan"},
		{Key: "company_whatsapp", Value: "628567990037", Description: "Nomor WhatsApp (format internasional tanpa +)"},
		{Key: "social_instagram", Value: "https://www.instagram.com/revocreative_id/", Description: "URL Instagram"},
		{Key: "social_linkedin", Value: "https://www.linkedin.com/company/101702844", Description: "URL LinkedIn"},
		{Key: "footer_tagline", Value: "Membangun digital experience yang berdampak.", Description: "Tagline pendek di footer"},
	}

	for _, s := range defaults {
		var existing models.SiteSetting
		err := db.Where("key = ?", s.Key).First(&existing).Error
		if err == gorm.ErrRecordNotFound {
			if err := db.Create(&s).Error; err != nil {
				return err
			}
			log.Printf("✓ seed setting: %s", s.Key)
		} else if err != nil {
			return err
		}
	}
	return nil
}
