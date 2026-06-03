---

description: "Task list for Expand Go Backend"
---

# Tasks: Expand Go Backend

**Input**: Design documents from `/specs/002-expand-go-backend/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/, quickstart.md

**Tests**: INCLUDED — the project constitution mandates Test-First (NON-NEGOTIABLE). Write each
test, confirm it fails, then implement to green.

**Organization**: Grouped by user story (US1 CORS, US2 Validation, US3 Routing & Standards).

**Status: ✅ IMPLEMENTED** (2026-06-02). All phases done via fanned-out `backend-go` agents.
Backend `internal/` coverage **98.8%** (gate ≥95% in `make test-backend`); `go vet`/`gofmt` clean;
end-to-end smoke test green (probes preserved, 404/405 envelopes, CORS allow/deny). T003 (shared
`testutil`) was skipped — agents used `net/http/httptest` directly; no shared helper needed.
T002 (golangci config) was already present from feature 001.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files/packages, no incomplete dependencies)
- **[Story]**: US1 / US2 / US3 (omitted for Setup / Foundational / Polish)
- All paths are relative to repo root.

## Path Conventions

Go module at `backend/` (module `github.com/wonderlic-calc/backend`). Packages under
`backend/internal/`, entrypoint `backend/cmd/server/`.

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Dependencies and tooling so every package can compile and test.

- [ ] T001 Add backend dependencies in `backend/` — `go get github.com/gin-gonic/gin github.com/gin-contrib/cors github.com/go-playground/validator/v10 github.com/stretchr/testify`, then `go mod tidy` (creates `backend/go.sum`, pins versions)
- [ ] T002 [P] Add `backend/.golangci.yml` enabling gofmt/goimports/govet/errcheck (constitution III)
- [ ] T003 [P] Add `backend/internal/testutil/gin.go` test helper (`gin.SetMode(gin.TestMode)`, `PerformRequest(engine, method, path, headers, body)` returning `*httptest.ResponseRecorder`)

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Config, error envelope, middleware skeleton, router, and migrated probes that ALL
stories build on.

**⚠️ CRITICAL**: No user story work begins until this phase is complete.

### Tests (write first, confirm failing)

- [ ] T004 [P] `backend/internal/config/config_test.go` — table-driven: default `ADDR`, custom `ADDR`, retained timeouts
- [ ] T005 [P] `backend/internal/httpx/error_test.go` — table-driven: `WriteError` sets status, `Content-Type: application/json`, body `{"error":{code,message,fields?}}`; `fields` omitted when nil
- [ ] T006 [P] `backend/internal/middleware/logging_test.go` — request log emits method/path/status/latency as JSON via slog
- [ ] T007 [P] `backend/internal/middleware/recovery_test.go` — a panicking handler yields 500 with `internal` envelope; panic logged; no leaked internals
- [ ] T008 [P] `backend/internal/health/handler_test.go` — table-driven (gin): `GET /healthz` → 200 `{"status":"ok"}`, `GET /readyz` → 200 `{"status":"ready"}`

### Implementation

- [ ] T009 [P] `backend/internal/config/config.go` — `Config{Addr string; AllowedOrigins []string; Read/Write/IdleTimeout}`; `Load()` reads `ADDR` (default `:8080`) + timeouts; `AllowedOrigins` field present (populated in US1)
- [ ] T010 [P] `backend/internal/httpx/error.go` — `Error{Code,Message,Fields}`, envelope `{"error":Error}`, code consts (`validation_failed/bad_request/not_found/method_not_allowed/internal`), `WriteError(c *gin.Context, status int, code, msg string, fields map[string]string)`
- [ ] T011 [P] `backend/internal/middleware/logging.go` — `RequestLogger(logger *slog.Logger) gin.HandlerFunc` (JSON: method, path, status, latency)
- [ ] T012 `backend/internal/middleware/recovery.go` — `Recovery(logger *slog.Logger) gin.HandlerFunc` mapping panics → `httpx.WriteError(...,internal,500)` (depends on T010)
- [ ] T013 [P] `backend/internal/health/handler.go` — migrate to `gin.HandlerFunc`s (`Healthz`, `Readyz`); same JSON bodies/status; `Register(r gin.IRouter)` for `GET /healthz`, `GET /readyz`
- [ ] T014 `backend/internal/server/router.go` — `NewRouter(cfg config.Config, logger *slog.Logger) *gin.Engine`: `gin.New()`, chain Recovery → RequestLogger, register health probes (depends on T009–T013)
- [ ] T015 `backend/cmd/server/main.go` — refactor: `config.Load()` → `server.NewRouter()` → `http.Server` with retained graceful shutdown/timeouts (depends on T014)

**Checkpoint**: `make run-backend` serves; `/healthz` + `/readyz` unchanged; `make test-backend` green for foundational packages.

---

## Phase 3: User Story 1 - Frontend can call the backend across origins (Priority: P1) 🎯 MVP

**Goal**: Configured frontend origin(s) can call the API in-browser (incl. preflight); other
origins are denied (FR-001/002/003).

**Independent Test**: Preflight `OPTIONS` from `http://localhost:5173` returns
`Access-Control-Allow-Origin` for it; a disallowed origin does not.

### Tests (write first, confirm failing)

- [ ] T016 [P] [US1] `backend/internal/middleware/cors_test.go` — table-driven: allowed origin preflight grants origin/methods/headers; disallowed origin not granted; allowed method vs disallowed method
- [ ] T017 [P] [US1] `backend/internal/config/config_test.go` — add cases: `CORS_ALLOWED_ORIGINS` comma-split, trimmed, empty → default `["http://localhost:5173"]`

### Implementation

- [ ] T018 [US1] `backend/internal/config/config.go` — parse `CORS_ALLOWED_ORIGINS` into `AllowedOrigins` (comma-split, trim, drop empties, default `http://localhost:5173`)
- [ ] T019 [US1] `backend/internal/middleware/cors.go` — `CORS(origins []string) gin.HandlerFunc` via `gin-contrib/cors` (methods `GET,POST,PUT,PATCH,DELETE,OPTIONS`; headers `Origin,Content-Type,Accept,Authorization`; MaxAge 12h; no wildcard-with-credentials)
- [ ] T020 [US1] `backend/internal/server/router.go` — insert `CORS(cfg.AllowedOrigins)` into the middleware chain; `backend/internal/server/router_test.go` — preflight allowed/denied through the real router

**Checkpoint**: Browser CORS works for configured origin; denied otherwise.

---

## Phase 4: User Story 2 - Malformed requests rejected with clear errors (Priority: P2)

**Goal**: Declarative validation rejects bad input with the standard envelope before any handler
logic; valid input passes (FR-006/007).

**Independent Test**: POST valid + each invalid class (missing/wrong-type/out-of-range/malformed)
to a validated test endpoint; valid accepted, others → 400 with consistent envelope.

### Tests (write first, confirm failing)

- [ ] T021 [P] [US2] `backend/internal/httpx/bind_test.go` — table-driven: missing required → `validation_failed`+`fields`; wrong type / malformed JSON → `bad_request`; out-of-range (`gt`) and `oneof` → `validation_failed`; valid → no error
- [ ] T022 [P] [US2] `backend/internal/server/router_test.go` — register an in-test handler using `BindJSON`; assert each invalid class returns the standard envelope and valid passes

### Implementation

- [ ] T023 [US2] `backend/internal/httpx/bind.go` — `BindJSON(c *gin.Context, dst any) bool`: `ShouldBindJSON` → on `validator.ValidationErrors` build `fields` map + `validation_failed`; on decode/syntax error → `bad_request`; writes envelope and returns false on failure

**Checkpoint**: All invalid-input classes rejected uniformly; valid input proceeds.

---

## Phase 5: User Story 3 - Consistent, discoverable routing & quality standards (Priority: P3)

**Goal**: All routes centrally registered with a versioned `/api/v1` group; unknown route/method
return the standard envelope; tests are table-driven with ≥95% coverage from one command
(FR-005/008/011/012).

**Independent Test**: Inspect `router.go` (all routes in one place); `GET /nope` → 404 envelope,
`DELETE /healthz` → 405 envelope; `make test-backend` green with coverage ≥95%.

### Tests (write first, confirm failing)

- [ ] T024 [P] [US3] `backend/internal/server/router_test.go` — add cases: unknown route → 404 `not_found` envelope; bad method on known path → 405 `method_not_allowed` envelope; probes unchanged after full chain
- [ ] T025 [P] [US3] `backend/internal/server/router_test.go` — assert `/api/v1` group exists and applies the same middleware (probe via a registered test route under the group)

### Implementation

- [ ] T026 [US3] `backend/internal/server/router.go` — add `NoRoute` → `httpx.WriteError(...not_found,404)`, `NoMethod` (`HandleMethodNotAllowed=true`) → `...method_not_allowed,405`; create `r.Group("/api/v1")` (empty, ready for future routes)
- [ ] T027 [US3] `Makefile` — change `test-backend` to run `go test ./... -coverprofile=coverage.out -covermode=atomic` then fail if `go tool cover -func` total `< 95.0%`; add `cover-backend` (HTML report)
- [ ] T028 [P] [US3] `backend/README.md` — document routing conventions (central registration, `/api/v1` group, error envelope, table-driven + testify + ≥95% coverage standard)

**Checkpoint**: Central routing verified; 404/405 standardized; coverage gate enforced.

---

## Phase 6: Polish & Cross-Cutting Concerns

- [ ] T029 [P] Run `gofmt -w .`, `goimports -w .`, `go vet ./...`, `golangci-lint run` in `backend/` — clean
- [ ] T030 Run `cd backend && go mod tidy`; confirm `go.mod`/`go.sum` pinned and minimal
- [ ] T031 Verify coverage ≥95%: `make test-backend` green; close gaps with table-driven cases
- [ ] T032 Run `specs/002-expand-go-backend/quickstart.md` curl checks against a running server (CORS allowed/denied, probes, 404/405)
- [ ] T033 [P] Confirm `contracts/backend-openapi.yaml` (v0.2.0) matches implemented envelope/probes; update if drift

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (P1)**: none — start immediately (T001 unblocks all compilation)
- **Foundational (P2)**: after Setup — BLOCKS all user stories
- **User Stories (P3–P5)**: after Foundational; US1/US2 are independent; US3 builds on the router
- **Polish (P6)**: after desired stories

### Within Foundational

- T010 (error envelope) before T012 (recovery uses it) and before US3 404/405
- T009–T013 before T014 (router) before T015 (main)

### User Story Dependencies

- **US1 (P1)**: Foundational only. Touches config.go + new cors.go + router chain.
- **US2 (P2)**: Foundational only. Self-contained in httpx (+ router test). Independent of US1.
- **US3 (P3)**: Foundational + uses router; verifies central registration & coverage. Best last (depends on final router shape).

### Shared-file notes (NOT parallel with each other)

- `config/config.go`: T009 → T018 (sequential)
- `server/router.go`: T014 → T020 → T026 (sequential)
- `server/router_test.go`: T020/T022/T024/T025 append to same file (serialize)

---

## Parallel Opportunities (by package — safe to fan out)

**Foundational tests wave** (different files):
```
T004 config_test · T005 httpx error_test · T006 logging_test · T007 recovery_test · T008 health handler_test
```
**Foundational impl wave** (different packages; T012 after T010):
```
T009 config · T010 httpx/error · T011 middleware/logging · T013 health   →  then T012 recovery → T014 router → T015 main
```
**Stories** (after Foundational): US1 and US2 packages are disjoint (`middleware`+`config` vs `httpx`) → can run in parallel; US3 last.

---

## Implementation Strategy

### MVP (US1)
1. Phase 1 Setup → 2. Phase 2 Foundational → 3. Phase 3 US1 → validate CORS → demo.

### Incremental
Foundational → US1 (CORS) → US2 (validation) → US3 (routing/standards/coverage) → Polish. Each is an independently testable increment.

### Parallel fan-out (this run)
- Wave A: T001 (deps) — orchestrator
- Wave B (parallel agents): config | httpx | middleware(logging/recovery) | health — tests-first then impl per package
- Wave C: server/router + main (wires B)
- Wave D (parallel): US1 CORS additions | US2 bind/validation
- Wave E: US3 router (404/405, /api/v1), Makefile gate, README; then Polish + coverage close-out

---

## Notes

- [P] = different files/packages, no incomplete deps.
- Verify each test FAILS before implementing (constitution II).
- During waves, run `go test ./internal/<pkg>` (not `./...`) until all packages exist.
- Keep `main.go` thin so ≥95% coverage is achievable; exercise router via `server` tests.
- Commit after each phase/checkpoint.
