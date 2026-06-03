---
description: "Task list for Large-Number & Edge-Case Handling"
---

# Tasks: Large-Number & Edge-Case Handling

**Input**: Design documents from `/specs/004-large-number-handling/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/, quickstart.md

**Tests**: INCLUDED — constitution Principle II (Test-First, NON-NEGOTIABLE) requires a failing test before each fix.

**Organization**: Tasks grouped by user story (US1 P1, US2 P2, US3 P3).

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel — only across **different files/workspaces** with no incomplete dependency.
- Owning subagent noted in each task: **(backend-go)** or **(frontend-react)**.

## Subagent Fan-Out Map

| Subagent | Owns | Files |
|----------|------|-------|
| `backend-go` | All backend logic, tests, contract | `backend/internal/calc/*`, `specs/.../contracts/` |
| `frontend-react` | Display verification tests | `frontend/src/components/calculator/*.test.tsx` |

> ⚠️ **Same-file serialization**: US1, US2, US3 all touch `backend/internal/calc/evaluator.go` (`formatResult`) and `evaluator_test.go`. Within the backend these stories are **sequential**, not parallel. The only true cross-story parallelism is **backend ↔ frontend** (separate workspaces).

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Confirm green baseline before changing anything.

- [X] T001 [P] (backend-go) Confirm backend baseline is green: run `make -C backend test` and `make -C backend lint`; record current coverage %
- [X] T002 [P] (frontend-react) Confirm frontend baseline is green: run `make -C frontend test`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Shared artifact used by all three stories.

**⚠️ CRITICAL**: Blocks US1–US3.

- [X] T003 (backend-go) Add documented constant `maxExactInt = 1 << 53` (// 9007199254740992) in `backend/internal/calc/evaluator.go` near `formatResult` (no logic change yet)

**Checkpoint**: Constant available; story work can begin (backend serially, frontend in parallel).

---

## Phase 3: User Story 1 - Large results never silently wrong (Priority: P1) 🎯 MVP

**Goal**: The reported expression and any out-of-range integer result return the correct value (scientific notation) or an explicit error — never the corrupted `-9223372036854775808`.

**Independent Test**: POST `-9223372036854775808 - 999999999999999999999999999999999999999999`; assert result is NOT `-9223372036854775808` (it is `-1e+42`).

### Tests for User Story 1 (write first, must FAIL) ⚠️

- [X] T004 [US1] (backend-go) Add failing regression table case in `backend/internal/calc/evaluator_test.go`: `-9223372036854775808 - 999999999999999999999999999999999999999999` → expect `-1e+42` AND assert `!= "-9223372036854775808"`; run and confirm RED

### Implementation for User Story 1

- [X] T005 [US1] (backend-go) Fix `formatResult` in `backend/internal/calc/evaluator.go`: gate the `strconv.FormatInt(int64(v),10)` branch on `math.Abs(v) <= maxExactInt`; else `strconv.FormatFloat(v,'g',12,64)`
- [X] T006 [US1] (backend-go) Run `make -C backend test`; confirm T004 now GREEN and no existing case regressed (e.g. `12 + 7 * 3` → `33`)

**Checkpoint**: Corruption fixed and verified — MVP deliverable.

---

## Phase 4: User Story 2 - Readable scientific notation (Priority: P2)

**Goal**: Large and very-small results render in compact scientific notation in the API and on screen.

**Independent Test**: Evaluate a very large and a very small result; confirm both are scientific notation in the payload and fully visible in the UI.

### Tests for User Story 2 (write first) ⚠️

- [X] T007 [US2] (backend-go) Add table cases in `backend/internal/calc/evaluator_test.go`: boundary `9007199254740992` → plain `9007199254740992`; `9007199254740992 + 9007199254740992` → scientific; `1e40 * 2` → `2e+40`; `1 / 1e20` → `1e-20` (write first; T005 already makes large cases pass, small/boundary lock behavior)
- [X] T008 [P] [US2] (frontend-react) Add Vitest case in `frontend/src/components/calculator/CalculatorDisplay.test.tsx`: a scientific-notation result string (e.g. `-1e+42`) renders fully/legibly (no clipping, sizing branch applies)

### Implementation for User Story 2

- [X] T009 [US2] (backend-go) Run `make -C backend test`; confirm all T007 cases pass (no code change expected beyond T005 — investigate if any fail)
- [X] T010 [P] [US2] (backend-go) Update `specs/004-large-number-handling/contracts/calculate-openapi.yaml` is already PATCHed to 0.3.1; verify examples match actual output (`-1e+42`, `1e-20`) and adjust if Go `%g` differs

**Checkpoint**: Scientific notation verified backend + frontend.

---

## Phase 5: User Story 3 - Explicit errors for unrepresentable results (Priority: P3)

**Goal**: Overflow-to-infinity and division-by-zero return explicit errors surfaced by the UI.

**Independent Test**: Evaluate `1e400` and `1 / 0`; both return `calculation_error` (422) and the UI shows the message.

### Tests for User Story 3 (write first) ⚠️

- [X] T011 [US3] (backend-go) Add table cases in `backend/internal/calc/evaluator_test.go`: `1e400` → `ErrNonFinite`; `1 / 0` → `ErrDivideByZero` (lock existing behavior)
- [X] T012 [P] [US3] (frontend-react) Add/extend Vitest case in `frontend/src/components/calculator/CalculatorDisplay.test.tsx`: an error message renders in `role="alert"` without corrupting prior value

### Implementation for User Story 3

- [X] T013 [US3] (backend-go) Run `make -C backend test`; confirm error cases pass (no code change expected — these paths already exist; investigate if any fail)

**Checkpoint**: All three stories independently verified.

---

## Phase 6: Polish & Cross-Cutting Concerns

- [X] T014 [P] (backend-go) Ensure `formatResult` doc comment describes the threshold + scientific-notation rule in `backend/internal/calc/evaluator.go`
- [X] T015 (backend-go) Run `make -C backend lint` + `make -C backend test`; confirm coverage ≥95% (matches/improves baseline from T001)
- [X] T016 [P] (frontend-react) Run `make -C frontend test` + lint; confirm green
- [X] T017 (backend-go) Execute `quickstart.md` manual curl check; confirm `{"result":"-1e+42",...}`
- [X] T018 (backend-go) **(discovered at integration boundary)** Allow `e`/`E` in the HTTP handler charset gate `allowedCharsRE` in `backend/internal/calculate/handler.go` so scientific-notation input literals reach the evaluator (FR-006); add API table cases in `handler_test.go` (large result `-1e+42`, `1e3+0`→`1000`, `1/1e20`→`1e-20`, `1e400`→422)

> **Note (T018)**: The end-to-end curl in T017 exposed that the handler's `allowedCharsRE` rejected `e`/`E` before the evaluator ran — so `1 / 1e20` and `1e400` returned `validation_failed` instead of evaluating. The evaluator already supported `e`; only the handler gate was stricter than the contract. Fixed by widening the regex to `^[0-9eE+\-*/().\ ]+$`.

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (P1)** → no deps.
- **Foundational (P2)** → after Setup; T003 blocks all stories.
- **US1 (P3)** → after T003. **MVP.**
- **US2 (P4)** → backend tasks (T007, T009, T010) depend on T005 (same file). Frontend task T008 depends only on T002.
- **US3 (P5)** → backend tasks depend on T005. Frontend task T012 depends only on T002.
- **Polish (P6)** → after all desired stories.

### Cross-Story Reality

- US1 → US2 → US3 **backend** tasks are **serial** (shared `evaluator.go` / `evaluator_test.go`).
- T005 (the single `formatResult` fix) mechanically satisfies US2's large-number behavior; US2/US3 backend tasks are mostly **lock-in tests**.

### Parallel Opportunities (true [P])

- T001 ∥ T002 (backend vs frontend baseline).
- Frontend display tests T008, T012 run anytime after T002 — fully parallel to the backend stream.
- T010 (contract verify) ∥ frontend tasks.
- Polish: T015 ∥ T016.

---

## Parallel Example: Subagent Fan-Out

```text
# Wave 1 — baseline (parallel):
backend-go:     T001  (make -C backend test/lint)
frontend-react: T002  (make -C frontend test)

# Wave 2 — backend core stream (serial) ∥ frontend stream (parallel):
backend-go:     T003 → T004(RED) → T005(fix) → T006 → T007 → T009 → T010 → T011 → T013
frontend-react: T008  ,  T012   (independent display tests)

# Wave 3 — polish (parallel):
backend-go:     T014 → T015 → T017
frontend-react: T016
```

---

## Implementation Strategy

### MVP First (User Story 1 only)

1. Phase 1 Setup → Phase 2 Foundational (T003).
2. US1: T004 (RED) → T005 (fix) → T006 (GREEN).
3. **STOP & VALIDATE**: corruption gone. Ship the MVP — this alone closes the reported bug.

### Incremental

- Add US2 (scientific-notation lock-in + frontend display test).
- Add US3 (explicit-error lock-in + frontend alert test).
- Polish: lint, coverage, quickstart curl.

---

## Notes

- The entire production change is **one function** (`formatResult`) + **one constant**. US2/US3 are predominantly test and contract verification — by design, the fix is surgical (plan.md, research.md).
- Verify each `⚠️ write-first` test fails (or locks intended behavior) before/with implementation.
- Commit after each story checkpoint.
