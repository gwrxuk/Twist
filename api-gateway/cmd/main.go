package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/twist/api-gateway/internal/config"
	"github.com/twist/api-gateway/internal/handlers"
	"github.com/twist/api-gateway/internal/middleware"
	"github.com/twist/api-gateway/pkg/database"
	"github.com/twist/api-gateway/pkg/logger"
	"github.com/twist/api-gateway/pkg/metrics"

	"go.uber.org/zap"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		fmt.Printf("Warning: failed to load .env file: %v\n", err)
	}

	// Initialize logger
	log := logger.NewLogger()
	defer log.Sync()

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load configuration", zap.Error(err))
	}

	// Initialize DB connection
	db, err := database.NewPostgresDB(cfg.Database)
	if err != nil {
		log.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	// Initialize Redis
	redisClient, err := database.NewRedisClient(cfg.Redis)
	if err != nil {
		log.Fatal("Failed to connect to Redis", zap.Error(err))
	}
	defer redisClient.Close()

	// Initialize metrics
	metricsClient := metrics.NewPrometheusClient()

	// Set Gin mode
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize Gin router
	router := gin.New()

	// Apply middleware
	router.Use(gin.Recovery())
	router.Use(middleware.Logger(log))
	router.Use(middleware.Metrics(metricsClient))

	// Initialize handlers
	h := handlers.NewHandler(db, redisClient, log, cfg)

	// Set up API routes
	api := router.Group("/api/v1")
	{
		// Public routes
		api.GET("/health", h.HealthCheck)

		// Auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/login", h.Login)
			auth.POST("/register", h.Register)
		}

		// Protected routes
		protected := api.Group("/")
		protected.Use(middleware.Auth(cfg.JWT.Secret))
		{
			// Node management
			nodes := protected.Group("/nodes")
			{
				nodes.GET("", h.ListNodes)
				nodes.GET("/:id", h.GetNode)
				nodes.POST("", h.CreateNode)
				nodes.PUT("/:id", h.UpdateNode)
				nodes.DELETE("/:id", h.DeleteNode)
			}

			// User management
			users := protected.Group("/users")
			{
				users.GET("/me", h.GetCurrentUser)
				users.PUT("/me", h.UpdateCurrentUser)
			}

			// API key management
			apiKeys := protected.Group("/api-keys")
			{
				apiKeys.GET("", h.ListAPIKeys)
				apiKeys.POST("", h.CreateAPIKey)
				apiKeys.DELETE("/:id", h.DeleteAPIKey)
			}
		}
	}

	// Metrics endpoint
	router.GET("/metrics", h.Metrics)

	// Start HTTP server
	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler: router,
	}

	// Start the server in a goroutine
	go func() {
		log.Info("Starting API Gateway server", zap.String("address", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	// Create a deadline for server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown", zap.Error(err))
	}

	log.Info("Server exiting")
}
