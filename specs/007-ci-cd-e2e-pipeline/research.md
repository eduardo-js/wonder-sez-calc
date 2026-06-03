# Phase 0 Research: CI & E2E Pipeline

## R1. Pipeline platform & structure

**Decision**: Single workflow `.github/workflows/ci.yml` with three jobs: `frontend`, `backend`
(parallel), and `e2e` (`needs: [frontend, backend]`). Triggers: `pull_request` → `main` and `push`
→ `main`.

**Rationale**: Parallel quality jobs give fast feedback; gating e2e on both avoids spending runner
time on e2e when a cheaper gate already failed. Matches FR-001/FR-002 and the < 15 min budget (SC-005).

**Alternatives considered**: One monolithic job (`make all`) — simpler but serial and slower; separate
workflow files per job — more files, harder to read for a 3-job pipeline.

## R2. Local/CI parity (SC-006, FR-013)

**Decision**: CI runs the exact commands the local workflow uses:
- frontend job: `npm ci` → `npm run lint` → `npm run typecheck` → `npm run test` → `npm run build`
- backend job: `go vet ./...` → `make test-backend` (the coverage-gate target) → `go build ./cmd/server`

**Rationale**: These are the same underlying commands the Makefile invokes (`make lint/test/build`).
Reusing `make test-backend` keeps the 95% coverage gate identical to local (FR-003). Calling the npm
scripts directly (rather than `make`) lets frontend/backend run as independent jobs while staying
command-identical; the Makefile's `nvm` guard is a no-op when Node is already on PATH, so `make`
targets also remain CI-safe.

**Alternatives considered**: Invoke `make lint test build` in one job — maximal parity but serial and
mixes workspaces; rejected for the parallel split (still command-identical per workspace).

## R3. Toolchain setup & caching

**Decision**: `actions/setup-node@v4` with `node-version-file: .nvmrc` and `cache: npm`;
`actions/setup-go@v5` with `go-version-file: backend/go.mod` and built-in module caching.

**Rationale**: Pinning to `.nvmrc` (24) and `go.mod` (1.25) guarantees CI uses the same versions as
local. Built-in caching keeps runs within the time budget without custom cache keys.

**Alternatives considered**: Hardcode versions in YAML — drifts from `.nvmrc`/`go.mod`; manual
`actions/cache` — unnecessary now that setup actions cache natively.

## R4. E2E tool

**Decision**: **Playwright** (`@playwright/test`), Chromium only for v1, in a dedicated top-level
`e2e/` package.

**Rationale**: First-class GitHub Actions support, built-in trace viewer + screenshots/video on
failure (satisfies FR-007/SC-004), auto-waiting (less flake), and `webServer`/baseURL config. A
separate `e2e/` package keeps it out of the frontend Vitest scope (Principle I) so unit and e2e
runners stay distinct.

**Alternatives considered**: Cypress — solid but heavier CI image and weaker trace tooling than
Playwright traces; Selenium — more boilerplate, more flake. Putting Playwright inside `frontend/` —
would entangle e2e specs with the Vitest glob and component test config.

## R5. Booting the stack in CI for e2e

**Decision**: The e2e job uses the feature-006 Docker Compose stack: `docker compose up -d --build`,
then poll `http://localhost:8080/healthz` until healthy (bounded timeout), run Playwright against
`http://localhost:5173`, then `docker compose down -v` in an `always()` step.

**Rationale**: Reuses the real, already-tested run topology (no bespoke CI-only startup), exercising
the genuine frontend↔backend contract (FR-005, Principle IV). `/healthz` is the existing readiness
signal (FR-006); the compose healthcheck already gates the frontend on backend health.

**Alternatives considered**: Run `vite preview` + `go run` directly on the runner — diverges from the
real stack and re-implements wiring; Playwright `webServer` launching dev servers — same divergence,
and wouldn't validate the container path.

**Fail-fast (FR-008)**: the health poll has a max attempt count / timeout; on timeout the job prints
`docker compose logs` and exits non-zero rather than hanging.

## R6. Failure diagnostics / artifacts (FR-007, SC-004)

**Decision**: On e2e failure, upload the Playwright HTML report + traces/screenshots and the
`docker compose logs` output via `actions/upload-artifact@v4` with `if: always()`.

**Rationale**: Lets a maintainer diagnose from the run alone, no local repro (SC-004).

## R7. Concurrency & re-runnability (FR-010, FR-011)

**Decision**: Top-level `concurrency: { group: ci-${{ github.ref }}, cancel-in-progress: true }`.
Jobs are stateless and deterministic, so re-running a commit yields the same decision (FR-010).

**Rationale**: Cancels superseded runs on new pushes to the same branch/PR (FR-011), saving runner
minutes; statelessness gives idempotent gating.

## R8. Forked-PR secret safety (FR-009)

**Decision**: The pipeline needs **no secrets** (CI + e2e only; CD deferred). Use the default
`GITHUB_TOKEN` with least privilege (`permissions: contents: read`). Because no secret is referenced,
forked-PR runs execute all gates safely.

**Rationale**: With CD out of scope there is no registry/deploy credential to leak; explicit minimal
`permissions` documents intent and avoids accidental exposure if secrets are added later.

**Alternatives considered**: `pull_request_target` (has secret access) — unnecessary and riskier;
rejected.

## R9. Branch reality

**Decision**: Gate `main`. The current branch will be pushed and set as the repo's default branch
`main` (user-confirmed); no `master`/`main` dual-trigger needed.

**Rationale**: Single protected default avoids ambiguous gating; matches the resolved spec (FR-001).

## E2E scenario scope (v1)

- Primary journey: load UI → input an expression (e.g. `2+2`) → activate equals → assert displayed
  result (`4`) — proving the browser→backend→browser round trip.
- One negative/diagnostic check: simulate/observe a failure path producing an artifact (kept minimal).
- Cross-browser/device matrices: out of scope v1 (Chromium only).

All NEEDS CLARIFICATION resolved (CD deferred, default branch = `main`).
