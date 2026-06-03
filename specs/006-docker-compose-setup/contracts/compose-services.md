# Contract: Compose Service Interface

The orchestration contract the implementation MUST satisfy. This is the verifiable boundary for
the feature (no new HTTP API is introduced).

## Stack invariants

- `docker compose up` brings up **two** services: `backend`, `frontend`.
- No host Node or Go toolchain is required; only a container engine + Compose v2.
- All configuration is inline in `docker-compose.yml`; **no `.env` file exists**.

## `backend` contract

| Aspect       | Contract                                                                 |
|--------------|--------------------------------------------------------------------------|
| Host port    | `8080` → container `8080`                                                 |
| Health       | `GET /healthz` → `200 {"status":"ok"}`; used as compose `healthcheck`    |
| API          | `POST /api/v1/calculate` reachable from the host browser                 |
| Env          | `ADDR=:8080`, `CORS_ALLOWED_ORIGINS=http://localhost:5173`               |
| Readiness    | Reports `healthy` once `/healthz` succeeds                               |

**Acceptance**
- `curl localhost:8080/healthz` returns `{"status":"ok"}`.
- `docker compose ps` shows backend `healthy` before frontend is depended-on ready.

## `frontend` contract

| Aspect       | Contract                                                                 |
|--------------|--------------------------------------------------------------------------|
| Host port    | `5173` → container `5173`                                                 |
| Serves       | Calculator UI at `http://localhost:5173`                                  |
| Backend URL  | `VITE_API_BASE_URL=http://localhost:8080` (browser → host-mapped backend) |
| Depends on   | `backend` healthy                                                         |
| Dev loop     | Editing `frontend/src/**` hot-reloads in the running container (US3)      |

**Acceptance**
- Opening `http://localhost:5173` renders the calculator.
- Performing a calculation issues `POST http://localhost:8080/api/v1/calculate` and shows the result
  (no CORS error in console).

## CORS alignment

`CORS_ALLOWED_ORIGINS` (backend) MUST contain the frontend origin `http://localhost:5173`. Changing
either host port REQUIRES updating both the port mapping and this origin.

## Make target contract

`make help` MUST list, with descriptions: `docker-up`, `docker-down`, `docker-build`,
`docker-rebuild`, `docker-logs`. Pre-existing targets (`install build test run-* dev` …) MUST behave
identically to before (FR-013).

## Failure contract

- A host port already in use MUST abort `up` with a clear "port is already allocated" error (FR-011).
- A build failure MUST stop startup and surface the failing stage.
