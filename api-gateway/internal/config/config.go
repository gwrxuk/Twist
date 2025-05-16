package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Environment string
	Server      ServerConfig
	Database    DatabaseConfig
	Redis       RedisConfig
	JWT         JWTConfig
	Services    ServicesConfig
}

type ServerConfig struct {
	Host string
	Port int
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

type JWTConfig struct {
	Secret        string
	ExpiryMinutes int
}

type ServicesConfig struct {
	CoreEngine    ServiceConfig
	SmartContract ServiceConfig
	Analytics     ServiceConfig
}

type ServiceConfig struct {
	Host string
	Port int
}

func LoadConfig() (*Config, error) {
	// Set default config values
	viper.SetDefault("environment", "development")
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.port", 8000)
	viper.SetDefault("database.sslmode", "disable")
	viper.SetDefault("jwt.expiry_minutes", 60)

	// Read configuration from environment variables
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// Map environment variables to config
	// Server
	mapEnvToConfig("SERVER_HOST", "server.host")
	mapEnvToConfig("SERVER_PORT", "server.port")

	// Database
	mapEnvToConfig("DB_HOST", "database.host")
	mapEnvToConfig("DB_PORT", "database.port")
	mapEnvToConfig("DB_USER", "database.user")
	mapEnvToConfig("DB_PASSWORD", "database.password")
	mapEnvToConfig("DB_NAME", "database.dbname")
	mapEnvToConfig("DB_SSLMODE", "database.sslmode")

	// Redis
	mapEnvToConfig("REDIS_HOST", "redis.host")
	mapEnvToConfig("REDIS_PORT", "redis.port")
	mapEnvToConfig("REDIS_PASSWORD", "redis.password")
	mapEnvToConfig("REDIS_DB", "redis.db")

	// JWT
	mapEnvToConfig("JWT_SECRET", "jwt.secret")
	mapEnvToConfig("JWT_EXPIRY_MINUTES", "jwt.expiry_minutes")

	// Services
	mapEnvToConfig("CORE_ENGINE_HOST", "services.core_engine.host")
	mapEnvToConfig("CORE_ENGINE_PORT", "services.core_engine.port")
	mapEnvToConfig("SMART_CONTRACT_HOST", "services.smart_contract.host")
	mapEnvToConfig("SMART_CONTRACT_PORT", "services.smart_contract.port")
	mapEnvToConfig("ANALYTICS_HOST", "services.analytics.host")
	mapEnvToConfig("ANALYTICS_PORT", "services.analytics.port")

	// Validate required fields
	if err := validateRequiredConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

func mapEnvToConfig(envVar, configPath string) {
	if val, ok := os.LookupEnv(envVar); ok {
		viper.Set(configPath, val)
	}
}

func validateRequiredConfig() error {
	requiredFields := []string{
		"database.host",
		"database.port",
		"database.user",
		"database.password",
		"database.dbname",
		"jwt.secret",
	}

	var missingFields []string
	for _, field := range requiredFields {
		if !viper.IsSet(field) {
			missingFields = append(missingFields, field)
		}
	}

	if len(missingFields) > 0 {
		return fmt.Errorf("missing required configuration: %s", strings.Join(missingFields, ", "))
	}

	return nil
}
