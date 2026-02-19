# Unit tests

All feature and unit tests live under this directory. They can be run as a whole or by package.

## Layout

- `app/charts/` – chart suggestion logic (bar, line, pie, area, table)
- `app/metrics/` – metrics calculator (time series, aggregates, data quality, std dev, period labels)
- `app/queryrunner/` – query validation (schema, SELECT-only)
- `app/service/` – reports service (perf suggestions, convertMetrics to API types)
- `app/story/` – narrative sanitizer (remove fabricated "previous period")
- `web/` – report HTML formatting (charts, data quality, perf suggestions, narrative)

Server middleware tests remain in `cmd/server/` (they test unexported middleware).

## Run all unit tests

```bash
make test-unit
```

## Run one package

```bash
go test ./test/unit/app/metrics/... -v
go test ./test/unit/web/... -v
```

## Run one test

```bash
go test ./test/unit/app/service/... -run TestBuildPerfSuggestions_LimitApplied -v
```

See `docs/qa-features.md` for the full feature → test mapping.
