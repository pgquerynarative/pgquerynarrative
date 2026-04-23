# UI overview

The web UI is a React SPA (Vite, Tailwind CSS, shadcn/ui) built from [`frontend/`](../frontend/) and served at `/` when running the app. It provides:

- **Query editor** — Schema browser, SQL suggestions, and **Ask** (natural language → SQL + report).
- **Saved queries** — List, create, edit, delete; run and generate reports from a saved query.
- **Reports** — List reports; open a report to view narrative, metrics (time-series, forecast, correlations, cohorts when applicable), and export (HTML/PDF).
- **Connection-aware workflows** — Connection picker in Query Runner/Ask; schema browser refreshes by selected connection; Saved Queries and Reports show `connection_id` badges and support filtering by connection.

Read-only analytics settings (e.g. trend threshold, anomaly sigma) are shown in **Settings → Analytics**; values come from [Configuration](configuration.md#metrics).

**See also:** [Quick start](getting-started/quickstart.md) · [API reference](api/README.md) · [Documentation index](README.md)
