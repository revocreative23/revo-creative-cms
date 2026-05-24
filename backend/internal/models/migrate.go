package models

import (
	"log"

	"gorm.io/gorm"
)

// AllModels mengembalikan semua model untuk di-migrate.
// Tambahkan model baru di sini saat dibuat.
func AllModels() []interface{} {
	return []interface{}{
		&User{},
		&SiteSetting{},
		&Logo{},
		&PortfolioItem{},
		&Product{},
	}
}

func AutoMigrate(db *gorm.DB) error {
	log.Println("→ menjalankan AutoMigrate...")
	if err := db.AutoMigrate(AllModels()...); err != nil {
		return err
	}
	log.Println("✓ semua tabel ter-migrate")
	return nil
}
