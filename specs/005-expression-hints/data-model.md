# Phase 1 Data Model: Expression Hints

The hint is **pure derived state** — no additions to `CalculatorState`, no new reducer
actions, and no backend/persistence model changes. (An earlier design that added
`hint`/`hintStatus` state and a debounced preview was removed: the frontend must not
evaluate; see research.md and spec FR-003.)

## No state changes

`CalculatorState` (`frontend/src/hooks/calculatorTypes.ts`) is unchanged:
`display`, `current`, `previous`, `operator`, `overwrite`, `error`, `status`, `errorMsg`.

## Derived hint expression

The hint text is computed (not stored) by a pure helper in
`frontend/src/lib/expression.ts`:

```ts
buildHintExpression(previous, operator, current, overwrite): string
```

Rules:

| State | Hint text |
|-------|-----------|
| `previous !== null && operator !== null && overwrite === true` (operator chosen, RHS not typed) | `` `${previous} ${operator}` `` (e.g. `12 +`) |
| `previous !== null && operator !== null && overwrite === false` (RHS being typed) | `` `${previous} ${operator} ${current}` `` (e.g. `12 + 3`) |
| `previous === null && overwrite === false` (typing the first number) | `current` (e.g. `12`) |
| otherwise (fresh / post-equals / post-operator overwrite with no operator) | `""` (hidden) |

Notes:
- After `CHOOSE_OP`, `current` still holds the **left** operand while `overwrite` is true,
  so the helper deliberately omits `current` and shows only `previous operator`.
- After `EVALUATE_SUCCESS` and `CLEAR`, `overwrite` is true and `previous`/`operator` are
  null → hint is `""` (cleared for the next entry).

## Invariants

- The hint performs **no** evaluation and triggers **no** request (FR-003, SC-002).
- The hint text matches the expression sent to the backend on equals (`buildExpression`)
  for the same operands/operator (FR-005, SC-001).
- The hint is suppressed whenever an error alert is shown (FR-006), handled in the display.
