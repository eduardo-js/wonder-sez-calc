---
description: "Task list for 003-wire-frontend-backend"
---

# Tasks: Wire Frontend and Backend (Server-Side Calculation)

**Input**: Design documents from `/specs/003-wire-frontend-backend/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/calculate-openapi.yaml

**Tests**: INCLUDED — Test-First is NON-NEGOTIABLE per constitution Principle II. Write tests, confirm they fail, then implement.

**Organization**: Tasks grouped by user story (P1 → P2 → P3) for independent implementation and testing.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: US1 / US2 / US3 (maps to spec.md user stories)

## Path Conventions

Web-app monorepo: `backend/internal/...`, `backend/cmd/...`, `frontend/src/...`

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Confirm both workspaces build/test green before changes.

- [X] T001 Verify backend builds and tests pass: `make test-backend`
- [X] T002 [P] Verify frontend builds and tests pass: `make test-frontend`
- [X] T003 [P] Add `VITE_API_BASE_URL` (default `http://localhost:8080`) to `frontend/.env.example` and document in `frontend/README.md` if present

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Shared contract/error scaffolding both backend and frontend stories depend on.

**⚠️ CRITICAL**: No user story work begins until this phase completes.

- [X] T004 Add `CodeCalculation = "calculation_error"` constant to `backend/internal/httpx/error.go`
- [X] T005 [P] Define frontend API/result/error TypeScript types in `frontend/src/lib/apiTypes.ts` matching `contracts/calculate-openapi.yaml` (`CalculationRequest`, `CalculationResult`, `ApiError`)
- [X] T006 [P] Extend `CalculatorState` in `frontend/src/hooks/calculatorTypes.ts` with `status: "idle" | "loading" | "error"` and `errorMsg: string | null`; update `INITIAL_STATE`

**Checkpoint**: Error code + shared types exist — user stories can begin.

---

## Phase 3: User Story 1 - Compute a valid expression via the backend (Priority: P1) 🎯 MVP

**Goal**: Frontend sends an expression to the backend; backend evaluates valid expressions with correct precedence and returns JSON; frontend displays the backend result. No client-side arithmetic.

**Independent Test**: Enter `2 + 3 * 4`, press `=`, see `14` from the backend response; confirm no local compute path runs.

### Tests for User Story 1 ⚠️ (write first, confirm failing)

- [X] T007 [P] [US1] Table-driven evaluator tests (happy path: precedence, parens, decimals, unary minus) in `backend/internal/calc/evaluator_test.go`
- [X] T008 [P] [US1] Handler test for `POST /api/v1/calculate` 200 success in `backend/internal/calculate/handler_test.go`
- [X] T009 [P] [US1] Router test asserting `/api/v1/calculate` is registered in `backend/internal/server/router_test.go`
- [X] T010 [P] [US1] Frontend api client test (valid expression → result) in `frontend/src/lib/api.test.ts`
- [X] T011 [P] [US1] Hook test: `EVALUATE` calls api and displays backend result in `frontend/src/hooks/useCalculator.test.ts`

### Implementation for User Story 1

- [X] T012 [US1] Implement pure evaluator (tokenize → shunting-yard → eval; `+ - * /`, parens, decimals, unary minus) in `backend/internal/calc/evaluator.go`
- [X] T013 [US1] Implement result formatter (integers plain; else ≤12 significant digits) in `backend/internal/calc/evaluator.go`
- [X] T014 [US1] Implement calculate handler: bind `CalculationRequest`, evaluate, return `CalculationResult` (200) in `backend/internal/calculate/handler.go`
- [X] T015 [US1] Add `Register(group *gin.RouterGroup)` for `POST /calculate` in `backend/internal/calculate/register.go`
- [X] T016 [US1] Mount `calculate.Register` on the `/api/v1` group in `backend/internal/server/router.go`
- [X] T017 [P] [US1] Implement typed `calculate()` client (fetch, JSON) in `frontend/src/lib/api.ts`
- [X] T018 [P] [US1] Add expression builder from operands in `frontend/src/lib/expression.ts`
- [X] T019 [US1] Rewire `useCalculator` `EVALUATE` to call `calculate()` async and set display from response in `frontend/src/hooks/useCalculator.ts`
- [X] T020 [US1] Remove client arithmetic: delete `frontend/src/lib/calculator.ts` and `frontend/src/lib/calculator.test.ts`; remove all `compute` imports/usages

**Checkpoint**: Valid expressions compute on the backend and display correctly. MVP demoable.

---

## Phase 4: User Story 2 - Receive and surface validation errors (Priority: P2)

**Goal**: Invalid/undefined input → backend returns structured error; frontend shows an accessible alert; display not updated with a result.

**Independent Test**: Submit `5 + * 2` → 400 `validation_failed`; submit `1/0` → 422 `calculation_error`; both surface a user alert.

### Tests for User Story 2 ⚠️ (write first, confirm failing)

- [X] T021 [P] [US2] Evaluator error tests (malformed, unbalanced parens, empty, divide-by-zero, non-finite → typed sentinels) in `backend/internal/calc/evaluator_test.go`
- [X] T022 [P] [US2] Handler tests: malformed expr → 400 `validation_failed` (fields.expression), divide-by-zero/non-finite → 422 `calculation_error`, bad JSON → 400 `bad_request` in `backend/internal/calculate/handler_test.go`
- [X] T023 [P] [US2] api client test: error envelope parsed and thrown/returned as typed error in `frontend/src/lib/api.test.ts`
- [X] T024 [P] [US2] Component test: backend error renders `role="alert"` with message; display unchanged in `frontend/src/components/calculator/Calculator.validation.test.tsx`

### Implementation for User Story 2

- [X] T025 [US2] Add typed error sentinels (`ErrInvalidExpression`, `ErrDivideByZero`, `ErrNonFinite`) and return them from the evaluator in `backend/internal/calc/evaluator.go`
- [X] T026 [US2] Map evaluator errors → wire codes/HTTP in handler (`validation_failed`/400 with `fields.expression`, `calculation_error`/422) plus input length/charset validation in `backend/internal/calculate/handler.go`
- [X] T027 [US2] Parse error envelope in `calculate()` and surface as typed error in `frontend/src/lib/api.ts`
- [X] T028 [US2] Set `status:"error"`/`errorMsg` on failed `EVALUATE` (display unchanged) in `frontend/src/hooks/useCalculator.ts`; clear on next input/CLEAR
- [X] T029 [US2] Add accessible alert region (`role="alert"`) bound to `errorMsg` in `frontend/src/components/calculator/CalculatorDisplay.tsx`

**Checkpoint**: Invalid input is rejected by the backend and surfaced to the user; US1 still works.

---

## Phase 5: User Story 3 - Prevent concurrent requests with in-progress feedback (Priority: P3)

**Goal**: While a request is in flight, the `=` control shows a spinner and is disabled; duplicate presses issue no requests; control restores on resolve, including network failure/timeout.

**Independent Test**: Trigger a calc; during pending state the `=` shows a spinner and repeated presses send no extra requests; on resolve/timeout the control restores.

### Tests for User Story 3 ⚠️ (write first, confirm failing)

- [X] T030 [P] [US3] Hook test: second `EVALUATE` while `loading` issues no request; spinner state toggles on resolve in `frontend/src/hooks/useCalculator.test.ts`
- [X] T031 [P] [US3] api client test: `AbortController` timeout rejects with connectivity error in `frontend/src/lib/api.test.ts`
- [X] T032 [P] [US3] Component test: `=` button renders spinner + `disabled` while loading, restores after in `frontend/src/components/calculator/CalculatorButton.test.tsx`

### Implementation for User Story 3

- [X] T033 [US3] Add in-flight guard (ignore `EVALUATE` when `status==="loading"`) and set `loading` on dispatch in `frontend/src/hooks/useCalculator.ts`
- [X] T034 [US3] Add `AbortController` + timeout (default 5s) to `calculate()`; map abort/network failure to connectivity error in `frontend/src/lib/api.ts`
- [X] T035 [US3] On resolve/error/timeout, restore `status:"idle"` and re-enable submissions in `frontend/src/hooks/useCalculator.ts`
- [X] T036 [US3] Render spinner in `=` control and set `disabled` while loading in `frontend/src/components/calculator/CalculatorButton.tsx` (and wire prop in `CalculatorKeypad.tsx`)

**Checkpoint**: Single in-flight request enforced; spinner shows/restores; failures never leave the control stuck.

---

## Phase 6: Polish & Cross-Cutting Concerns

- [X] T037 [P] Verify structured request logging covers calculate (incl. validation failures) via existing middleware; add a handler log line if missing in `backend/internal/calculate/handler.go`
- [X] T038 [P] Run `gofmt`/`go vet`/lint (backend) and ESLint/Prettier/`tsc --noEmit` (frontend); fix issues
- [X] T039 Confirm backend coverage ≥95%: `cd backend && go test ./... -cover`
- [X] T040 Run `quickstart.md` end-to-end validation (curl cases + UI checks)

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: no dependencies.
- **Foundational (Phase 2)**: depends on Setup; BLOCKS all user stories.
- **User Stories (Phase 3–5)**: depend on Foundational. US2 builds on US1's evaluator/handler/hook; US3 builds on US1's request path. Recommended order P1 → P2 → P3.
- **Polish (Phase 6)**: depends on all targeted stories.

### User Story Dependencies

- **US1 (P1)**: after Foundational. Independent MVP.
- **US2 (P2)**: after US1 (extends evaluator errors, handler mapping, hook/UI for alerts) — independently testable via error inputs.
- **US3 (P3)**: after US1 (extends request lifecycle) — independently testable via in-flight behavior.

### Within Each User Story

- Tests first and failing → implement → green.
- Backend: evaluator → handler → router.
- Frontend: api client/types → hook → component.

### Parallel Opportunities

- Setup: T002, T003.
- Foundational: T005, T006 (T004 independent backend file).
- US1 tests T007–T011 in parallel; impl T017/T018 parallel to backend T012–T016.
- US2 tests T021–T024 parallel; US3 tests T030–T032 parallel.
- Backend and frontend tracks can proceed in parallel within a story (different files).

---

## Parallel Example: User Story 1

```bash
# Tests first (parallel):
Task: "Evaluator happy-path tests in backend/internal/calc/evaluator_test.go"
Task: "Handler 200 test in backend/internal/calculate/handler_test.go"
Task: "Router registration test in backend/internal/server/router_test.go"
Task: "api client test in frontend/src/lib/api.test.ts"
Task: "Hook displays backend result test in frontend/src/hooks/useCalculator.test.ts"

# Then implementation — backend and frontend in parallel:
Task: "Evaluator in backend/internal/calc/evaluator.go"
Task: "calculate() client in frontend/src/lib/api.ts"
Task: "expression builder in frontend/src/lib/expression.ts"
```

---

## Implementation Strategy

### MVP First (User Story 1)

1. Phase 1 Setup → 2. Phase 2 Foundational → 3. Phase 3 US1 → 4. STOP & validate `2 + 3 * 4 → 14` from backend → 5. demo.

### Incremental Delivery

1. Foundation ready → 2. US1 (MVP, backend-computed results) → 3. US2 (errors + alerts) → 4. US3 (lock + spinner). Each increment is independently testable and adds value without breaking prior stories.

### Parallel Team Strategy

After Foundational: one dev on backend (`calc`/`calculate`), one on frontend (`lib`/`hooks`/`components`); integrate per story at the contract boundary.

---

## Notes

- [P] = different files, no incomplete dependencies.
- Constitution: tests fail before implementation; gofmt/vet/lint + strict TS clean; ≥95% backend coverage; contract is source of truth.
- Commit after each task or logical group.
- FR coverage: FR-001 (T020), FR-002/004/007 (US1), FR-003/005/006/008 (US2), FR-009/010/011/012 (US3), FR-013 (contract-aligned types T005 + handler), FR-014 (T037).
