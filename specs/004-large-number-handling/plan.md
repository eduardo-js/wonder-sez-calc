# Implementation Plan: Large-Number & Edge-Case Handling

**Branch**: `004-large-number-handling` | **Date**: 2026-06-02 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `/specs/004-large-number-handling/spec.md`

## Summary

Fix the result-formatting defect that silently corrupts large-magnitude results. The evaluator computes correctly in `float64` but `formatResult` casts integer-valued results to `int64`, which saturates above the int64 range — turning `≈ -1e42` into `-9223372036854775808`. Replace the unconditional int64 cast with a magnitude-gated decision: render a plain integer only when the value is within the exact-integer range (|v| ≤ 2⁵³); otherwise emit scientific notation via `%g` (≤12 significant digits). Truly unrepresentable results (overflow to ±Inf, division by zero) keep returning the existing explicit errors. Frontend already renders the backend result string verbatim and its display already truncates/auto-sizes; work there is confined to test coverage confirming `e`-notation and error strings display correctly. No new dependencies, no breaking contract change.

## Technical Context

**Language/Version**: Go 1.22+ (backend); TypeScript 5.x + React 18 (frontend)

**Primary Dependencies**: Backend — stdlib only (`strconv`, `math`); existing gin/validator unchanged. Frontend — Vite, React, Vitest + Testing Library (already present). **No new deps.**

**Storage**: N/A (stateless calculation)

**Testing**: Go `testing` (table-driven, maintain ≥95% coverage); Vitest + Testing Library

**Target Platform**: Linux server (backend :8080); modern browsers via Vite dev (:5173)

**Project Type**: Web application (monorepo: `frontend/` + `backend/`)

**Performance Goals**: Formatting is O(1); no measurable latency change.

**Constraints**: No regression for in-range whole numbers and existing non-integer formatting; zero silently-wrong results; expression ≤256 chars (unchanged).

**Scale/Scope**: One function changed in `backend/internal/calc/evaluator.go` plus tests; frontend test-only additions. Small surface, high correctness bar.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Status | Notes |
|-----------|--------|-------|
| I. Monorepo Cohesion & Boundaries | ✅ | Change is backend-internal (`calc`); frontend untouched except tests; integration via existing HTTP/JSON contract. |
| II. Test-First (NON-NEGOTIABLE) | ✅ | Bug fix → **regression test first**: add the reported expression as a failing table case (asserting it is NOT `-9223372036854775808`), confirm red, then fix `formatResult`. |
| III. Maintainability & Code Quality | ✅ | Single-responsibility change to a pure function; documented threshold constant; gofmt/vet/lint + strict TS. |
| IV. Contract-Driven Integration | ✅ | Result field stays a string; scientific notation already satisfies it. Contract gets a PATCH (0.3.0→0.3.1): clarified `result` description + large-number example. No breaking change. |
| V. Observability & Operability | ✅ | Unrepresentable results still surface explicit typed errors (`ErrNonFinite`, `ErrDivideByZero`) → wire `calculation_error`; FE alert region unchanged. |

**No new dependencies** → no dependency justification required. **No violations** → Complexity Tracking empty.

## Project Structure

### Documentation (this feature)

```text
specs/004-large-number-handling/
├── plan.md              # This file
├── spec.md              # Feature spec
├── research.md          # Phase 0
├── data-model.md        # Phase 1
├── quickstart.md        # Phase 1
├── contracts/
│   └── calculate-openapi.yaml   # Phase 1 (PATCH: 0.3.1)
└── checklists/
    └── requirements.md
```

### Source Code (repository root)

```text
backend/
└── internal/
    └── calc/
        ├── evaluator.go        # EDIT: formatResult — magnitude-gated integer vs scientific
        └── evaluator_test.go   # EDIT: add regression + boundary + large/small table cases

frontend/
└── src/
    └── components/calculator/
        └── CalculatorDisplay.test.tsx   # EDIT: assert e-notation result renders fully/legibly
```

**Structure Decision**: Web-application layout (Option 2). The defect is isolated to one pure function in the backend `calc` domain package; fixing it there keeps transport and frontend untouched (constitution I, III). Frontend changes are test-only because it already renders the result string verbatim with truncation/auto-sizing.

## Complexity Tracking

> No constitution violations. Section intentionally empty.
