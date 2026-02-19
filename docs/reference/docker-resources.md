# Docker resources

Information about Docker storage, memory, and resource usage.

**Storage:** Images (PostgreSQL, app, builder), volumes (e.g. `pgquerynarrative_data`), and build cache.  
**Memory:** PostgreSQL typically 200 MB–1 GB; app 50–200 MB. Limits are set in `docker-compose.yml` (PostgreSQL 512 MB, app 256 MB). The default image is `postgres:18-alpine` for a smaller footprint.

**Freeing space:**  
- `docker image prune -a` — remove unused images  
- `docker builder prune -a` — clear build cache  
- `docker system prune -a` — remove unused data (use with care)  
- `docker volume prune` — remove unused volumes (includes database data)

**Docker Desktop:** Use Settings → Resources to cap memory, disk, and CPU.
