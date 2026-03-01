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
		// Default to a widely supported Docker API version.
		_ = os.Setenv("DOCKER_API_VERSION", "1.41")
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
