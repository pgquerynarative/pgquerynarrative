package story

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/pgquerynarrative/pgquerynarrative/app/llm"
	"github.com/pgquerynarrative/pgquerynarrative/app/metrics"
)

// Generator creates narratives from query results
type Generator struct {
	llmClient llm.Client
}

// NewGenerator creates a new narrative generator
func NewGenerator(llmClient llm.Client) *Generator {
	return &Generator{
		llmClient: llmClient,
	}
}

// Generate creates a narrative from query results and metrics
func (g *Generator) Generate(ctx context.Context, sql string, columns []string, rows [][]interface{}, calcMetrics *metrics.Metrics) (*NarrativeContent, error) {
	// Convert metrics to JSON (compact format uses less memory than MarshalIndent)
	metricsJSON, err := json.Marshal(calcMetrics)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal metrics: %w", err)
	}

	// Build prompt
	prompt := llm.BuildNarrativePrompt(sql, columns, rows, string(metricsJSON))

	// Generate narrative using LLM
	response, err := g.llmClient.Generate(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to generate narrative: %w", err)
	}

	// Parse response
	narrative, err := ParseNarrative(response)
	if err != nil {
		return nil, fmt.Errorf("failed to parse narrative: %w", err)
	}

	return narrative, nil
}
