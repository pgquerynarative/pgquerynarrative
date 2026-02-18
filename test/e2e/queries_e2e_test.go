package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	goahttp "goa.design/goa/v3/http"

	queriesServer "github.com/pgquerynarrative/pgquerynarrative/api/gen/http/queries/server"
	"github.com/pgquerynarrative/pgquerynarrative/api/gen/queries"
	"github.com/pgquerynarrative/pgquerynarrative/app/queryrunner"
	"github.com/pgquerynarrative/pgquerynarrative/app/service"
	"github.com/pgquerynarrative/pgquerynarrative/test/testhelpers"
)

func TestQueriesE2E(t *testing.T) {
	ctx := context.Background()
	container := testhelpers.RunPostgresContainer(t, ctx)
	t.Cleanup(func() { _ = container.Terminate(ctx) })

	connStr, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("failed to get connection string: %v", err)
	}

	migrationsPath, err := filepath.Abs("../../app/db/migrations")
	if err != nil {
		t.Fatalf("failed to resolve migrations path: %v", err)
	}

	m, err := migrate.New("file://"+migrationsPath, connStr)
	if err != nil {
		t.Fatalf("failed to create migrator: %v", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		t.Fatalf("failed to run migrations: %v", err)
	}

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		t.Fatalf("failed to create pool: %v", err)
	}
	defer pool.Close()

	_, err = pool.Exec(ctx, `
		INSERT INTO demo.sales (id, date, product_category, product_name, quantity, unit_price, total_amount, region, sales_rep)
		VALUES (gen_random_uuid(), CURRENT_DATE, 'Electronics', 'Alpha', 5, 10.00, 50.00, 'North', 'A. Lee')
	`)
	if err != nil {
		t.Fatalf("failed to seed data: %v", err)
	}

	validator := queryrunner.NewValidator([]string{"demo"}, 10000)
	runner := queryrunner.NewRunner(pool, validator, 1000, 30*time.Second)
	queriesService := service.NewQueriesService(pool, pool, runner, 0)

	endpoints := queries.NewEndpoints(queriesService)

	mux := goahttp.NewMuxer()
	dec := goahttp.RequestDecoder
	enc := goahttp.ResponseEncoder
	errHandler := func(ctx context.Context, w http.ResponseWriter, err error) {
		_ = goahttp.ErrorEncoder(enc, nil)(ctx, w, err)
	}

	httpServer := queriesServer.New(endpoints, mux, dec, enc, errHandler, nil)
	queriesServer.Mount(mux, httpServer)

	testServer := httptest.NewServer(mux)
	t.Cleanup(testServer.Close)

	runPayload := map[string]interface{}{
		"sql":   "SELECT product_category, SUM(total_amount) AS total FROM demo.sales GROUP BY product_category",
		"limit": 100,
	}
	payloadBytes, _ := json.Marshal(runPayload)
	resp, err := http.Post(testServer.URL+"/api/v1/queries/run", "application/json", bytes.NewReader(payloadBytes))
	if err != nil {
		t.Fatalf("run query request failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status: %d", resp.StatusCode)
	}

	savePayload := map[string]interface{}{
		"name": "Sample Query",
		"sql":  "SELECT * FROM demo.sales",
		"tags": []string{"demo"},
	}
	saveBytes, _ := json.Marshal(savePayload)
	resp, err = http.Post(testServer.URL+"/api/v1/queries/saved", "application/json", bytes.NewReader(saveBytes))
	if err != nil {
		t.Fatalf("save query request failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected save status: %d", resp.StatusCode)
	}
	var saved queries.SavedQuery
	if err := json.NewDecoder(resp.Body).Decode(&saved); err != nil {
		t.Fatalf("failed to decode saved query: %v", err)
	}
	if saved.ID == "" {
		t.Fatalf("expected saved query ID")
	}
}
