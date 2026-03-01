package narrative

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pgquerynarrative/pgquerynarrative/api/gen/queries"
	"github.com/pgquerynarrative/pgquerynarrative/api/gen/reports"
	schema "github.com/pgquerynarrative/pgquerynarrative/api/gen/schema"
	suggestions "github.com/pgquerynarrative/pgquerynarrative/api/gen/suggestions"
	"github.com/pgquerynarrative/pgquerynarrative/app/catalog"
	"github.com/pgquerynarrative/pgquerynarrative/app/config"
	"github.com/pgquerynarrative/pgquerynarrative/app/db"
	"github.com/pgquerynarrative/pgquerynarrative/app/embedding"
	"github.com/pgquerynarrative/pgquerynarrative/app/llm"
	"github.com/pgquerynarrative/pgquerynarrative/app/queryrunner"
	"github.com/pgquerynarrative/pgquerynarrative/app/service"
	pkgsuggestions "github.com/pgquerynarrative/pgquerynarrative/app/suggestions"
)

// Client provides access to narrative capabilities: running queries, generating
// reports, schema discovery, saved queries, and the Goa service implementations
// for HTTP or embedded use. All methods accept context.Context; cancellation
// is propagated to the underlying operations. Call Close when done to release
// database connection pools; Close is idempotent and safe to call multiple times.
type Client struct {
	pools          *db.Pools
	queriesService *service.QueriesService
	reportsService *service.ReportsService
	schemaService  *service.SchemaService
	suggester      *pkgsuggestions.Suggester
}

// NewClient builds a narrative client from the given config. It creates
// database pools, query runner, LLM client, and all services. The returned
// client must be closed to release resources.
func NewClient(ctx context.Context, cfg Config) (*Client, error) {
	dbCfg := config.DatabaseConfig{
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
	}
	pools, err := db.NewPools(ctx, dbCfg)
	if err != nil {
		return nil, err
	}

	allowedSchemas := cfg.AllowedSchemas
	if len(allowedSchemas) == 0 {
		allowedSchemas = []string{"demo"}
	}
	maxQueryLength := cfg.MaxQueryLength
	if maxQueryLength <= 0 {
		maxQueryLength = 10000
	}
	maxRows := cfg.MaxRowsPerQuery
	if maxRows <= 0 {
		maxRows = 1000
	}

	validator := queryrunner.NewValidator(allowedSchemas, maxQueryLength)
	runner := queryrunner.NewRunner(pools.ReadOnly, validator, maxRows, cfg.Database.QueryTimeout)
	llmClient := newLLMClient(LLMConfig(cfg.LLM))

	var queriesService *service.QueriesService
	var reportsService *service.ReportsService
	var suggester *pkgsuggestions.Suggester
	embeddingStore := embedding.NewStore(pools.App)

	if cfg.Embedding.BaseURL != "" && cfg.Embedding.Model != "" {
		emb := embedding.NewOllamaEmbedder(cfg.Embedding.BaseURL, cfg.Embedding.Model)
		queriesService = service.NewQueriesServiceWithEmbedding(
			pools.ReadOnly, pools.App, runner, cfg.Metrics.TrendThresholdPercent,
			emb, embeddingStore, cfg.Embedding.Model,
		)
		reportsService = service.NewReportsServiceWithRAG(
			pools.ReadOnly, pools.App, runner, llmClient, cfg.Metrics.TrendThresholdPercent,
			emb, embeddingStore,
		)
		suggester = pkgsuggestions.NewSuggesterWithEmbedding(pools.App, emb, embeddingStore)
	} else {
		queriesService = service.NewQueriesService(
			pools.ReadOnly, pools.App, runner, cfg.Metrics.TrendThresholdPercent,
		)
		reportsService = service.NewReportsService(
			pools.ReadOnly, pools.App, runner, llmClient, cfg.Metrics.TrendThresholdPercent,
		)
		suggester = pkgsuggestions.NewSuggester(pools.App)
	}

	catalogLoader := catalog.NewLoader(pools.ReadOnly, allowedSchemas)
	schemaService := service.NewSchemaService(catalogLoader)

	return &Client{
		pools:          pools,
		queriesService: queriesService,
		reportsService: reportsService,
		schemaService:  schemaService,
		suggester:      suggester,
	}, nil
}

func newLLMClient(cfg LLMConfig) llm.Client {
	switch cfg.Provider {
	case "ollama":
		return llm.NewOllamaClient(cfg.BaseURL, cfg.Model)
	case "gemini":
		return llm.NewGeminiClient(cfg.APIKey, cfg.Model)
	case "claude":
		return llm.NewClaudeClient(cfg.APIKey, cfg.Model)
	case "openai":
		return llm.NewOpenAIClient(cfg.APIKey, cfg.Model)
	case "groq":
		return llm.NewGroqClient(cfg.APIKey, cfg.Model)
	default:
		return llm.NewOllamaClient(cfg.BaseURL, cfg.Model)
	}
}

// Close releases database connection pools. Call once when the client is no longer needed.
func (c *Client) Close() {
	if c.pools != nil {
		c.pools.Close()
	}
}

// Ready returns nil if the database pools are reachable (for readiness probes).
func (c *Client) Ready(ctx context.Context) error {
	if c.pools == nil {
		return errors.New("client not initialized")
	}
	return c.pools.Health(ctx)
}

// QueriesService returns the queries service for use with Goa endpoints or direct calls.
func (c *Client) QueriesService() queries.Service {
	return c.queriesService
}

// ReportsService returns the reports service for use with Goa endpoints or direct calls.
func (c *Client) ReportsService() reports.Service {
	return c.reportsService
}

// SchemaService returns the schema service for use with Goa endpoints or direct calls.
func (c *Client) SchemaService() schema.Service {
	return c.schemaService
}

// SuggestionsService returns the suggestions service for use with Goa endpoints or direct calls.
func (c *Client) SuggestionsService() suggestions.Service {
	return c.suggester
}

// AppPool returns the application database pool for use by the server (e.g. audit logging).
// Do not close the returned pool; Close the Client instead.
func (c *Client) AppPool() *pgxpool.Pool {
	if c == nil || c.pools == nil {
		return nil
	}
	return c.pools.App
}
