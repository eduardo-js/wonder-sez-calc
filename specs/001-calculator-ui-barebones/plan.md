# Implementation Plan: Calculator UI Barebones

**Branch**: `001-calculator-ui-barebones` | **Date**: 2026-06-02 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `specs/001-calculator-ui-barebones/spec.md`

## Summary

Stand up a monorepo barebones: a responsive React + TypeScript calculator UI built with
shadcn/ui (reusable button components, client-side calculation + input validation, component
and flow tests), plus a scaffolding-only Go backend (health/readiness endpoints), all
orchestrated by a root Makefile. Calculation runs client-side in v1; the frontend↔backend
calculation contract is deferred. Approach and tech choices are detailed in
[research.md](./research.md).

## Technical Context

**Language/Version**: TypeScript 5.x on Node 20 LTS (frontend); Go 1.23+ (backend)

**Primary Dependencies**: React 18+, Vite 5+, Tailwind CSS, shadcn/ui (Radix), Vitest, React
Testing Library (frontend); Go standard library `net/http` + `log/slog` (backend)

**Storage**: N/A (no persistence in v1)

**Testing**: Vitest + React Testing Library + jsdom (frontend); Go `testing` + `net/http/httptest`,
table-driven (backend)

**Target Platform**: Modern browsers (mobile/tablet/desktop) for the UI; Linux/container for the Go service

**Project Type**: Web application — monorepo with `frontend/` + `backend/`

**Performance Goals**: Instant UI feedback on button press (<100ms perceived); 60fps layout; not
a throughput-bound feature

**Constraints**: Responsive with no horizontal scroll/overlap at sm/md/lg; strict TypeScript
(`strict: true`, no implicit `any`); structured JSON logs + health/readiness on backend; no
cross-boundary imports

**Scale/Scope**: Single calculator screen (~4 reusable components + logic/validation modules);
backend = 2 operational endpoints. Small barebones scope.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Status | Notes |
|-----------|--------|-------|
| I. Monorepo Cohesion & Boundaries | ✅ PASS | `frontend/` + `backend/` independent workspaces; `contracts/` for shared contract; no cross-imports |
| II. Test-First (NON-NEGOTIABLE) | ✅ PASS | Tasks will order failing tests before implementation; Vitest+RTL (FE), table-driven `testing` (BE) |
| III. Maintainability & Code Quality | ✅ PASS | ESLint+Prettier, strict TS; gofmt/go vet/golangci-lint; single-responsibility modules |
| IV. Contract-Driven Integration | ✅ PASS (scoped) | OpenAPI for backend health endpoints; calc contract explicitly deferred (documented in spec assumptions) |
| V. Observability & Operability | ✅ PASS | Backend `/healthz` + `/readyz`, `log/slog` JSON logs, explicit error handling; FE surfaces errors (divide-by-zero) gracefully |

**Result**: PASS — no violations. Complexity Tracking not required.

## Project Structure

### Documentation (this feature)

```text
specs/001-calculator-ui-barebones/
├── plan.md              # This file
├── spec.md              # Feature spec
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output (backend OpenAPI + frontend component contract)
└── tasks.md             # Phase 2 output (/speckit-tasks — NOT created here)
```

### Source Code (repository root)

```text
frontend/
├── src/
│   ├── components/
│   │   ├── ui/                  # shadcn primitives (Button, etc.) — vendored
│   │   └── calculator/
│   │       ├── Calculator.tsx          # container, wires hook + display + keypad
│   │       ├── CalculatorDisplay.tsx
│   │       ├── CalculatorKeypad.tsx
│   │       └── CalculatorButton.tsx     # reusable button (FR-002)
│   ├── lib/
│   │   ├── calculator.ts        # pure compute (FR-003/005)
│   │   ├── validation.ts        # input rules (FR-004)
│   │   ├── keypad.ts            # button config array (data source)
│   │   └── utils.ts            # shadcn cn() helper
│   ├── hooks/
│   │   └── useCalculator.ts    # reducer (state transitions from data-model.md)
│   ├── App.tsx
│   ├── main.tsx
│   └── index.css               # Tailwind directives
├── tests/                       # or co-located *.test.tsx
│   ├── CalculatorButton.test.tsx
│   ├── CalculatorDisplay.test.tsx
│   ├── CalculatorKeypad.test.tsx
│   ├── Calculator.test.tsx     # integration flow (FR-009)
│   └── calculator.logic.test.ts # pure logic + validation (FR-009)
├── index.html
├── package.json
├── tsconfig.json
├── vite.config.ts              # + Vitest config
├── tailwind.config.ts
├── components.json             # shadcn config
└── .eslintrc / .prettierrc

backend/
├── cmd/server/main.go          # wires router + slog, starts http.Server
├── internal/health/
│   ├── handler.go              # /healthz + /readyz
│   └── handler_test.go         # table-driven httptest
├── go.mod
└── go.sum

contracts/                       # repo-root shared contract location (constitution)
└── README.md                   # points to feature OpenAPI; future calc contract lives here

Makefile                         # root orchestration (FR-011)
```

**Structure Decision**: Web-application monorepo (Option 2). `frontend/` is an independent npm
project (Vite/React/TS); `backend/` is an independent Go module; `contracts/` holds the shared
API contract per the constitution. The root `Makefile` is the single orchestration entry point
delegating into each workspace.

## Complexity Tracking

No constitution violations — section intentionally empty.
