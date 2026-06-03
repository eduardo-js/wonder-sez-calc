# Implementation Plan: Expand Go Backend

**Branch**: `002-expand-go-backend` | **Date**: 2026-06-02 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `specs/002-expand-go-backend/spec.md`

## Summary

Promote the scaffolding-only Go backend onto a real HTTP foundation: adopt **Gin** as the
router, apply **CORS** for the configured frontend origin(s) via middleware, add **declarative
request validation** (`go-playground/validator`, Gin-native struct tags) with a single
structured **error envelope**, and centralize **route registration** behind a versioned
`/api/v1` group while preserving the existing `/healthz` and `/readyz` behavior. Establish
backend **test standards** — `stretchr/testify` (assert/require/mock), mocks for collaborators,
table-driven as the default, coverage **≥95%** enforced from the Makefile. No business
(calculator) endpoints are introduced here; the validation/error/routing pattern is proven by
tests. Detailed decisions in [research.md](./research.md).

## Technical Context

**Language/Version**: Go 1.23+ (existing `backend/` module)

**Primary Dependencies**: `github.com/gin-gonic/gin` (HTTP framework + binding),
`github.com/gin-contrib/cors` (CORS middleware), `github.com/go-playground/validator/v10`
(validation, bundled via Gin binding), `github.com/stretchr/testify` (assert/require/mock).
Retain `log/slog` for structured logging.

**Storage**: N/A (no persistence in this feature)

**Testing**: Go `testing` + `testify` + `net/http/httptest`, table-driven; `gin.TestMode`;
coverage via `go test -coverprofile`

**Target Platform**: Linux/container HTTP service

**Project Type**: Web application — monorepo `frontend/` + `backend/` (this feature touches `backend/` only)

**Performance Goals**: Not throughput-bound; middleware overhead negligible. Correct CORS
preflight and validation rejection are the bar.

**Constraints**: Allowed origins configurable via env (no code change); structured JSON logs +
explicit error handling (no swallowed errors); single error envelope across all endpoints;
existing probe contracts unchanged; pinned/justified deps; **≥95%** backend coverage.

**Scale/Scope**: One central router + ~3 middleware (CORS, recovery→envelope, request log),
a validation/error helper package, migrated health handlers. Small, foundational scope.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Status | Notes |
|-----------|--------|-------|
| I. Monorepo Cohesion & Boundaries | ✅ PASS | Changes confined to `backend/`; no cross-boundary imports; shared contract stays in `contracts/` |
| II. Test-First (NON-NEGOTIABLE) | ✅ PASS | Tasks order failing table-driven tests before impl; testify + httptest; ≥95% coverage gate |
| III. Maintainability & Code Quality | ✅ PASS | gofmt/goimports, go vet, golangci-lint; single-responsibility packages; exported docs; pinned deps |
| IV. Contract-Driven Integration | ✅ PASS | OpenAPI extended with the standard `Error` schema + versioned API group; probes contract preserved |
| V. Observability & Operability | ✅ PASS | slog JSON request logging middleware; recovery→structured envelope; CORS denials and validation failures logged; no swallowed errors |

**Result**: PASS — no violations. Complexity Tracking not required.

## Project Structure

### Documentation (this feature)

```text
specs/002-expand-go-backend/
├── plan.md              # This file
├── spec.md              # Feature spec
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output (extended OpenAPI: Error schema + /api/v1 group)
└── tasks.md             # Phase 2 output (/speckit-tasks — NOT created here)
```

### Source Code (repository root)

```text
backend/
├── cmd/server/
│   └── main.go                 # build config, build router, start http.Server (graceful shutdown retained)
├── internal/
│   ├── config/
│   │   ├── config.go           # env config: Addr, AllowedOrigins (FR-002), timeouts
│   │   └── config_test.go      # table-driven: defaults, parsing, comma-split origins
│   ├── server/
│   │   ├── router.go           # central route registration (FR-005): engine, middleware chain, /api/v1 group, NoRoute/NoMethod
│   │   └── router_test.go      # table-driven httptest: routes, CORS, 404/405 envelope, validation via test handler
│   ├── middleware/
│   │   ├── cors.go             # gin-contrib/cors from config (FR-001/003)
│   │   ├── cors_test.go        # preflight allowed/denied, methods/headers
│   │   ├── logging.go          # slog JSON request logging (FR-010)
│   │   ├── logging_test.go
│   │   ├── recovery.go         # panic → 500 standard envelope (no swallowed errors)
│   │   └── recovery_test.go
│   ├── httpx/
│   │   ├── error.go            # Error envelope type + writers (FR-007/008)
│   │   ├── error_test.go
│   │   ├── bind.go             # BindJSON helper: decode + validate → envelope (FR-006)
│   │   └── bind_test.go        # table-driven: missing/wrong-type/out-of-range/malformed
│   └── health/
│       ├── handler.go          # migrated to gin.HandlerFunc; same bodies (FR-009)
│       └── handler_test.go     # table-driven, gin context
├── go.mod                      # + gin, gin-contrib/cors, validator, testify (pinned)
└── go.sum                      # generated

contracts/                      # repo-root shared contract location (constitution)
└── README.md                   # updated to point at the extended 002 OpenAPI

Makefile                        # test-backend gains coverage + ≥95% gate (FR-012)
```

**Structure Decision**: Keep the existing Go module layout. Routing concerns split by single
responsibility: `config` (env), `server` (central registration + middleware wiring),
`middleware` (CORS/logging/recovery), `httpx` (error envelope + bind/validate), `health`
(migrated handlers). `main.go` stays thin: load config → build router → serve with graceful
shutdown. No throwaway business endpoints; the validation/error pattern is proven via
`server`/`httpx` tests using an in-test handler registered on the real router.

## Complexity Tracking

No constitution violations — section intentionally empty.
