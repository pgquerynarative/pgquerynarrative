// Package config provides configuration management for the application.
// It loads configuration from environment variables with sensible defaults.
package config

import (
	"os"
	"strconv"
	"time"
)

// Config holds all application configuration.
type Config struct {
	Server   ServerConfig   // HTTP server configuration
	Database DatabaseConfig // Database connection settings
	Security SecurityConfig // Security settings
	LLM      LLMConfig      // LLM provider configuration
	Metrics  MetricsConfig  // Metrics and period-comparison settings
}

// MetricsConfig holds settings for metrics and period-over-period comparison.
type MetricsConfig struct {
	// TrendThresholdPercent is the minimum absolute % change to label as "up" or "down" (vs "flat"). Default 0.5.
	TrendThresholdPercent float64
}

// ServerConfig contains HTTP server settings.
type ServerConfig struct {
	Host         string        // Server host (default: "0.0.0.0")
	Port         int           // Server port (default: 8080)
	ReadTimeout  time.Duration // Request read timeout
	WriteTimeout time.Duration // Response write timeout
}

// DatabaseConfig contains PostgreSQL connection settings.
type DatabaseConfig struct {
	Host             string        // Database host
	Port             int           // Database port
	Database         string        // Database name
	User             string        // Application user (full access)
	Password         string        // Application user password
	MaxConnections   int           // Maximum connection pool size
	ReadOnlyUser     string        // Read-only user (for queries)
	ReadOnlyPassword string        // Read-only user password
	SSLMode          string        // SSL mode (disable, require, etc.)
	QueryTimeout     time.Duration // Maximum query execution time
}

// SecurityConfig contains security-related settings.
type SecurityConfig struct {
	AuthEnabled bool // Whether authentication is enabled (future feature)
}

// LLMConfig contains settings for the LLM provider used for narrative generation.
type LLMConfig struct {
	Provider string // Provider name: "ollama", "gemini", "claude", "openai", "groq"
	Model    string // Model name (e.g., "llama3.2", "gpt-4o-mini")
	APIKey   string // API key (for cloud providers)
	BaseURL  string // Base URL (for local providers like Ollama)
}

// Load reads configuration from environment variables and returns a Config struct.
// All settings have sensible defaults for local development.
func Load() Config {
	return Config{
		Server: ServerConfig{
			Host:         getEnv("PGQUERYNARRATIVE_HOST", "0.0.0.0"),
			Port:         getEnvInt("PGQUERYNARRATIVE_PORT", 8080),
			ReadTimeout:  getEnvDuration("PGQUERYNARRATIVE_READ_TIMEOUT", 15*time.Second),
			WriteTimeout: getEnvDuration("PGQUERYNARRATIVE_WRITE_TIMEOUT", 60*time.Second),
		},
		Database: DatabaseConfig{
			Host:             getEnv("DATABASE_HOST", "localhost"),
			Port:             getEnvInt("DATABASE_PORT", 5432),
			Database:         getEnv("DATABASE_NAME", "pgquerynarrative"),
			User:             getEnv("DATABASE_USER", "pgquerynarrative_app"),
			Password:         getEnv("DATABASE_PASSWORD", "pgquerynarrative_app"),
			MaxConnections:   getEnvInt("DATABASE_MAX_CONNECTIONS", 10),
			ReadOnlyUser:     getEnv("DATABASE_READONLY_USER", "pgquerynarrative_readonly"),
			ReadOnlyPassword: getEnv("DATABASE_READONLY_PASSWORD", "pgquerynarrative_readonly"),
			SSLMode:          getEnv("DATABASE_SSL_MODE", "disable"),
			QueryTimeout:     getEnvDuration("QUERY_TIMEOUT", 30*time.Second),
		},
		Security: SecurityConfig{
			AuthEnabled: getEnvBool("SECURITY_AUTH_ENABLED", false),
		},
		LLM: LLMConfig{
			Provider: getEnv("LLM_PROVIDER", "ollama"),
			Model:    getEnv("LLM_MODEL", "llama3.2"),
			APIKey:   getEnv("LLM_API_KEY", ""),
			BaseURL:  getEnv("LLM_BASE_URL", "http://localhost:11434"),
		},
		Metrics: MetricsConfig{
			TrendThresholdPercent: getEnvFloat("PERIOD_TREND_THRESHOLD_PERCENT", 0.5),
		},
	}
}

// getEnvFloat retrieves a float environment variable or returns a default value.
func getEnvFloat(key string, fallback float64) float64 {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseFloat(value, 64); err == nil {
			return parsed
		}
	}
	return fallback
}

// getEnv retrieves an environment variable or returns a default value.
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// getEnvInt retrieves an integer environment variable or returns a default value.
// If the value cannot be parsed as an integer, the default is returned.
func getEnvInt(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			return parsed
		}
	}
	return fallback
}

// getEnvBool retrieves a boolean environment variable or returns a default value.
// Accepts: "true", "1", "yes", "on" (case-insensitive) for true.
func getEnvBool(key string, fallback bool) bool {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseBool(value); err == nil {
			return parsed
		}
	}
	return fallback
}

// getEnvDuration retrieves a duration environment variable or returns a default value.
// Accepts formats like "30s", "5m", "1h", etc.
func getEnvDuration(key string, fallback time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if parsed, err := time.ParseDuration(value); err == nil {
			return parsed
		}
	}
	return fallback
}
