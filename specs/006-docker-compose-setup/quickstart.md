# Quickstart: Containerized Stack

Run the full wonderlic-calc stack (frontend + backend) in containers — no Node or Go on your host.

## Prerequisites

- Docker Engine + Docker Compose v2 (`docker compose version` works).
- Nothing else. No Node, no Go.

## Start everything

```bash
make docker-up
```

- Builds both images on first run, then starts both services detached.
- Frontend: http://localhost:5173
- Backend health: http://localhost:8080/healthz → `{"status":"ok"}`

## Smoke test (acceptance)

```bash
curl -s localhost:8080/healthz            # {"status":"ok"}
```

Then open http://localhost:5173, perform a calculation, and confirm the result returns (no CORS
error in the browser console). This exercises the browser → `localhost:8080` → backend path.

## Everyday commands

| Command              | What it does                                  |
|----------------------|-----------------------------------------------|
| `make docker-up`     | Start the stack (build if needed)             |
| `make docker-logs`   | Follow combined logs                          |
| `make docker-down`   | Stop and remove containers                    |
| `make docker-build`  | Rebuild images without starting               |
| `make docker-rebuild`| Rebuild images and restart (apply code changes)|

`make help` lists these alongside the existing targets.

## Iterating on code (US3)

- **Frontend**: edit `frontend/src/**` — Vite HMR updates the browser live (source is bind-mounted).
- **Backend**: edit `backend/**`, then `make docker-rebuild` to rebuild + restart.

## Configuration

All config is inline in `docker-compose.yml` — there is **no `.env` file**:

| Service  | Variable               | Default                  |
|----------|------------------------|--------------------------|
| backend  | `ADDR`                 | `:8080`                  |
| backend  | `CORS_ALLOWED_ORIGINS` | `http://localhost:5173`  |
| frontend | `VITE_API_BASE_URL`    | `http://localhost:8080`  |

To use different host ports, edit the `ports` mapping **and** `CORS_ALLOWED_ORIGINS` /
`VITE_API_BASE_URL` together (they must stay aligned).

## Troubleshooting

- **`port is already allocated`**: another process holds 5173 or 8080. Stop it, or change the host
  port (and the aligned env values), then `make docker-up`.
- **CORS error in console**: confirm `CORS_ALLOWED_ORIGINS` includes the exact frontend origin.
- **Frontend can't reach backend**: confirm `make docker-logs` shows backend `healthy`; the frontend
  waits for backend health before it is considered ready.

## Non-container workflow still works

`make dev`, `make run-frontend`, `make run-backend`, `make test` are unchanged for contributors with
local toolchains.
