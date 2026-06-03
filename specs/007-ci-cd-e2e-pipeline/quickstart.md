# Quickstart: CI & E2E Pipeline

## What CI does

On every PR and push to `main`, GitHub Actions runs:

1. **frontend** — `npm ci`, lint, typecheck, unit tests, build.
2. **backend** — `go vet`, tests + 95% coverage gate, build.
3. **e2e** (after 1 & 2 pass) — boots the Docker Compose stack, waits for `/healthz`, runs Playwright
   against `http://localhost:5173`, tears down. On failure, uploads Playwright report/traces + compose logs.

Overall status is the required merge gate. Superseded runs are auto-cancelled.

## Run e2e locally

Prereqs: Docker (feature 006 stack) + Node 24.

```bash
# Option A: one shot (brings stack up, runs e2e, tears down)
make e2e

# Option B: manual
make docker-up                 # start stack
cd e2e && npm ci
npx playwright install chromium
npx playwright test            # runs against http://localhost:5173
npx playwright show-report     # view results / traces
make docker-down               # stop stack
```

## E2E scenarios (v1)

`e2e/tests/calculator.spec.ts` runs 12 browser-driven scenarios against the live stack:

- Core round trip: `2 + 2 =` → `4` (proves browser → frontend → backend → result).
- All four operators, decimals, chained ops, clear/reset.
- Expression hint line (`role="status"`) and division-by-zero error alert (`role="alert"`) + recovery.
- Large-number scientific-notation round trip.

See `contracts/e2e-scenarios.md` for the full table.

## Verifying the pipeline itself

- Open a PR to `main` with a deliberately broken unit test → CI fails, merge blocked.
- Open a clean PR → CI passes, mergeable.
- Break the frontend↔backend contract (e.g. change the response shape) → the e2e job fails with a trace artifact.

## Parity

CI runs the same commands as local (`npm run lint/typecheck/test/build`, `go vet`, `make test-backend`,
`go build`). If it passes locally via `make lint test build` + `make e2e`, it passes in CI.

## Notes

- Pipeline needs no secrets (CI + e2e only; CD deferred) → forked-PR runs are safe.
- Chromium only for v1; multi-browser is a future addition.
