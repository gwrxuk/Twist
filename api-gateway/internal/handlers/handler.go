package handlers

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/twist/api-gateway/internal/config"
	"go.uber.org/zap"
)

// Handler holds all the dependencies for our handlers
type Handler struct {
	db          *pgxpool.Pool
	redisClient *redis.Client
	logger      *zap.Logger
	config      *config.Config
}

// NewHandler creates a new Handler instance
func NewHandler(db *pgxpool.Pool, redisClient *redis.Client, logger *zap.Logger, config *config.Config) *Handler {
	return &Handler{
		db:          db,
		redisClient: redisClient,
		logger:      logger,
		config:      config,
	}
}
