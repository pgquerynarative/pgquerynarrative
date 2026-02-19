# PgQueryNarrative documentation

PgQueryNarrative turns SQL query results into business narratives with AI. This documentation covers installation, configuration, the API, and development.

---

## Getting started

| Document | Description |
|----------|-------------|
| [Installation](getting-started/installation.md) | Prerequisites, database setup, and running the application |
| [Quick start](getting-started/quickstart.md) | Minimal steps to run with Docker or local PostgreSQL |
| [LLM setup](getting-started/llm-setup.md) | Configure an LLM provider for report generation (Ollama, OpenAI, Claude, Gemini, Groq) and MCP |

---

## User guides

| Document | Description |
|----------|-------------|
| [Configuration](configuration.md) | Environment variables (server, database, LLM, metrics) |
| [CLI usage](usage/cli-usage.md) | Command-line interface for queries, saved queries, and reports |

---

## API

| Document | Description |
|----------|-------------|
| [API reference](api/README.md) | REST endpoints, request/response formats, error codes |
| [API examples](api/examples.md) | cURL examples for running queries, saving queries, and generating reports |

---

## Features

| Document | Description |
|----------|-------------|
| [Period comparison](features/period-comparison.md) | Automatic period-over-period comparison, trend, and configuration |

---

## Reference

| Document | Description |
|----------|-------------|
| [Troubleshooting](reference/troubleshooting.md) | Common issues and solutions |
| [Docker resources](reference/docker-resources.md) | Containers, storage, and resource usage |
| [PostgreSQL extension](reference/postgres-extension.md) | Run queries and generate reports from SQL |

---

## Development

| Document | Description |
|----------|-------------|
| [Development setup](development/setup.md) | Build, test, code generation, and workflow |
| [Testing](development/testing.md) | Running unit and integration tests, QA checklist |

For contributing and security, see the [.github](https://github.com/pgquerynarrative/pgquerynarrative/tree/main/.github) directory (CONTRIBUTING.md, SECURITY.md).

---

## Changelog

Release history and unreleased changes: [changelog](../changelog/README.md).
