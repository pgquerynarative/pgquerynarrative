// Package service provides business logic for queries, reports, and schema.
package service

import (
	"context"

	schema "github.com/pgquerynarrative/pgquerynarrative/gen/schema"

	"github.com/pgquerynarrative/pgquerynarrative/app/catalog"
)

// SchemaService implements the schema service API (catalog discovery).
type SchemaService struct {
	loader *catalog.Loader
}

// NewSchemaService creates a schema service that returns allowed schemas,
// tables, and columns from the database (read-only pool).
func NewSchemaService(loader *catalog.Loader) *SchemaService {
	return &SchemaService{loader: loader}
}

// Get returns the list of allowed schemas with their tables and columns.
func (s *SchemaService) Get(ctx context.Context) (*schema.SchemaResult, error) {
	return s.loader.Load(ctx)
}
