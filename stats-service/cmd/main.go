package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/urlshortener/stats-service/internal/consumer"
	"github.com/urlshortener/stats-service/internal/handlers"
	"github.com/urlshortener/stats-service/internal/models"
	"github.com/urlshortener/stats-service/internal/repository"
	"github.com/urlshortener/stats-service/internal/service"
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
	if err := db.AutoMigrate(&models.Click{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Initialize Redis
	redisAddr := os.Getenv("REDIS_URL")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	// Initialize layers
	statsRepo := repository.NewStatsRepository(db)
	statsService := service.NewStatsService(statsRepo)
	statsHandler := handlers.NewStatsHandler(statsService)

	// Start consumer in background
	ctx, cancel := context.WithCancel(context.Background())
	clickConsumer := consumer.NewClickConsumer(redisClient, statsService)
	go clickConsumer.Start(ctx)

	// Setup Gin
	r := gin.Default()

	// CORS
	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
	}))

	// Routes
	r.GET("/health", statsHandler.Health)

	api := r.Group("/api/stats")
	{
		api.GET("/overall", statsHandler.GetOverallStats)
		api.GET("/recent", statsHandler.GetRecentClicks)
		api.GET("/:code", statsHandler.GetURLStats)
	}

	// Graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		log.Println("Shutting down...")
		cancel()
		redisClient.Close()
		os.Exit(0)
	}()

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8083"
	}

	log.Printf("Stats Service starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
