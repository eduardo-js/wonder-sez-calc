---
description: "Task list for Calculator UI Barebones"
---

# Tasks: Calculator UI Barebones

**Input**: Design documents from `/specs/001-calculator-ui-barebones/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/, quickstart.md

**Tests**: INCLUDED — the spec explicitly requires component render tests (FR-008) and
calculation/validation tests (FR-009), and the constitution mandates test-first (Principle II).

**Organization**: Tasks grouped by user story (US1 calculation, US2 validation, US3 responsive)
for independent implementation and testing.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no incomplete dependencies)
- **[Story]**: US1 / US2 / US3 (setup, foundational, polish have no story label)

## Path Conventions

Monorepo web app per plan.md: `frontend/` (Vite/React/TS), `backend/` (Go module), `contracts/`,
root `Makefile`.

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Monorepo scaffolding, toolchain, orchestration.

- [X] T001 Create monorepo root layout (`frontend/`, `backend/`, `contracts/`) and `contracts/README.md` pointing to `specs/001-calculator-ui-barebones/contracts/backend-openapi.yaml`
- [X] T002 [P] Add `.nvmrc` at repo root containing the latest Node LTS (`24`)
- [X] T003 Scaffold frontend with Vite (React + TypeScript) in `frontend/` (`package.json`, `index.html`, `src/main.tsx`, `src/App.tsx`)
- [X] T004 Install and configure Tailwind CSS in `frontend/` (`tailwind.config.ts`, `postcss.config.js`, Tailwind directives in `frontend/src/index.css`)
- [X] T005 Initialize shadcn/ui in `frontend/` (`components.json`, `frontend/src/lib/utils.ts` with `cn()`) and vendor the `Button` primitive into `frontend/src/components/ui/button.tsx`
- [X] T006 [P] Configure Vitest + React Testing Library + jsdom in `frontend/vite.config.ts` and `frontend/src/test/setup.ts` (add `@testing-library/jest-dom`, `@testing-library/user-event`)
- [X] T007 [P] Configure ESLint + Prettier and strict TypeScript (`strict: true`, no implicit `any`) in `frontend/.eslintrc.cjs`, `frontend/.prettierrc`, `frontend/tsconfig.json`
- [X] T008 [P] Initialize Go module in `backend/` (`go.mod`, module path) with package dirs `backend/cmd/server/` and `backend/internal/health/`
- [X] T009 Create root `Makefile` with phony targets: `install`, `build`, `test`, `test-frontend`, `test-backend`, `lint`, `fmt`, `run-frontend`, `run-backend`, `dev`, `clean`, `all` (delegating into `frontend/` via npm and `backend/` via go)

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Shared component shells, state infrastructure, and backend scaffolding required by all stories.

**⚠️ CRITICAL**: No user story work begins until this phase is complete.

- [X] T010 [P] Define calculator types and initial state (`CalculatorState`, action union) in `frontend/src/hooks/calculatorTypes.ts` per data-model.md
- [X] T011 [P] Create the keypad config array (`CalculatorButton` entries: digits, operators, decimal, equals, clear) in `frontend/src/lib/keypad.ts`
- [X] T012 Create `useCalculator` reducer skeleton (all actions wired, no compute/validation logic yet) in `frontend/src/hooks/useCalculator.ts` (depends on T010)
- [X] T013 [P] Create `CalculatorButton` component shell (props: label/value/variant/onPress/ariaLabel) in `frontend/src/components/calculator/CalculatorButton.tsx`
- [X] T014 [P] Create `CalculatorDisplay` component shell (props: value/error) in `frontend/src/components/calculator/CalculatorDisplay.tsx`
- [X] T015 Create `CalculatorKeypad` component rendering one `CalculatorButton` per keypad config entry in `frontend/src/components/calculator/CalculatorKeypad.tsx` (depends on T011, T013)
- [X] T016 Create `Calculator` container wiring `useCalculator` → display + keypad, mount in `frontend/src/App.tsx` in `frontend/src/components/calculator/Calculator.tsx` (depends on T012, T014, T015)
- [X] T017 Write table-driven test for health/readiness handlers (httptest) in `backend/internal/health/handler_test.go`
- [X] T018 Implement `/healthz` + `/readyz` JSON handlers in `backend/internal/health/handler.go` (make T017 pass)
- [X] T019 Implement server entrypoint with `log/slog` JSON logging and `http.Server` in `backend/cmd/server/main.go` (depends on T018)

**Checkpoint**: App shell renders, reducer dispatches, backend health endpoints pass — stories can begin.

---

## Phase 3: User Story 1 - Perform a basic calculation (Priority: P1) 🎯 MVP

**Goal**: User composes an arithmetic expression with buttons and gets the correct result.

**Independent Test**: Press `1 2 + 7 =` → display shows `19`; new digit after result starts fresh; `C` resets to `0`.

### Tests for User Story 1 (write first, ensure they FAIL) ⚠️

- [X] T020 [P] [US1] Unit tests for pure compute (`+ - * /`, result formatting, overwrite-after-equals) in `frontend/src/lib/calculator.test.ts`
- [X] T021 [P] [US1] `CalculatorButton` render test (correct label, variant class, fires `onPress(value)` once, accessible name) in `frontend/src/components/calculator/CalculatorButton.test.tsx`
- [X] T022 [P] [US1] `Calculator` integration test for the basic flow (`12+7=19`, new entry after result, clear→`0`) in `frontend/src/components/calculator/Calculator.test.tsx`

### Implementation for User Story 1

- [X] T023 [US1] Implement pure compute function (`+ - * /`) in `frontend/src/lib/calculator.ts` (make T020 pass)
- [X] T024 [US1] Implement `INPUT_DIGIT`, `CHOOSE_OP`, `EVALUATE`, `CLEAR` transitions in `frontend/src/hooks/useCalculator.ts` (depends on T023)
- [X] T025 [US1] Finalize `CalculatorButton` rendering (variant→shadcn styling, aria) in `frontend/src/components/calculator/CalculatorButton.tsx` (make T021 pass)
- [X] T026 [US1] Render display value/formatting in `frontend/src/components/calculator/CalculatorDisplay.tsx` (make T022 pass)

**Checkpoint**: Basic calculator fully functional and testable — MVP deliverable.

---

## Phase 4: User Story 2 - Trust the input is valid (Priority: P2)

**Goal**: Malformed input is prevented and division-by-zero shows a clear, recoverable error.

**Independent Test**: Second `.` is ignored; leading/repeated operators normalized; `÷0` shows error state; `C` recovers.

### Tests for User Story 2 (write first, ensure they FAIL) ⚠️

- [X] T027 [P] [US2] Unit tests for validation rules (single decimal, leading operator, repeated operator, divide-by-zero) in `frontend/src/lib/validation.test.ts`
- [X] T028 [P] [US2] `Calculator` integration test for validation + error recovery in `frontend/src/components/calculator/Calculator.validation.test.tsx`

### Implementation for User Story 2

- [X] T029 [US2] Implement validation helpers (decimal/operator guards) in `frontend/src/lib/validation.ts` (make T027 pass)
- [X] T030 [US2] Apply validation in `INPUT_DECIMAL` and operator transitions of `frontend/src/hooks/useCalculator.ts` (depends on T029)
- [X] T031 [US2] Add divide-by-zero `error` state in `frontend/src/lib/calculator.ts` + reducer, and render error indication (`role="alert"`) in `frontend/src/components/calculator/CalculatorDisplay.tsx`
- [X] T032 [US2] Ensure `CLEAR` recovers from `error` to initial state in `frontend/src/hooks/useCalculator.ts` (make T028 pass)

**Checkpoint**: US1 + US2 both work independently; input is trustworthy.

---

## Phase 5: User Story 3 - Use the calculator on any device (Priority: P3)

**Goal**: Layout adapts across mobile/tablet/desktop with reachable, legible, tappable controls.

**Independent Test**: Keypad renders all buttons; no horizontal scroll/overlap at sm/md/lg; tap targets adequate.

### Tests for User Story 3 (write first, ensure they FAIL) ⚠️

- [X] T033 [P] [US3] `CalculatorKeypad` test: renders one button per config entry and a grid container in `frontend/src/components/calculator/CalculatorKeypad.test.tsx`
- [X] T034 [P] [US3] Responsive structure test: keypad grid + button sizing classes present for breakpoints in `frontend/src/components/calculator/Calculator.responsive.test.tsx`

### Implementation for User Story 3

- [X] T035 [US3] Implement responsive CSS Grid (sm/md/lg) and adequate tap-target sizing in `frontend/src/components/calculator/CalculatorKeypad.tsx` (make T033 pass)
- [X] T036 [US3] Responsive display sizing + long-result overflow/truncation in `frontend/src/components/calculator/CalculatorDisplay.tsx` (make T034 pass)

**Checkpoint**: All three user stories independently functional and responsive.

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Quality gates and end-to-end verification.

- [X] T037 [P] Backend: ensure `gofmt`/`go vet` clean and add `golangci-lint` config in `backend/.golangci.yml`
- [X] T038 [P] Frontend: `make lint` and `make fmt` clean; `tsc --noEmit` type-check passes
- [X] T039 Verify all `Makefile` targets work end-to-end (`install`, `test`, `build`, `dev`, `clean`)
- [X] T040 Update repo `README.md` with monorepo dev instructions referencing `make` targets and `.nvmrc`
- [X] T041 Run `quickstart.md` end-to-end verification (calc flow, validation, responsive resize, `curl /healthz` + `/readyz`)

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies — start immediately.
- **Foundational (Phase 2)**: Depends on Setup — BLOCKS all user stories.
- **User Stories (Phase 3–5)**: All depend on Foundational. US1→US2→US3 by priority, but each is independently testable. US2/US3 build on shared reducer/components but layer in distinct behavior.
- **Polish (Phase 6)**: Depends on the targeted stories being complete.

### Within Each User Story

- Tests written FIRST and FAIL before implementation (constitution Principle II).
- Logic modules → reducer wiring → component rendering.

### Parallel Opportunities

- Setup: T002, T006, T007, T008 are [P] (distinct files) after T003 establishes `frontend/`.
- Foundational: T010, T011, T013, T014 [P]; backend T017–T019 parallel to frontend foundational work.
- Per story: all `[P]` test tasks run together; logic vs. component files often parallelizable.

---

## Parallel Example: User Story 1

```bash
# Tests first (parallel):
Task: "Unit tests for compute in frontend/src/lib/calculator.test.ts"            # T020
Task: "CalculatorButton render test in .../CalculatorButton.test.tsx"            # T021
Task: "Calculator basic-flow integration test in .../Calculator.test.tsx"        # T022

# Then implementation:
Task: "Implement compute in frontend/src/lib/calculator.ts"                       # T023
```

---

## Implementation Strategy

### MVP First (User Story 1 only)

1. Phase 1 Setup → 2. Phase 2 Foundational → 3. Phase 3 US1 → **STOP & VALIDATE** (`12+7=19`) → demo MVP.

### Incremental Delivery

Setup + Foundational → US1 (MVP) → US2 (validation) → US3 (responsive). Each story tested independently before the next.

### Parallel Team Strategy

After Foundational: frontend dev drives US1→US2→US3; a second dev can own backend scaffolding (T017–T019) and Makefile/quickstart verification in parallel.

---

## Notes

- [P] = different files, no incomplete dependencies.
- Backend (T017–T019) is scaffolding-only this feature (health/readiness); calc contract deferred.
- `.nvmrc` pins latest Node LTS (`24`) per request; quickstart still notes Node 20+ as the minimum.
- Verify each test fails before implementing; commit after each task or logical group.
