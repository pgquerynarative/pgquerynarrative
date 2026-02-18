# Docker resources

**Storage:** Images (postgres, app, builder), volumes (e.g. `pgquerynarrative_data`), build cache. **Memory:** Postgres ~200 MB–1 GB; app ~50–200 MB. Limits in `docker-compose.yml`: Postgres 512 MB, app 256 MB. Image: `postgres:18-alpine` (smaller than default).

**Free space:** `docker image prune -a`, `docker builder prune -a`, `docker system prune -a`. `docker volume prune` removes volumes (including DB data).

**Docker Desktop:** Settings → Resources to cap memory, disk, CPU.
