# Phase 1 Data Model: CI & E2E Pipeline

No runtime/persistence data. "Entities" are pipeline configuration objects.

## Entity: Workflow (`.github/workflows/ci.yml`)

| Field        | Value / rule                                                              |
|--------------|---------------------------------------------------------------------------|
| name         | `CI`                                                                       |
| on           | `pull_request` (branches: `main`) + `push` (branches: `main`)             |
| permissions  | `contents: read` (least privilege; no secrets used)                       |
| concurrency  | group `ci-${{ github.ref }}`, `cancel-in-progress: true`                  |
| jobs         | `frontend`, `backend` (parallel), `e2e` (`needs: [frontend, backend]`)    |

## Entity: Job `frontend`

| Field    | Value / rule                                                                   |
|----------|--------------------------------------------------------------------------------|
| runs-on  | `ubuntu-latest`                                                                |
| setup    | checkout; `setup-node` (`node-version-file: .nvmrc`, `cache: npm`)             |
| steps    | `npm ci` → `npm run lint` → `npm run typecheck` → `npm run test` → `npm run build` (in `frontend/`) |
| pass/fail| any non-zero step fails the job                                                |

## Entity: Job `backend`

| Field    | Value / rule                                                                   |
|----------|--------------------------------------------------------------------------------|
| runs-on  | `ubuntu-latest`                                                                |
| setup    | checkout; `setup-go` (`go-version-file: backend/go.mod`)                       |
| steps    | `go vet ./...` → `make test-backend` (coverage gate ≥ 95%) → `go build ./cmd/server` |
| pass/fail| any non-zero step fails; coverage below threshold fails (FR-003)              |

## Entity: Job `e2e`

| Field      | Value / rule                                                                 |
|------------|------------------------------------------------------------------------------|
| runs-on    | `ubuntu-latest`                                                              |
| needs      | `[frontend, backend]` — only runs if both gates pass                         |
| setup      | checkout; `setup-node` (`.nvmrc`, cache); `npm ci` in `e2e/`; `npx playwright install --with-deps chromium` |
| start      | `docker compose up -d --build`                                               |
| wait       | poll `http://localhost:8080/healthz` until ok, bounded timeout (FR-006/008)  |
| run        | `npx playwright test` against `http://localhost:5173`                        |
| artifacts  | on `always()`: upload Playwright report/traces + `docker compose logs` (FR-007) |
| teardown   | `docker compose down -v` on `always()`                                       |

## Entity: Playwright config (`e2e/playwright.config.ts`)

| Field          | Value / rule                                                   |
|----------------|----------------------------------------------------------------|
| baseURL        | `http://localhost:5173`                                        |
| projects       | Chromium only (v1)                                             |
| trace          | `on-first-retry` (or `retain-on-failure`)                     |
| screenshot     | `only-on-failure`                                              |
| retries        | `2` in CI (flake tolerance, FR — edge case), `0` locally       |
| reporter       | `html` (uploaded as artifact)                                  |

## Entity: E2E spec (`e2e/tests/calculator.spec.ts`)

| Field        | Value / rule                                                     |
|--------------|------------------------------------------------------------------|
| scenario     | load UI → enter `2+2` → equals → expect result `4`               |
| assertion    | displayed result equals backend-computed value (real round trip) |
| diagnostic   | failure produces trace + screenshot artifact                     |

## Entity: Makefile target (edit)

| Target     | Action                                                              |
|------------|---------------------------------------------------------------------|
| `test-e2e` | `cd e2e && npm ci && npx playwright test` (assumes stack already up) |
| `e2e`      | bring up stack + wait healthz + `test-e2e` + tear down (local convenience) |

Both `## `-documented so `make help` lists them.

## Relationships / ordering

```
frontend ─┐
          ├─ both pass ─► e2e ─► (artifacts on failure)
backend  ─┘
```

Gate decision (overall workflow status) = AND of all jobs; used as the required merge check (FR-004).
