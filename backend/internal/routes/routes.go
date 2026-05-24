package routes

import (
	"net/http"

	"revocreative-cms/internal/config"
	"revocreative-cms/internal/handlers"
	"revocreative-cms/internal/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Register(r *gin.Engine, db *gorm.DB, cfg *config.Config) {
	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "app": "revocreative-cms"})
	})

	authH := handlers.NewAuthHandler(db)
	settingsH := handlers.NewSettingsHandler(db)
	logosH := handlers.NewLogosHandler(db)
	portfolioH := handlers.NewPortfolioHandler(db)
	productsH := handlers.NewProductsHandler(db)
	uploadH := handlers.NewUploadHandler(cfg)

	api := r.Group("/api")
	{
		// --- Public (no auth) ---
		api.POST("/auth/login", authH.Login)
		api.POST("/auth/logout", authH.Logout)

		api.GET("/settings", settingsH.ListPublic)
		api.GET("/logos/active", logosH.GetActive)
		api.GET("/portfolio", portfolioH.ListPublic)
		api.GET("/products", productsH.ListPublic)

		// --- Protected ---
		authed := api.Group("")
		authed.Use(middleware.RequireAuth())
		{
			authed.GET("/auth/me", authH.Me)

			admin := authed.Group("/admin")
			{
				// upload (multipart)
				admin.POST("/upload", uploadH.Upload)

				// settings
				admin.GET("/settings", settingsH.ListAdmin)
				admin.PUT("/settings/:key", settingsH.Update)

				// logos
				admin.GET("/logos", logosH.ListAdmin)
				admin.POST("/logos", logosH.Create)
				admin.PUT("/logos/:id/activate", logosH.Activate)
				admin.DELETE("/logos/:id", logosH.Delete)

				// portfolio
				admin.GET("/portfolio", portfolioH.ListAdmin)
				admin.GET("/portfolio/:id", portfolioH.GetOne)
				admin.POST("/portfolio", portfolioH.Create)
				admin.PUT("/portfolio/:id", portfolioH.Update)
				admin.DELETE("/portfolio/:id", portfolioH.Delete)

				// products
				admin.GET("/products", productsH.ListAdmin)
				admin.GET("/products/:id", productsH.GetOne)
				admin.POST("/products", productsH.Create)
				admin.PUT("/products/:id", productsH.Update)
				admin.DELETE("/products/:id", productsH.Delete)
			}
		}
	}
}
