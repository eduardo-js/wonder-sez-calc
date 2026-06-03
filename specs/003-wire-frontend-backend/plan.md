# Implementation Plan: Wire Frontend and Backend (Server-Side Calculation)

**Branch**: `003-wire-frontend-backend` | **Date**: 2026-06-02 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `/specs/003-wire-frontend-backend/spec.md`

## Summary

Move all arithmetic to the backend. The frontend submits a raw expression string to a new `POST /api/v1/calculate` endpoint; the backend validates and evaluates it (standard operator precedence) and returns a JSON result or the standard error envelope. The frontend displays the result, surfaces errors in an accessible alert region, shows a spinner on the `=` control while a request is in flight, and blocks concurrent requests. Client-side `compute()` is removed.

## Technical Context

**Language/Version**: Go 1.22+ (backend); TypeScript 5.x + React 18 (frontend)

**Primary Dependencies**: Backend — gin, go-playground/validator (both already present; **no new deps**). Frontend — Vite, React, Vitest + Testing Library (already present); `fetch` + `AbortController` (web standard).

**Storage**: N/A (stateless calculation)

**Testing**: Go `testing` (table-driven, ≥95% coverage); Vitest + Testing Library

**Target Platform**: Linux server (backend :8080); modern browsers via Vite dev (:5173)

**Project Type**: Web application (monorepo: `frontend/` + `backend/`)

**Performance Goals**: Spinner appears/clears within 100 ms of press/resolve (SC-005); calculation latency dominated by network, not compute.

**Constraints**: Single in-flight request (SC-004); client timeout default 5s; expression ≤256 chars; no client-side arithmetic (FR-001).

**Scale/Scope**: Single endpoint; one frontend integration path. Small surface, high correctness bar.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Status | Notes |
|-----------|--------|-------|
| I. Monorepo Cohesion & Boundaries | ✅ | Work split across `frontend/` and `backend/`; no cross-boundary imports; integration via HTTP/JSON only. |
| II. Test-First (NON-NEGOTIABLE) | ✅ | Evaluator + handler tests and frontend api/hook tests written before impl; div-by-zero etc. as table cases. |
| III. Maintainability & Code Quality | ✅ | Pure evaluator, pure reducer; side effects isolated in hook. gofmt/vet/lint + strict TS. |
| IV. Contract-Driven Integration | ✅ | `contracts/calculate-openapi.yaml` is source of truth; FE client + BE handler validated against it; reuses 0.2.0 error envelope, adds `calculation_error`. |
| V. Observability & Operability | ✅ | Existing structured-log + recovery middleware applies to the new route; errors explicit (typed sentinels → wire codes); FE surfaces failures via alert region. |

**No new dependencies** → no dependency justification required. **No violations** → Complexity Tracking empty.

## Project Structure

### Documentation (this feature)

```text
specs/003-wire-frontend-backend/
├── plan.md              # This file
├── spec.md              # Feature spec
├── research.md          # Phase 0
├── data-model.md        # Phase 1
├── quickstart.md        # Phase 1
├── contracts/
│   └── calculate-openapi.yaml   # Phase 1
└── checklists/
    └── requirements.md
```

### Source Code (repository root)

```text
backend/
├── internal/
│   ├── calc/                 # NEW: pure expression evaluator
│   │   ├── evaluator.go      #   tokenize → shunting-yard → eval; typed errors
│   │   └── evaluator_test.go #   table-driven (precedence, parens, errors)
│   ├── calculate/            # NEW: HTTP handler + DTOs
│   │   ├── handler.go        #   bind → evaluate → map errors → JSON
│   │   ├── handler_test.go
│   │   └── register.go       #   Register(group) under /api/v1
│   ├── httpx/error.go        # EDIT: add CodeCalculation = "calculation_error"
│   └── server/router.go      # EDIT: mount calculate.Register on /api/v1 group
└── cmd/server/main.go        # unchanged

frontend/
├── src/
│   ├── lib/
│   │   ├── api.ts            # NEW: typed calculate() client (fetch+AbortController+timeout)
│   │   ├── api.test.ts       # NEW
│   │   ├── expression.ts     # NEW: build expression string from operands
│   │   └── calculator.ts     # REMOVE (client arithmetic) + delete calculator.test.ts
│   ├── hooks/
│   │   ├── calculatorTypes.ts# EDIT: add status + errorMsg; EVALUATE no longer computes
│   │   └── useCalculator.ts  # EDIT: async EVALUATE via api.ts; lock; spinner state
│   └── components/calculator/
│       ├── CalculatorButton.tsx   # EDIT: spinner when `=` is loading + disabled
│       ├── CalculatorDisplay.tsx  # EDIT/ADD: alert region (role="alert")
│       └── *.test.tsx             # EDIT: cover spinner, lock, alert, backend result
```

**Structure Decision**: Web-application layout (Option 2). New backend packages `calc` (pure domain) and `calculate` (HTTP adapter) keep evaluation logic testable and separate from transport (constitution III). Frontend isolates the side-effecting request in `lib/api.ts` and the hook, keeping the reducer pure.

## Complexity Tracking

> No constitution violations. Section intentionally empty.
