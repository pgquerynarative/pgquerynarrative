## [Unreleased]

### Added

- **Logging (zerolog):** `LOG_LEVEL` (debug, info, warn, error) and `LOG_PRETTY` (colorful console output for local dev). When set, app uses zerolog for leveled, structured logs; request logging uses one message per request (`http request`) with level by status (4xx/5xx → error, /health, /ready, /version → debug, else info). Documented in Configuration – Logging.
- **Migration 000011:** Set `default_transaction_read_only = on` for the readonly role so the database enforces read-only at session level.
- **Configurable windows:** `METRICS_MAX_TIMESERIES_PERIODS` (default 24, range 2–120) controls the maximum number of periods included in time-series output ("last N" for charts and period history). Configuration doc now has a "Configurable windows" subsection listing all window-related vars (trend periods, moving avg, max time-series periods, seasonal lag, min periods for seasonality). Settings UI shows max time-series periods when present.

### Planned (Release 2)

Additional analytics: further cohort metrics and seasonal adjustments.
