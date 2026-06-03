# Data Model: Wire Frontend and Backend (Server-Side Calculation)

**Feature**: 003-wire-frontend-backend | **Date**: 2026-06-02

## Backend

### CalculationRequest (wire input)

| Field        | Type   | Rules                                                                 |
|--------------|--------|----------------------------------------------------------------------|
| `expression` | string | required; non-empty after trim; max 256 chars; allowed chars only: digits, `.`, `+ - * /`, `(`, `)`, spaces |

Validation tags (gin/validator): `binding:"required"`. Length/charset enforced in the service layer (returns `validation_failed` with `fields.expression`).

### CalculationResult (wire output, 200)

| Field        | Type   | Notes                                                        |
|--------------|--------|-------------------------------------------------------------|
| `result`     | string | formatted value (integers plain; else ≤12 significant digits) |
| `expression` | string | echo of the normalized/received expression                  |

### CalculationError (wire output, 4xx — existing envelope)

```json
{ "error": { "code": "...", "message": "...", "fields": { "expression": "..." } } }
```

| Code                 | HTTP | When                                                      |
|----------------------|------|----------------------------------------------------------|
| `validation_failed`  | 400  | empty, too long, illegal chars, malformed/unbalanced expr |
| `calculation_error`  | 422  | division by zero, non-finite (overflow/Inf/NaN)          |
| `bad_request`        | 400  | unparseable JSON body                                     |
| `internal`           | 500  | unexpected failure                                        |

### Internal evaluation (service layer, not on the wire)

- **Token**: `{kind: number|operator|lparen|rparen, value}`.
- **Evaluator**: tokenize → shunting-yard (RPN) → fold. Precedence `* /` > `+ -`; left-associative; unary minus at expr start or after `(`/operator.
- **Errors** (typed sentinels): `ErrInvalidExpression`, `ErrDivideByZero`, `ErrNonFinite` → mapped to wire codes by the handler.

## Frontend

### CalculatorState (extended)

Existing fields (`display`, `current`, `previous`, `operator`, `overwrite`, `error`) plus:

| Field      | Type                              | Notes                                              |
|------------|-----------------------------------|----------------------------------------------------|
| `status`   | `"idle" \| "loading" \| "error"`  | drives `=` spinner and request lock                |
| `errorMsg` | `string \| null`                  | backend error message for the alert region         |

### State transitions (request lifecycle)

```
idle --(press =, valid parts)--> loading        [issue POST /calculate; ignore further = while loading]
loading --(200)--> idle                          [display result]
loading --(4xx error)--> error                   [errorMsg = backend message; alert shown; display unchanged]
loading --(network fail/timeout)--> error        [errorMsg = generic connectivity message]
error --(any input / CLEAR)--> idle              [clear errorMsg]
```

- Side effects (fetch, AbortController, timeout) live in the hook/effect; the reducer stays pure.
- Client-side `compute()` is **removed** — no arithmetic on the frontend (FR-001).
