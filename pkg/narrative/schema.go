package narrative

import (
	"context"

	schema "github.com/pgquerynarrative/pgquerynarrative/api/gen/schema"
)

// GetSchema returns the list of allowed schemas with their tables and columns
// (from information_schema). Use this for discovery when building queries
// programmatically. Context cancellation is propagated.
func (c *Client) GetSchema(ctx context.Context) (*schema.SchemaResult, error) {
	return c.schemaService.Get(ctx, &schema.GetPayload{})
}
