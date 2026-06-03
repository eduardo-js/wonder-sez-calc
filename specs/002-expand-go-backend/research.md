# Phase 0 Research: Expand Go Backend

All Technical Context items are resolved (no NEEDS CLARIFICATION). The spec Assumptions fixed
the stack; this document records the decisions and rejected alternatives.

## 1. HTTP framework

- **Decision**: `github.com/gin-gonic/gin` (v1.10.x).
- **Rationale**: User-preferred. Mature, widely used, built-in route groups, path params,
  middleware chain, and request binding/validation out of the box. Minimal ceremony to migrate
  the two existing probes.
- **Alternatives considered**:
  - `net/http` + `ServeMux` (current): no built-in binding/validation/middleware grouping;
    more hand-rolled plumbing for the standards we want.
  - chi / echo: comparable, but Gin was explicitly requested and bundles validator.

## 2. CORS

- **Decision**: `github.com/gin-contrib/cors`, configured from env.
- **Rationale**: Official Gin middleware; handles preflight (`OPTIONS`), allowed
  origins/methods/headers, and exposes a typed config. Avoids hand-written header logic.
- **Config**: `CORS_ALLOWED_ORIGINS` (comma-separated). Default `http://localhost:5173`
  (Vite dev). Only listed origins are granted access (FR-003) — no wildcard with credentials.
- **Allowed methods/headers**: `GET, POST, PUT, PATCH, DELETE, OPTIONS`; headers
  `Origin, Content-Type, Accept, Authorization`. MaxAge for preflight caching.
- **Alternatives considered**: hand-rolled CORS middleware — more error-prone for preflight
  edge cases; rejected.

## 3. Input validation — package vs. reflect

- **Decision**: `github.com/go-playground/validator/v10` via Gin's binding (`ShouldBindJSON`
  + `binding`/`validate` struct tags). **Package, not hand-rolled reflect.** (Confirmed by user.)
- **Rationale**: De-facto standard, already bundled with Gin; declarative struct tags
  (`required`, `min`, `max`, `oneof`, etc.); maintained; far less custom reflection code to
  test to ≥95%.
- **Pattern**: a thin `httpx.BindJSON[T]`-style helper decodes the body, runs validation, and
  on failure translates `validator.ValidationErrors` into the standard error envelope with a
  `fields` map. On malformed/empty/oversized JSON it returns the same envelope, 400.
- **Alternatives considered**: manual `reflect`-based validation — rejected (reinvents a solved
  problem, higher bug surface, harder to hit coverage target).

## 4. Routing structure

- **Decision**: One central `server.NewRouter(cfg, logger)` builds the `*gin.Engine`, installs
  the middleware chain (recovery → request logging → CORS), registers probes at top level
  (`/healthz`, `/readyz`) and application routes under a versioned `r.Group("/api/v1")`, and
  sets `NoRoute`/`NoMethod` handlers that emit the standard envelope (FR-008).
- **Rationale**: Single discoverable place to see all routes (FR-005); versioned prefix lets
  the API evolve without breaking probes; consistent middleware application.
- **Probe preservation**: same paths, same JSON bodies (`{"status":"ok"}` / `{"status":"ready"}`),
  same 200s (FR-009).

## 5. Error envelope

- **Decision**: single JSON shape for all client/server errors:
  ```json
  { "error": { "code": "validation_failed", "message": "…", "fields": { "amount": "required" } } }
  ```
  `fields` present only for validation errors. Codes: `validation_failed` (400),
  `not_found` (404), `method_not_allowed` (405), `internal` (500).
- **Rationale**: FR-007 consistency across endpoints; machine-readable `code` + human `message`;
  optional per-field detail for forms.
- **Implementation**: `httpx.Error` type + `httpx.WriteError(c, status, code, msg, fields)`;
  recovery middleware maps panics → `internal`/500 without leaking internals.

## 6. Logging & observability

- **Decision**: keep `log/slog` JSON; a request-logging middleware logs method, path, status,
  latency, and (on error) the error. Recovery middleware logs the panic with stack context
  then writes the envelope. CORS denials are visible via status in request logs.
- **Rationale**: Constitution V — structured logs, explicit error handling, no swallowed errors.
- **Note**: do not use Gin's default text logger; wire slog so logs stay JSON and consistent
  with `main.go`.

## 7. Test standards & coverage

- **Decision**: `github.com/stretchr/testify` — `require` for fatal preconditions, `assert`
  for non-fatal checks, `mock` for collaborator doubles. **Table-driven by default.**
  `gin.SetMode(gin.TestMode)` in tests; `httptest.NewRecorder` + `http.NewRequest` for HTTP.
- **Mocks**: prefer mocking external collaborators (e.g. a future readiness dependency) via
  `testify/mock`; handlers take interfaces so collaborators are injectable.
- **Coverage ≥95%**: `go test ./... -coverprofile=coverage.out -covermode=atomic`; a Makefile
  `test-backend` (or `cover-backend`) target computes `go tool cover -func` total and fails the
  build if `< 95%`. Tiny, untestable glue (e.g. `main` wiring) kept minimal so the threshold is
  achievable; `main` stays thin and is exercised indirectly where practical.
- **Rationale**: matches user-defined project standards; table-driven keeps cases dense and
  uniform; testify reduces assertion boilerplate.
- **Alternatives considered**: stdlib-only asserts (current) — verbose; rejected for new code
  per the agreed standard. Probe tests migrate to testify + table-driven.

## 8. Configuration

- **Decision**: small `config` package reading env with sane defaults: `ADDR` (`:8080`),
  `CORS_ALLOWED_ORIGINS` (`http://localhost:5173`), timeouts retained from current `main.go`.
- **Rationale**: FR-002 per-environment config without code changes; no config library needed
  for this scope (plain `os.Getenv` + parsing), fully unit-testable.

## Dependency justification (constitution dependency policy)

| Dependency | Need | Maintenance | License |
|------------|------|-------------|---------|
| gin-gonic/gin | HTTP framework, routing, binding (user-chosen) | Active, very widely used | MIT |
| gin-contrib/cors | Correct CORS/preflight handling | Active, official contrib | MIT |
| go-playground/validator/v10 | Declarative validation (Gin-bundled) | Active, standard | MIT |
| stretchr/testify | Assertions + mocks (project standard) | Active, standard | MIT |

All versions will be pinned in `go.mod`/`go.sum`.
