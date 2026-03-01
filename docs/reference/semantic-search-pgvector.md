# Semantic search (pgvector)

Semantic search (embedding-based similarity) is supported in two ways: **in-memory** (default) and **inside PostgreSQL** with the pgvector extension. It powers [GET /api/v1/suggestions/similar](../api/README.md#suggestions) and RAG context in report generation.

## Flow

1. **Text → vector:** The app sends text to an embedding model (e.g. Ollama `nomic-embed-text`). The model returns a fixed-size vector (e.g. 768 dimensions). No embedding logic inside Postgres.
2. **Storage:** Vectors in `app.query_embeddings`: **JSONB** (always); **vector(768)** (optional when pgvector is enabled via migration 000007).
3. **Similarity search:** With pgvector: `ORDER BY embedding_vector <=> $query_vector LIMIT k` (cosine distance, HNSW index). Without: app loads rows and computes cosine similarity in Go.

## Enabling pgvector

1. **Install the extension:**
   - **Debian/Ubuntu:** `apt install postgresql-16-pgvector` (or your PG version).
   - **Docker:** Use an image that includes pgvector (e.g. `ankane/pgvector`).
   - **Mac (Homebrew):** `./tools/db/install-pgvector-mac.sh` or `brew install pgvector` then restart Postgres. pgvector does **not** require `shared_preload_libraries`; `CREATE EXTENSION vector` in the DB is enough.
2. **Create extension (superuser):** `./tools/db/ensure-pgvector-extension.sh` or `psql -d pgquerynarrative -c 'CREATE EXTENSION IF NOT EXISTS vector;'`.
3. **Run migrations:** `make migrate` (includes 000007 if pgvector is present). If 000007 fails, the app continues with in-memory similarity.
4. **Backfill (optional):** Existing rows have only JSONB; re-save a query or run a backfill so `embedding_vector` is populated. New saves populate both.

## Embeddings config

See [Configuration – Embeddings](../configuration.md#embeddings). Set `EMBEDDING_BASE_URL` (or use default when `LLM_PROVIDER=ollama`) and optionally `EMBEDDING_MODEL`. Used by `GET /api/v1/suggestions/similar` and RAG context in report generation.

## Summary

| Component | Where | Role |
|-----------|--------|------|
| Embedding model | App (Ollama etc.) | Text → vector |
| Vector storage | Postgres | `app.query_embeddings` (jsonb + optional vector column) |
| Similarity search | Postgres (pgvector) or app (in-memory) | k-NN by cosine |
| RAG / narrative | App + LLM | Retrieve similar queries, then generate text |

## See also

- [Configuration](../configuration.md) — Embeddings variables
- [API reference](../api/README.md) — Suggestions endpoints
- [Troubleshooting](troubleshooting.md) — Common issues
- [Documentation index](../README.md)
