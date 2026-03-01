# Deployment

Build and deploy PgQueryNarrative with Docker, Kubernetes, or Helm. For first-time run see [Quick start](../getting-started/quickstart.md) or [Installation](../getting-started/installation.md).

**Development vs production:** `make start-docker` uses root [docker-compose.yml](../../docker-compose.yml) and [Dockerfile](../../Dockerfile) (includes React frontend build). Below: **production** image (`deploy/docker/Dockerfile`) and Compose/Kubernetes/Helm.

---

## Docker

### Build

From repo root:

```bash
docker build -f deploy/docker/Dockerfile -t pgquerynarrative:latest .
```

Multi-stage build: Go app + migrate binary, then minimal Alpine image (server binary, migrations, entrypoint). Production image does not bundle the React SPA; for full UI use the root Dockerfile or serve static assets separately.

### Run with Docker Compose

```bash
docker compose -f deploy/docker/docker-compose.yml up -d
```

Or build and run: `docker compose -f deploy/docker/docker-compose.yml up -d --build`.

Set env or use `.env`. Required: `DATABASE_*`, `LLM_*` (see [Configuration](../configuration.md)). Optional: `PGQUERYNARRATIVE_SEED=true` for demo seed. App waits for Postgres, runs migrations, then starts. API and [health endpoints](operations.md#health-checks): http://localhost:8080.

### Pre-built image

Point Compose at your registry image:

```yaml
services:
  app:
    image: your-registry/pgquerynarrative:1.0.0
```

---

## Kubernetes

Manifests: `deploy/kubernetes/`. PostgreSQL is external; app connects via `DATABASE_HOST` and credentials from a Secret.

### Prerequisites

- Cluster and `kubectl` configured.
- PostgreSQL reachable from the cluster. Create DB and roles, run migrations once if DB is empty (see [Installation](../getting-started/installation.md)).

### Apply order

1. Namespace (optional).
2. Secret (edit with real credentials; do not commit).
3. ConfigMap (DB host, LLM).
4. Deployment, Service.
5. Ingress (optional).

```bash
kubectl apply -f deploy/kubernetes/namespace.yaml
kubectl apply -f deploy/kubernetes/secret.yaml
kubectl apply -f deploy/kubernetes/configmap.yaml
kubectl apply -f deploy/kubernetes/deployment.yaml
kubectl apply -f deploy/kubernetes/service.yaml
kubectl apply -f deploy/kubernetes/ingress.yaml   # optional
```

Set `image` in `deployment.yaml` to your image. Ensure ConfigMap has correct `DATABASE_HOST`. Probes: [Operations – Health checks](operations.md#health-checks) (GET /health, GET /ready).

### Access

- No Ingress: `kubectl port-forward -n pgquerynarrative svc/pgquerynarrative 8080:8080` → http://localhost:8080.
- With Ingress: configure controller and DNS for host in `ingress.yaml`.

---

## Helm

Chart: `deploy/helm/pgquerynarrative/`. Deploys app with ConfigMap, Secret, Deployment, Service, optional Ingress.

### Install

```bash
helm install pgqn ./deploy/helm/pgquerynarrative -n pgquerynarrative --create-namespace
```

Override: `--set image.repository=... --set image.tag=1.0.0 --set database.host=... --set secret.databasePassword=xxx` or `-f my-values.yaml`.

### Upgrade / uninstall

```bash
helm upgrade pgqn ./deploy/helm/pgquerynarrative -n pgquerynarrative
helm uninstall pgqn -n pgquerynarrative
```

### Chart values

See `deploy/helm/pgquerynarrative/values.yaml`. Key keys: **image**, **database**, **secret**, **llm**, **ingress.enabled** / **ingress.host**, **seed**.

---

## Summary

| Method | Path | Use case |
|--------|------|----------|
| Docker Compose | `deploy/docker/` | Single host, dev or staging |
| Kubernetes | `deploy/kubernetes/` | Raw manifests |
| Helm | `deploy/helm/pgquerynarrative/` | Parameterized K8s install |

## See also

- [Configuration](../configuration.md) · [Operations](operations.md) · [Installation](../getting-started/installation.md) · [Documentation index](../README.md)
