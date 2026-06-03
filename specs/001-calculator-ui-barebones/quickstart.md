# Quickstart: Calculator UI Barebones

**Feature**: 001-calculator-ui-barebones

## Prerequisites

- **Node 20 LTS or newer** (local machine currently has Node 16 — upgrade required; Vite 5+/shadcn need ≥18).
- **Go 1.23 or newer** (not currently installed locally — install required).
- **GNU Make** (present).

## One-command workflow (from repo root)

```bash
make install     # frontend npm install + backend go mod download
make dev         # run frontend (Vite) and backend (Go) together
make test        # run frontend (Vitest) + backend (go test) suites
make lint        # eslint + go vet / golangci-lint
make build       # production build of both workspaces
make fmt         # prettier + gofmt
make clean       # remove build artifacts
```

Per-stack targets also exist: `make run-frontend`, `make run-backend`, `make test-frontend`,
`make test-backend`.

## Verify the feature end-to-end

1. `make install` — completes without error for both workspaces.
2. `make test` — frontend component/flow tests and backend handler tests all pass (green).
3. `make dev` — open the Vite URL (default http://localhost:5173):
   - Press `1 2 + 7 =` → display shows `19`.
   - Press `.` twice in one number → only one decimal appears.
   - Divide by zero → clear error state shown; press `C` → back to `0`.
   - Resize the browser narrow→wide → keypad stays aligned, no horizontal scroll.
4. Backend health: `curl localhost:8080/healthz` → `{"status":"ok"}`;
   `curl localhost:8080/readyz` → `{"status":"ready"}`.

## Layout created by this feature

```text
frontend/   # React + TS + Vite + Tailwind + shadcn, Vitest tests
backend/    # Go module, net/http health/readiness, table-driven tests
contracts/  # API contract placeholder (backend OpenAPI; future calc contract)
Makefile    # root orchestration
```
