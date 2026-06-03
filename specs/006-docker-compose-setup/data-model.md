# Phase 1 Data Model: Docker Compose Setup

This feature has no runtime/persistence data model. The "entities" are configuration objects:
the compose services and their inline environment. Documented here as the structural contract.

## Entity: Compose stack (`docker-compose.yml`)

| Field        | Value / rule                                                            |
|--------------|-------------------------------------------------------------------------|
| services     | exactly two: `backend`, `frontend`                                      |
| networks     | default bridge (implicit); no custom network required (see research R1) |
| volumes      | bind mounts for source (HMR); named volume for frontend `node_modules`  |

## Entity: `backend` service

| Field         | Value / rule                                                           |
|---------------|------------------------------------------------------------------------|
| build.context | `./backend`                                                            |
| build.dockerfile | `Dockerfile` (multi-stage, `golang:1.25-alpine` → slim runtime)     |
| ports         | `8080:8080` (host:container)                                           |
| environment   | `ADDR=:8080`, `CORS_ALLOWED_ORIGINS=http://localhost:5173` (inline)    |
| healthcheck   | test `wget -qO- http://localhost:8080/healthz` ; interval/retries set  |
| restart       | `unless-stopped`                                                       |

**Validation rules**
- `CORS_ALLOWED_ORIGINS` MUST include the frontend host origin or browser calls fail (FR-005).
- Container MUST expose/listen on `8080` to match the mapping and healthcheck.

## Entity: `frontend` service

| Field         | Value / rule                                                            |
|---------------|-------------------------------------------------------------------------|
| build.context | `./frontend`                                                            |
| build.dockerfile | `Dockerfile` (`node:24-alpine`, runs Vite dev server)               |
| command/CMD   | Vite dev server bound to `0.0.0.0:5173`                                  |
| ports         | `5173:5173`                                                             |
| environment   | `VITE_API_BASE_URL=http://localhost:8080` (inline)                      |
| depends_on    | `backend: { condition: service_healthy }`                               |
| volumes       | `./frontend:/app` (source) + `/app/node_modules` (anonymous, preserve)  |
| restart       | `unless-stopped`                                                        |

**Validation rules**
- Vite `server.host` MUST be truthy (`vite.config.ts`) so the dev server binds `0.0.0.0` and is
  reachable from the host; otherwise the mapped port serves nothing.
- `VITE_API_BASE_URL` MUST point at the host-mapped backend origin (`localhost:8080`), per research R1.

## Entity: Inline environment (no `.env`)

| Variable               | Service  | Default applied inline      | Overridable via |
|------------------------|----------|-----------------------------|-----------------|
| `ADDR`                 | backend  | `:8080`                     | edit compose / `docker compose` shell env |
| `CORS_ALLOWED_ORIGINS` | backend  | `http://localhost:5173`     | edit compose    |
| `VITE_API_BASE_URL`    | frontend | `http://localhost:8080`     | edit compose    |

No `.env` file is created or committed (FR-014). Values live in the committed `docker-compose.yml`.

## Entity: Make container targets (`Makefile`)

| Target           | Action                                                       |
|------------------|--------------------------------------------------------------|
| `docker-up`      | build if needed + start stack (detached)                     |
| `docker-down`    | stop + remove containers (and the anonymous volumes)         |
| `docker-build`   | build/rebuild images without starting                        |
| `docker-rebuild` | rebuild images + restart stack (backend refresh path, US3)   |
| `docker-logs`    | follow combined logs of both services                        |

All declared with `## <description>` so `make help` auto-lists them (FR-008).
