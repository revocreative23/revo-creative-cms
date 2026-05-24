package main

import (
	"log"
	"net/http"

	"revocreative-cms/internal/config"
	"revocreative-cms/internal/models"
	"revocreative-cms/internal/routes"
	"revocreative-cms/internal/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	db := config.NewDB(cfg)

	if err := models.AutoMigrate(db); err != nil {
		log.Fatalf("auto-migrate gagal: %v", err)
	}
	if err := services.SeedAll(db); err != nil {
		log.Fatalf("seed gagal: %v", err)
	}

	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// CORS — wajib AllowCredentials true supaya cookie session jalan dari SPA
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{cfg.CORSAllowedOrigin},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}))

	// Session — cookie httpOnly, di-encrypt pakai SESSION_SECRET
	store := cookie.NewStore([]byte(cfg.SessionSecret))
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   cfg.SessionMaxAge,
		HttpOnly: true,
		Secure:   cfg.AppEnv == "production", // di prod wajib HTTPS
		SameSite: http.SameSiteLaxMode,
	})
	r.Use(sessions.Sessions(cfg.SessionCookieName, store))

	// Serve uploaded files publik di /uploads/*
	r.Static("/uploads", cfg.UploadDir)

	routes.Register(r, db, cfg)

	addr := ":" + cfg.ServerPort
	log.Printf("server jalan di http://localhost%s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("server gagal start: %v", err)
	}
}
