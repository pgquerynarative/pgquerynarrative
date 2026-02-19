# Installation

## Prerequisites

- **Go 1.24+** (build from source)
- **PostgreSQL 16+** or Docker

Optional: [LLM setup](llm-setup.md) for report generation.

## Docker (recommended)

```bash
git clone https://github.com/your-org/pgquerynarrative.git
cd pgquerynarrative
make start-docker
```

Starts PostgreSQL, runs migrations and seed, then the app.

## Local (from source)

**1. Install Go and PostgreSQL** (macOS: `brew install go postgresql@18`; Linux: use distro packages).

**2. Clone and build**

```bash
git clone https://github.com/your-org/pgquerynarrative.git
cd pgquerynarrative
make setup
make generate
make build
```

**3. Database**

```bash
make db-init
make migrate
make seed
```

**4. Run**

```bash
make run
# or
./bin/server
```

## Verify

```bash
curl http://localhost:8080/api/v1/queries/saved
```

## PostgreSQL version

Supported: 16, 17, 18. Docker default: `postgres:18-alpine`. Override: `POSTGRES_IMAGE=postgres:17-alpine make start-docker`.

**See also:** [Configuration](../configuration.md), [Quick start](quickstart.md), [Troubleshooting](../reference/troubleshooting.md)
