# Embedded integration

Use PgQueryNarrative inside your own Go service: create a `narrative.Client` and call it directly (library) or mount its HTTP endpoints with the provided middleware (Chi, Gin, Echo). Config matches the [standalone server](installation.md); see [Configuration](../configuration.md).

## Library usage

```go
import (
    "github.com/pgquerynarrative/pgquerynarrative/pkg/narrative"
    "github.com/pgquerynarrative/pgquerynarrative/app/config"
)

cfg := narrative.FromAppConfig(config.Load())
client, err := narrative.NewClient(ctx, cfg)
if err != nil { ... }
defer client.Close()

result, err := client.RunQuery(ctx, "SELECT ... FROM demo.sales", 100)
report, err := client.GenerateReport(ctx, sql)
schema, err := client.GetSchema(ctx)
```

See `examples/library-usage/basic.go`.

## HTTP middleware (Chi, Gin, Echo)

Mount narrative endpoints on your router. Package: [pkg/narrative/middleware](https://github.com/pgquerynarrative/pgquerynarrative/tree/main/pkg/narrative/middleware).

| Framework | Mount call |
|-----------|------------|
| **Chi** | `narrativemw.MountChi(r, client, "/api")` |
| **Gin** | `narrativemw.MountGin(r, client, "/api")` |
| **Echo** | `narrativemw.MountEcho(e, client, "/api")` |

Mounted routes (with prefix `/api`):

| Method | Path | Description |
|--------|------|-------------|
| POST | /api/query/run | Body: `{"sql":"...", "limit": N}`. Run read-only SQL. |
| POST | /api/report/generate | Body: `{"sql":"..."}`. Generate narrative report (requires [LLM](llm-setup.md)). |
| GET | /api/schema | Allowed schemas, tables, columns. |
| GET | /api/suggestions/queries | Query: `intent`, `limit`. Suggested SQL. |

Use empty prefix `""` to mount at root (e.g. `/query/run`).

## Examples

| Example | Description |
|---------|-------------|
| `examples/library-usage/basic.go` | Client only |
| `examples/chi-integration/main.go` | Chi + middleware |
| `examples/gin-integration/main.go` | Gin + middleware |
| `examples/echo-integration/main.go` | Echo + middleware |

Build and run (set `DATABASE_*` and `LLM_*` as needed):

```bash
go build -o bin/example-chi ./examples/chi-integration
./bin/example-chi
```

## See also

- [Configuration](../configuration.md) · [API reference](../api/README.md) · [Documentation index](../README.md)
