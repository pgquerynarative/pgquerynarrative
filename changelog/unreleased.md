## [Unreleased]

### Added

- **Migration 000011:** Set `default_transaction_read_only = on` for the readonly role so the database enforces read-only at session level.
- **Configurable windows:** `METRICS_MAX_TIMESERIES_PERIODS` (default 24, range 2–120) controls the maximum number of periods included in time-series output ("last N" for charts and period history). Configuration doc now has a "Configurable windows" subsection listing all window-related vars (trend periods, moving avg, max time-series periods, seasonal lag, min periods for seasonality). Settings UI shows max time-series periods when present.

### Planned (Release 2)

Additional analytics: further cohort metrics, configurable windows, and seasonal adjustments.
