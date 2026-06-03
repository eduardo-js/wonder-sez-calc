---
description: "Task list for Expression Hints & UI Polish"
---

# Tasks: Expression Hints & UI Polish

**Input**: Design documents from `/specs/005-expression-hints/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/, quickstart.md

**Tests**: INCLUDED — constitution Principle II (Test-First).

**Organization**: Tasks grouped by user story (US1 P1, US2 P2, US3 P3). **Frontend-only, display-only feature** — the hint shows the in-progress expression text and performs **no** evaluation; the backend evaluates only on equals (no backend/contract change).

> **Design correction (2026-06-02):** A first pass implemented a live *result preview*
> (debounced backend calls, `hint`/`hintStatus` state, `api.ts` `signal`, a `debounce`
> helper). Per user direction the frontend must not evaluate — that machinery was
> **reverted/deleted** and replaced with a pure `buildHintExpression` helper. Tasks below
> reflect the final design.

## Format: `[ID] [P?] [Story] Description` — owner **(frontend-react)**

---

## Phase 1: Setup

- [X] T001 [P] (frontend-react) Confirm frontend baseline green: `make test-frontend` + `make lint` (57 tests)

---

## Phase 2: Foundational (Blocking Prerequisites)

- [X] T002 (frontend-react) Add pure `buildHintExpression(previous, operator, current, overwrite)` to `frontend/src/lib/expression.ts` (rules per data-model.md) + unit tests in `frontend/src/lib/expression.test.ts`

**Checkpoint**: Pure expression-text derivation available. No state/contract change.

---

## Phase 3: User Story 1 - See the expression being built (Priority: P1) 🎯 MVP

**Goal**: A muted secondary line shows the in-progress expression (`12` → `12 +` → `12 + 3`) while typing; the result appears only on equals.

**Independent Test**: Press `1 2 + 3` → hint line shows `12 + 3` with no result and no `calculate` call until `=`.

- [X] T003 [US1] (frontend-react) Failing test in `frontend/src/components/calculator/CalculatorDisplay.test.tsx`: `expression="12 + 3"` renders a muted secondary line (`role="status"`, `aria-live="polite"`) above the main value
- [X] T004 [P] [US1] (frontend-react) Failing integration test in `frontend/src/components/calculator/Calculator.test.tsx`: typing shows the expression line through stages; `calculate` is NOT called until `=`; after `=` the committed result shows
- [X] T005 [US1] (frontend-react) Add `expression?: string` prop to `frontend/src/components/calculator/CalculatorDisplay.tsx`; render muted secondary line above the main value when non-empty
- [X] T006 [US1] (frontend-react) Wire `expression={buildHintExpression(state.previous, state.operator, state.current, state.overwrite)}` in `frontend/src/components/calculator/Calculator.tsx`
- [X] T007 [US1] (frontend-react) Run suite; confirm GREEN, no eval-on-type

**Checkpoint**: Expression hint works; no frontend evaluation. MVP deliverable.

---

## Phase 4: User Story 2 - Graceful, non-distracting hint behavior (Priority: P2)

**Goal**: Hint hidden on fresh entry, cleared after equals, suppressed during error alerts.

**Independent Test**: Fresh calculator → no hint; error on equals → error shown, no competing hint; after equals → hint cleared.

- [X] T008 [US2] (frontend-react) Tests in `CalculatorDisplay.test.tsx`: empty `expression` → no hint line; error/errorMsg set → hint suppressed
- [X] T009 [US2] (frontend-react) Implement suppression in `CalculatorDisplay.tsx` (hide when error alert shown or `expression` empty); `buildHintExpression` returns `""` on fresh/post-equals state (covers cleared-after-equals via existing `overwrite`/null reset)
- [X] T010 [US2] (frontend-react) Run suite; confirm edge behaviors GREEN

**Checkpoint**: Hint is unobtrusive and trustworthy.

---

## Phase 5: User Story 3 - Clearer, more polished interface (Priority: P3)

**Goal**: Distinct hint/entry hierarchy; visible keyboard focus on keypad buttons; overflow truncation intact.

**Independent Test**: Tab the keypad → visible focus rings; hint line visibly secondary vs main line; long value truncates.

- [X] T011 [P] [US3] (frontend-react) Tests in `frontend/src/components/calculator/CalculatorButton.test.tsx`: `focus-visible:ring` affordance present; button keyboard-focusable and operable
- [X] T012 [P] [US3] (frontend-react) Add `focus-visible:ring-2 focus-visible:ring-sky-400 focus-visible:ring-offset-2 focus-visible:ring-offset-slate-800 focus-visible:outline-none` to `frontend/src/components/calculator/CalculatorButton.tsx`
- [X] T013 [US3] (frontend-react) Hierarchy polish in `CalculatorDisplay.tsx`: muted secondary hint line (`text-sm text-slate-400`, `flex-col`) above the primary value; truncation unchanged
- [X] T014 [US3] (frontend-react) Run suite; confirm GREEN

**Checkpoint**: All three stories verified.

---

## Phase 6: Polish & Cross-Cutting Concerns

- [X] T015 [P] (frontend-react) JSDoc on `buildHintExpression` and the `CalculatorDisplay` `expression` prop
- [X] T016 (frontend-react) Remove dead code from the reverted preview design: deleted `frontend/src/lib/debounce.ts` + test; reverted `calculatorTypes.ts` (no hint state), `useCalculator.ts` (no effect), `api.ts` (no `signal`)
- [X] T017 (frontend-react) `make test-frontend` (80 tests), `make lint`, `npm --prefix frontend run build` (strict tsc + bundle) — all green
- [ ] T018 (frontend-react) Run `quickstart.md` manual check against backend at `:8080`: expression line through `12 + 3`, no result until `=`, focus rings, distinct hint line

---

## Dependencies & Execution Order

- Setup (T001) → Foundational (T002, `buildHintExpression`) → US1 (T003–T007) → US2 (T008–T010) → US3 (T011–T014) → Polish (T015–T018).
- US3 button tasks (T011/T012) are independent of US1/US2 (separate file) and ran in parallel.
- US2/US3 display tasks share `CalculatorDisplay.tsx` with US1 → serialized there.

---

## Notes

- Frontend-only, display-only; **no backend/API/contract/state change**. The hint is pure derived text (`buildHintExpression`); evaluation stays on the existing equals path.
- Final state: 80 frontend tests passing, lint + strict build clean.
- Remaining: **T018** manual quickstart verification (requires backend running).
