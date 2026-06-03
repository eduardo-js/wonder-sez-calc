# Implementation Plan: Expression Hints & UI Polish

**Branch**: `005-expression-hints` | **Date**: 2026-06-02 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `/specs/005-expression-hints/spec.md`

## Summary

Add an **expression hint**: a secondary display line that shows the *in-progress
expression text* as the user types (`12`, then `12 +`, then `12 + 3`), plus targeted UI
polish. The frontend does **not** evaluate the hint — evaluation stays a backend
responsibility, on equals (`POST /api/v1/calculate`, unchanged). The hint is **pure
derived state**: a `buildHintExpression(previous, operator, current, overwrite)` helper in
`expression.ts` computes the text from existing state; `CalculatorDisplay` renders it as a
muted secondary line above the main value, suppressed while an error alert shows. UI polish
adds clear hint/entry visual hierarchy and keyboard `focus-visible` affordances on the
keypad. **Frontend-only. No new state, no new deps, no contract change.**

> **Note**: an earlier iteration implemented a live *result preview* via debounced backend
> calls — removed per user direction (the frontend must not evaluate). See research.md.

## Technical Context

**Language/Version**: TypeScript 5.x + React 18 (frontend only)

**Primary Dependencies**: Vite, React, Vitest + Testing Library (all present). **No new deps.**

**Storage**: N/A (stateless; preview reuses the calculate endpoint)

**Testing**: Vitest + React Testing Library + `userEvent`; mock `lib/api.calculate`

**Target Platform**: Modern browsers via Vite dev (:5173); backend at :8080

**Project Type**: Web application (monorepo: `frontend/` + `backend/`) — this feature touches only `frontend/`

**Performance Goals**: Hint is pure derived text, recomputed on render — instant, no network, no jank.

**Constraints**: Hint performs no evaluation and triggers no request (FR-003); hint text must match the expression sent on equals (FR-005); hint suppressed during error alerts (FR-006); existing EVALUATE/error behavior unchanged.

**Scale/Scope**: ~4 frontend files edited (display + Calculator wiring + `expression.ts` helper + keypad button focus ring) plus tests. No state-shape change. Small surface.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Status | Notes |
|-----------|--------|-------|
| I. Monorepo Cohesion & Boundaries | ✅ | Frontend-only, display-only change; no backend edit and no new integration surface (hint never calls the backend). |
| II. Test-First (NON-NEGOTIABLE) | ✅ | Failing tests first: `buildHintExpression` unit cases and `CalculatorDisplay` expression-line rendering; confirm red, then implement. |
| III. Maintainability & Code Quality | ✅ | Hint is a pure helper + derived render; no new state, no async, no dead code; strict TS, ESLint/Prettier clean. |
| IV. Contract-Driven Integration | ✅ | No contract change and no `calculate` option change — the hint is display-only. UI contract for the display prop documented in `contracts/`. |
| V. Observability & Operability | ✅ | Hint is suppressed during error alerts so it never competes with the assertive error message; existing EVALUATE error handling unchanged. |

**No new dependencies** → no dependency justification required. **No violations** → Complexity Tracking empty.

## Project Structure

### Documentation (this feature)

```text
specs/005-expression-hints/
├── plan.md              # This file
├── spec.md              # Feature spec
├── research.md          # Phase 0
├── data-model.md        # Phase 1
├── quickstart.md        # Phase 1
├── contracts/
│   └── hint-ui.md       # Phase 1 — UI contract for hint state + display props
└── checklists/
    └── requirements.md
```

### Source Code (repository root)

```text
frontend/
└── src/
    └── lib/
        ├── expression.ts               # EDIT: add pure buildHintExpression helper
        └── expression.test.ts          # NEW: unit tests for buildHintExpression
    components/calculator/
        ├── CalculatorDisplay.tsx       # EDIT: secondary expression line + hierarchy
        ├── CalculatorDisplay.test.tsx  # EDIT: expression-line render states
        ├── CalculatorButton.tsx        # EDIT: focus-visible ring (a11y polish)
        ├── CalculatorButton.test.tsx   # EDIT: focus-ring + keyboard-operability tests
        ├── Calculator.tsx              # EDIT: pass expression={buildHintExpression(...)}
        └── Calculator.test.tsx         # EDIT: hint shows while typing; no eval until equals
# No change to hooks/ (state shape untouched) or lib/api.ts (no contract/option change).
```

**Structure Decision**: Web-application layout (Option 2), frontend workspace only. The
hint is **pure derived text** over the existing binary-op state, so it is computed by a
helper in `expression.ts` and rendered by `CalculatorDisplay` — the reducer/hook and the
HTTP contract stay untouched, preserving constitution principles I–IV.

## Key Design Decisions

- **Pure derivation, no eval**: `buildHintExpression(previous, operator, current, overwrite)` returns the in-progress expression text; the frontend never evaluates it (FR-003). No new state, no effect, no async, no debounce, no `calculate` call for the hint.
- **Operator/overwrite rule**: after `CHOOSE_OP`, `current` still holds the left operand while `overwrite` is true → show `previous operator` only; once the RHS is typed (`overwrite === false`) → `previous operator current`. Fresh/post-equals (`overwrite === true`, no operator) → `""` (hidden, avoids redundant `0`).
- **Display hierarchy**: hint renders as a smaller, muted secondary line (`role="status"` + `aria-live="polite"`) above the main value; suppressed while an error alert is shown (FR-006/FR-007).
- **Consistency (FR-005)**: hint text and the equals-time `buildExpression` output use the same operands/operator, so what you see is what gets evaluated.
- **No contract/state change**: `CalculatorState`, the reducer, `api.ts`, and the HTTP contract are untouched.

## Complexity Tracking

> No constitution violations. Section intentionally empty.
</content>
