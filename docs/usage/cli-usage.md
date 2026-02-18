# CLI usage

Command-line access to the API. The app must be running (`make start-docker` or `make start-local`). The CLI runs inside a container when using Docker, or on the host when using a local binary.

## Commands

| Command | Description |
|---------|-------------|
| `make cli CMD='query "SQL"'` | Run query (optional limit: `query "SQL" 10`) |
| `make cli CMD='list'` | List saved queries |
| `make cli CMD='get "uuid"'` | Get saved query by ID |
| `make cli CMD='save "Name" "SQL"'` | Save query (optional: `"tags,a,b"`) |
| `make cli CMD='report "SQL"'` | Generate report (requires LLM) |

**Interactive:** `make cli-shell` then e.g. `pgquerynarrative query "SELECT * FROM demo.sales LIMIT 5"` or alias `pqn list`.

## Examples

```bash
make cli CMD='query "SELECT product_category, SUM(total_amount) FROM demo.sales GROUP BY product_category"'
make cli CMD='save "Top Products" "SELECT product_category, SUM(total_amount) FROM demo.sales GROUP BY product_category"'
make cli CMD='report "SELECT product_category, SUM(total_amount) FROM demo.sales GROUP BY product_category"'
```

## Environment (CLI container)

| Variable | Default | Description |
|----------|---------|-------------|
| `PGQUERYNARRATIVE_API_URL` | `http://app:8080` | API base URL (use `http://localhost:8080` when CLI runs on host) |
| `PGQUERYNARRATIVE_FORMAT` | `table` | Output: `table` or `json` |

JSON output: `docker compose run --rm -e PGQUERYNARRATIVE_FORMAT=json cli pgquerynarrative query "SELECT 1"`.

**Quoting:** Quote SQL in the outer command: `make cli CMD='query "SELECT * FROM demo.sales"'`. For SQL with single quotes escape: `'\''`.

**See also:** [API examples](../api/examples.md), [API reference](../api/README.md), [Troubleshooting](../troubleshooting.md)
