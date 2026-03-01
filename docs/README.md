# PgQueryNarrative documentation

PgQueryNarrative turns SQL query results into business narratives using an LLM. Run read-only SQL against PostgreSQL, compute metrics and chart suggestions, and generate narrative reports. The React SPA (built from `frontend/`) provides the web UI; the [REST API](api/README.md) and [CLI](usage/cli-usage.md) support automation.

**Recommended path:** [Quick start](getting-started/quickstart.md) → [LLM setup](getting-started/llm-setup.md) (for reports) → [Configuration](configuration.md).

---

## Documentation index

### Getting started

| Document | Description |
|----------|-------------|
| [Installation](getting-started/installation.md) | Prerequisites, database setup, Docker and local run methods, verification |
| [Quick start](getting-started/quickstart.md) | Minimal steps to run with Docker or local PostgreSQL |
| [LLM setup](getting-started/llm-setup.md) | Configure LLM for report generation (Ollama, OpenAI, Claude, Gemini, Groq) and MCP |
| [Embedded integration](getting-started/embedded.md) | Use as a Go library or mount HTTP endpoints in Chi, Gin, or Echo |

### User guides

| Document | Description |
|----------|-------------|
| [Configuration](configuration.md) | Environment variables (server, database, LLM, embeddings, MCP, metrics) |
| [CLI usage](usage/cli-usage.md) | Command-line interface for running queries, saved queries, and reports |

### API

| Document | Description |
|----------|-------------|
| [API reference](api/README.md) | REST endpoints, request/response shapes, error codes |
| [API examples](api/examples.md) | cURL examples for run, save, and reports |

### Reference

| Document | Description |
|----------|-------------|
| [Deployment](reference/deployment.md) | Docker build, Docker Compose, Kubernetes, Helm |
| [Operations](reference/operations.md) | Monitoring, health checks (`/health`, `/ready`), runbooks |
| [Troubleshooting](reference/troubleshooting.md) | Common issues and fixes |
| [PostgreSQL extension](reference/postgres-extension.md) | Call the API from SQL via `CREATE EXTENSION pgquerynarrative` |
| [Semantic search (pgvector)](reference/semantic-search-pgvector.md) | Embeddings, similar-query search, RAG in report generation |

### Development

| Document | Description |
|----------|-------------|
| [Development setup](development/setup.md) | Build, test, codegen, workflow, frontend build |
| [Testing](development/testing.md) | Unit, integration, and E2E tests |

---

**Contributing & security:** [.github/CONTRIBUTING.md](../.github/CONTRIBUTING.md) · [.github/SECURITY.md](../.github/SECURITY.md).

**Changelog:** [CHANGELOG.md](../CHANGELOG.md).
