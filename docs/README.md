# Documentation

Documentation for PgQueryNarrative: running queries, generating reports, configuration, and development.

## Getting started

| Document | Description |
|----------|-------------|
| [Installation](getting-started/installation.md) | Prerequisites, database setup, running the app |
| [Quick start](getting-started/quickstart.md) | Minimal steps to run with Docker or local PostgreSQL |
| [LLM setup](getting-started/llm-setup.md) | Configure LLM for report generation and MCP |

## Usage

| Document | Description |
|----------|-------------|
| [Configuration](configuration.md) | Environment variables (server, database, LLM, metrics) |
| [API reference](api/README.md) | REST endpoints, payloads, errors |
| [API examples](api/examples.md) | cURL examples for run, save, report |
| [CLI usage](usage/cli-usage.md) | Command-line interface (`make cli`, `cli-shell`) |

## Features

| Document | Description |
|----------|-------------|
| [Period comparison](features/period-comparison.md) | Automatic period-over-period (vs previous period), % change, trend, configuration |

## Reference

| Document | Description |
|----------|-------------|
| [Troubleshooting](troubleshooting.md) | Common issues and fixes |
| [Docker resources](docker-resources.md) | Containers and Compose |
| [PostgreSQL extension](postgres-extension.md) | Running queries and reports from SQL |

## Development

| Document | Description |
|----------|-------------|
| [Development setup](development/setup.md) | Build, test, code generation |
| [Contributing](../.github/CONTRIBUTING.md) | How to contribute |
| [Security](../.github/SECURITY.md) | Security policy |

## Changelog

[Changelog](../changelog/README.md) — release history and unreleased changes.
