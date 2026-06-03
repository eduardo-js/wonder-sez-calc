# Research: Calculator UI Barebones

**Feature**: 001-calculator-ui-barebones | **Date**: 2026-06-02

Consolidated technology decisions for the barebones monorepo (React frontend + Go backend,
Makefile-orchestrated). Each item: Decision / Rationale / Alternatives considered.

## 1. Frontend build tool

- **Decision**: **Vite** (React + TypeScript template).
- **Rationale**: shadcn/ui officially supports Vite; fast dev server and build; first-class
  Vitest integration sharing one config. Aligns with constitution's strict-TS requirement.
- **Alternatives considered**: Create React App (deprecated/unmaintained); Next.js (SSR/router
  overhead unneeded for a single-screen client-side calculator).

## 2. UI component library

- **Decision**: **shadcn/ui** (Radix UI primitives + Tailwind CSS), components vendored into
  `frontend/src/components/ui/`.
- **Rationale**: Explicitly requested. Copy-in (not a dependency) model gives full control and
  testability. `Button` primitive is the base for the reusable calculator buttons.
- **Alternatives considered**: MUI / Chakra (heavier runtime, not requested); raw Radix without
  shadcn (loses the styling/variant conventions the user asked for).

## 3. Styling

- **Decision**: **Tailwind CSS** with shadcn's tokens; CSS Grid for the button layout.
- **Rationale**: Required by shadcn. Tailwind responsive breakpoints (`sm`/`md`/`lg`) directly
  satisfy FR-007 responsive support. CSS Grid gives a stable, aligned calculator keypad.
- **Alternatives considered**: CSS Modules / plain CSS (incompatible with shadcn conventions).

## 4. Frontend testing

- **Decision**: **Vitest + React Testing Library + jsdom**, with `@testing-library/jest-dom`
  and `@testing-library/user-event`.
- **Rationale**: Constitution mandates a component/unit runner with Testing Library. Vitest
  reuses the Vite config. RTL satisfies FR-008 (render assertions) and FR-009 (flow/validation).
- **Alternatives considered**: Jest (separate Babel/transform config, slower with Vite); Cypress
  component testing (heavier, out of scope for barebones).

## 5. Calculator logic placement

- **Decision**: **Pure TypeScript module** in `frontend/src/lib/` (e.g. `calculator.ts`,
  `validation.ts`), consumed by a `useCalculator` hook. No backend dependency in v1.
- **Rationale**: Keeps calculation/validation unit-testable in isolation (FR-003/004/005),
  decoupled from React rendering. Matches the spec assumption that v1 computes client-side.
- **Alternatives considered**: Logic inside components (untestable in isolation); compute via Go
  backend now (premature; contract deferred per spec assumption).

## 6. Input validation approach

- **Decision**: Validate at the **reducer/state-transition layer** — each button press is an
  action; illegal transitions (second decimal, leading operator, repeated operator,
  divide-by-zero) are rejected or normalized before state updates.
- **Rationale**: Centralizes FR-004/FR-005 rules in one testable place, independent of UI.
- **Alternatives considered**: Per-component guards (scattered, hard to test); validating only
  on equals (lets malformed intermediate state exist).

## 7. Backend language/runtime & framework

- **Decision**: **Go (1.23+)**, standard library `net/http` only, exposing `/healthz`
  (liveness) and `/readyz` (readiness) returning structured JSON; structured logging via
  `log/slog`.
- **Rationale**: Constitution requires structured (JSON) logs and health/readiness endpoints
  (Principle V). Backend is scaffolding-only this feature, so stdlib avoids premature framework
  lock-in. `slog` is stdlib structured logging.
- **Alternatives considered**: Gin/Echo/Chi (unnecessary for two endpoints; adds a dependency to
  justify per constitution); no backend at all (violates FR-010/FR-012 monorepo scaffolding).

## 8. Backend testing

- **Decision**: Standard `testing` package, **table-driven** tests with `net/http/httptest`
  for the health/readiness handlers.
- **Rationale**: Constitution mandates table-driven Go tests. `httptest` verifies handlers
  without binding a port.
- **Alternatives considered**: testify (extra dependency, not needed for stdlib handlers).

## 9. Monorepo layout & boundaries

- **Decision**: Top-level `frontend/` (npm project) and `backend/` (Go module), plus a
  `contracts/` location at the repo root for the future API contract (placeholder this feature).
  No cross-boundary imports.
- **Rationale**: Matches constitution's canonical layout (Principle I, Technology & Structure).
  Each workspace independently buildable/testable.
- **Alternatives considered**: npm/Go workspaces tooling (overkill for two packages); single
  shared package (violates boundary principle).

## 10. Orchestration (Makefile)

- **Decision**: Root **Makefile** with phony targets that delegate to each workspace:
  `install`, `build`, `test`, `lint`, `fmt`, `run-frontend`, `run-backend`, `dev`, `clean`,
  plus aggregate `all`. Frontend targets `cd frontend && npm ...`; backend targets
  `cd backend && go ...`.
- **Rationale**: Satisfies FR-011 single orchestration entry point and SC-005 (one-command
  install/build/test/run for both stacks).
- **Alternatives considered**: Taskfile/Just (not requested; Make is specified); npm scripts
  only (can't cleanly orchestrate Go).

## 11. Toolchain prerequisites (local environment gap)

- **Decision**: Require **Node 20 LTS+** and **Go 1.23+**; document in quickstart.
- **Rationale**: Local machine currently has Node 16 and no Go (verified). Vite 5+/shadcn need
  Node 18+; using 20 LTS is the safe baseline. This is an environment prerequisite, not a design
  blocker.
- **Alternatives considered**: Pinning to Node 16 (unsupported by current Vite/shadcn).

## Resolved unknowns

All Technical Context items are resolved; no NEEDS CLARIFICATION remain.
