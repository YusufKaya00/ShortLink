package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/urlshortener/url-service/internal/handlers"
	"github.com/urlshortener/url-service/internal/middleware"
	"github.com/urlshortener/url-service/internal/models"
	"github.com/urlshortener/url-service/internal/repository"
	"github.com/urlshortener/url-service/internal/service"
	"github.com/urlshortener/url-service/pkg/redis"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Database connection
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:postgres123@localhost:5432/urlshortener?sslmode=disable"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto migrate
	if err := db.AutoMigrate(&models.URL{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Initialize Redis
	redisClient := redis.NewRedisClient()

	// Initialize layers
	urlRepo := repository.NewURLRepository(db)
	urlService := service.NewURLService(urlRepo, redisClient)
	urlHandler := handlers.NewURLHandler(urlService)

	// Setup Gin
	r := gin.Default()

	// CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// Routes
	r.GET("/health", urlHandler.Health)

	// Redirect route (public)
	r.GET("/:code", urlHandler.Redirect)

	api := r.Group("/api/urls")
	{
		// Public routes
		api.GET("/all", urlHandler.GetAllURLs)
		api.GET("/info/:code", urlHandler.GetURL)

		// Optional auth for creating URLs (works with or without auth)
		api.POST("", middleware.OptionalAuthMiddleware(), urlHandler.CreateURL)

		// Protected routes
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.GET("", urlHandler.GetUserURLs)
			protected.DELETE("/:id", urlHandler.DeleteURL)
		}
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	log.Printf("URL Service starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
