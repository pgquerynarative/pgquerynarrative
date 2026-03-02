package testhelpers

import (
	"context"
	"os"
	"testing"

	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func RunPostgresContainer(t *testing.T, ctx context.Context) *postgres.PostgresContainer {
	t.Helper()

	if os.Getenv("DOCKER_API_VERSION") == "" {
		// Default to 1.44 so newer Docker daemons accept the client (daemon may require minimum 1.44).
		_ = os.Setenv("DOCKER_API_VERSION", "1.44")
	}

	defer func() {
		if r := recover(); r != nil {
			t.Skipf("docker not available or API mismatch: %v", r)
		}
	}()

	// Use pgvector image so migration 000007_pgvector_embeddings can run.
	container, err := postgres.Run(ctx, "pgvector/pgvector:pg18",
		postgres.WithDatabase("pgquerynarrative"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
	)
	if err != nil {
		t.Skipf("docker not available: %v", err)
	}
	return container
}
