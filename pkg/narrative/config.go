// Package narrative provides a reusable client for PgQueryNarrative capabilities:
// running queries, generating reports, and exposing Goa service implementations
// for use by the standalone server or embedded applications.
package narrative

import (
	"time"

	"github.com/pgquerynarrative/pgquerynarrative/app/config"
)

// Config holds configuration for the narrative client. It can be built from
// environment (via app/config.Load) or supplied in code for library usage.
type Config struct {
	// Database holds PostgreSQL connection settings for both read-only and app pools.
	Database DatabaseConfig
	// LLM holds the LLM provider settings for narrative generation.
	LLM LLMConfig
	// Metrics holds trend threshold for period comparison (e.g. 0.5 for 0.5%).
	Metrics MetricsConfig
	// Embedding holds optional settings for RAG and similar-query retrieval. When BaseURL is empty, embeddings are disabled.
	Embedding EmbeddingConfig
	// AllowedSchemas is the list of schema names queries may access (e.g. []string{"demo"}).
	AllowedSchemas []string
	// MaxQueryLength is the maximum allowed query length in bytes.
	MaxQueryLength int
	// MaxRowsPerQuery is the maximum rows returned per query execution.
	MaxRowsPerQuery int
}

// EmbeddingConfig holds optional embedding model settings (e.g. Ollama nomic-embed-text).
type EmbeddingConfig struct {
	BaseURL string
	Model   string
}

// DatabaseConfig holds PostgreSQL connection settings.
type DatabaseConfig struct {
	Host             string
	Port             int
	Database         string
	User             string
	Password         string
	MaxConnections   int
	ReadOnlyUser     string
	ReadOnlyPassword string
	SSLMode          string
	QueryTimeout     time.Duration
}

// LLMConfig holds LLM provider settings.
type LLMConfig struct {
	Provider string // "ollama", "gemini", "claude", "openai", "groq"
	Model    string
	APIKey   string
	BaseURL  string
}

// MetricsConfig holds metrics and period-comparison settings.
type MetricsConfig struct {
	TrendThresholdPercent float64
}

// FromAppConfig converts app config into narrative config with default
// allowed schemas and limits. Use this when building a client from
// config.Load() in the standalone server.
func FromAppConfig(cfg config.Config) Config {
	return Config{
		Database: DatabaseConfig{
			Host:             cfg.Database.Host,
			Port:             cfg.Database.Port,
			Database:         cfg.Database.Database,
			User:             cfg.Database.User,
			Password:         cfg.Database.Password,
			MaxConnections:   cfg.Database.MaxConnections,
			ReadOnlyUser:     cfg.Database.ReadOnlyUser,
			ReadOnlyPassword: cfg.Database.ReadOnlyPassword,
			SSLMode:          cfg.Database.SSLMode,
			QueryTimeout:     cfg.Database.QueryTimeout,
		},
		LLM: LLMConfig{
			Provider: cfg.LLM.Provider,
			Model:    cfg.LLM.Model,
			APIKey:   cfg.LLM.APIKey,
			BaseURL:  cfg.LLM.BaseURL,
		},
		Metrics: MetricsConfig{
			TrendThresholdPercent: cfg.Metrics.TrendThresholdPercent,
		},
		Embedding: EmbeddingConfig{
			BaseURL: cfg.Embedding.BaseURL,
			Model:   cfg.Embedding.Model,
		},
		AllowedSchemas:  allowedSchemasOrDefault(cfg.Database.AllowedSchemas),
		MaxQueryLength:  10000,
		MaxRowsPerQuery: 1000,
	}
}

func allowedSchemasOrDefault(schemas []string) []string {
	if len(schemas) > 0 {
		return schemas
	}
	return []string{"public", "demo"}
}
