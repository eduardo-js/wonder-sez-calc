# Quickstart: Expression Hints & UI Polish

Frontend-only feature. Reuses the existing backend `calculate` endpoint.

## Prerequisites

- Backend running at `:8080` (so the hint can evaluate): `make backend-run` (or per repo scripts).
- Node 24 via `.nvmrc`.

## Run

```bash
cd frontend
npm install        # if needed
npm run dev        # Vite dev server at :5173
```

Open the app, then:

1. Type `12`, press `+`, type `3` → a muted secondary line shows `12 + 3 = 15` **before**
   pressing `=` (US1).
2. Press `+` and stop (trailing operator) → hint shows a muted "incomplete" indicator, no
   wrong value (US2).
3. Type `5`, `/`, `0` → hint shows a plain-language "can't evaluate yet" (divide-by-zero,
   US2).
4. Press `=` → the previewed result becomes the committed result; hint resets (US1 #2).
5. Tab through the keypad → each button shows a visible focus ring; entry vs hint lines are
   visually distinct (US3).

## Test (TDD — write/expect these first, confirm red, then implement)

```bash
cd frontend
npm test                                   # full Vitest run
npm test -- useCalculator                  # hook hint lifecycle: pending→loading→ready/invalid, debounce, stale-cancel
npm test -- CalculatorDisplay              # hint-line render states + a11y (aria-live polite)
npm test -- Calculator.test                # integration: preview appears before equals; commit matches preview
```

### Key test cases

- **Hook**: typing a full binary op dispatches a debounced `calculate`; `HINT_READY` sets
  `hint`; backend `calculation_error` → `HINT_INVALID`; a second input before the first
  resolves discards the stale result (generation guard / abort).
- **Display**: `hintStatus` × `hintExpr`/`hint` → correct line text; error alert suppresses
  hint; hint uses `aria-live="polite"`.
- **Integration**: preview value for `12 + 3` equals the committed value after `=`
  (FR-005 / SC-001), using a mocked `calculate`.

## Verify

```bash
cd frontend
npm run lint
npm run build        # type-check (strict) + bundle
```

All green + the manual steps above demonstrate US1–US3.
</content>
