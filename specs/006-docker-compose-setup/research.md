# Phase 0 Research: Docker Compose Setup

## R1. Browser→backend networking topology

**Decision**: The browser runs on the host and calls the backend via the **host-mapped port**
(`http://localhost:8080`), not via a compose service name. Compose service-name DNS
(`backend:8080`) is irrelevant for browser traffic because the browser is outside the docker
network.

**Rationale**: The frontend is a SPA executing in the user's browser; `fetch` originates from the
host, not from the frontend container. Only the host port mapping matters.

**Implication**: `VITE_API_BASE_URL=http://localhost:8080`; backend `CORS_ALLOWED_ORIGINS` must
include `http://localhost:5173` (the frontend's host origin). Both already match existing defaults.

**Alternatives considered**: A reverse proxy fronting both services on one origin (eliminates CORS)
— rejected as over-engineering for a 2-service local dev stack and out of spec scope.

## R2. Frontend container: dev server vs. static build

**Decision**: Run the **Vite dev server** in the frontend container (`vite --host 0.0.0.0 --port 5173`)
with the workspace bind-mounted.

**Rationale**: Directly satisfies US3 (hot reload while containerized) with zero extra tooling; keeps
the container definition simple; aligns with the local-dev scope (not production serving). Vite reads
`VITE_API_BASE_URL` from the process environment, so the inline compose `environment` value flows into
`import.meta.env` at dev-server start — no `.env` file needed (FR-014).

**Alternatives considered**: Multi-stage build → static assets served by nginx. Rejected for v1: more
moving parts, requires build-arg plumbing for the build-time `VITE_API_BASE_URL`, and loses HMR. Noted
as a future production-image follow-up (out of scope).

## R3. Backend container build

**Decision**: Multi-stage `Dockerfile` — `golang:1.25-alpine` builds `./cmd/server`, final stage runs
the compiled static binary on a minimal base. Inline `ADDR=:8080` and `CORS_ALLOWED_ORIGINS=http://localhost:5173`
via compose `environment`.

**Rationale**: Small runtime image, fast restarts, no host Go (FR-002, FR-009). Config already env-driven
(`config.Load`), so no code change.

**Backend refresh path (US3)**: `make docker-rebuild` (rebuild image + restart). A file-watch hot-reloader
(e.g. `air`) is deliberately **not** added — it is a new dependency requiring justification (Principle III)
and the rebuild path is fast enough for this stack. Frontend retains live HMR via R2.

**Alternatives considered**: `go run` with source mount for backend hot reload — rejected: pulls the full
Go toolchain into the runtime image (defeats the slim runtime stage) for marginal gain.

## R4. Inline environment, no `.env` files

**Decision**: Declare all configuration with compose `environment:` lists (and `Dockerfile ENV`/`EXPOSE`
where it documents the image). No `.env` files are created or committed.

**Rationale**: Explicit user directive. Values are non-secret (listen address, CORS origins, backend URL).
Inline values are visible in the committed compose file → zero-config, self-documenting startup (FR-015).
Sidesteps the repo `.gitignore` `.env*` rule entirely.

**Alternatives considered**: Per-workspace committed `.env` files (earlier draft) — explicitly rejected by user.

## R5. Healthcheck & startup ordering

**Decision**: Backend service defines a `healthcheck` hitting `/healthz`; frontend uses
`depends_on: { backend: { condition: service_healthy } }`.

**Rationale**: Satisfies FR-006 (frontend considered ready only once backend reachable) and gives a clean
failure signal. `/healthz` already returns `{"status":"ok"}`.

**Alternatives considered**: No ordering (frontend handles backend-down gracefully anyway) — kept the
graceful-degradation behavior but added `depends_on` for a smoother first-run experience.

## R6. Port collision / fail-fast

**Decision**: Map fixed host ports `5173:5173` and `8080:8080`. Rely on Docker's bind-failure error,
which aborts `up` with a clear "port is already allocated" message (FR-011). Document overriding host
ports in quickstart (FR-010).

**Rationale**: Native compose behavior already fails fast and clearly; no custom logic needed.

## R7. Make target naming

**Decision**: `docker-up`, `docker-down`, `docker-build`, `docker-rebuild`, `docker-logs`, all `## `-
documented so `make help`'s existing awk auto-lists them (FR-007, FR-008). Existing targets untouched (FR-013).

**Rationale**: Matches the repo's hyphenated, self-documenting target convention (`run-frontend`, `test-backend`).

## Environment matrix (resolved)

| Service  | Variable               | Value (inline)                | Source of truth        |
|----------|------------------------|-------------------------------|------------------------|
| backend  | `ADDR`                 | `:8080`                       | `config.Load` default  |
| backend  | `CORS_ALLOWED_ORIGINS` | `http://localhost:5173`       | `config.Load` default  |
| frontend | `VITE_API_BASE_URL`    | `http://localhost:8080`       | `frontend/src/lib/api.ts` default |

All NEEDS CLARIFICATION resolved.
