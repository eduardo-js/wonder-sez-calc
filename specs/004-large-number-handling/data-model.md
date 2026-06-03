# Phase 1 Data Model: Large-Number & Edge-Case Handling

No persistent storage. The only "entities" are the in-flight calculation value and its formatting decision.

## Entity: Calculation result value

| Field | Type | Notes |
|-------|------|-------|
| value | float64 | Finite numeric result from `evalRPN`. |
| isInteger | bool | `value == math.Trunc(value)` and finite. |
| magnitude | float64 | `math.Abs(value)`. |

**Formatting rule (state → output string):**

| Condition | Output form | Example |
|-----------|-------------|---------|
| `isInteger && magnitude <= 2^53` | Plain integer | `42`, `-9007199254740992` |
| `isInteger && magnitude > 2^53` | Scientific notation (`%g`, ≤12 sig) | `1e+42` |
| `!isInteger` (incl. very small) | Scientific/decimal (`%g`, ≤12 sig) | `3.14159265359`, `1e-20` |
| `±Inf` or `NaN` | — (not formatted) | → evaluation error |

Constant: `maxExactInt = 1 << 53` (`9007199254740992`).

## Entity: Evaluation error (unchanged)

| Sentinel | Meaning | Wire code | HTTP |
|----------|---------|-----------|------|
| `ErrInvalidExpression` | Malformed input | `validation_failed` | 400 |
| `ErrDivideByZero` | Division by zero | `calculation_error` | 422 |
| `ErrNonFinite` | Result ±Inf/NaN (incl. overflow, over-large literal) | `calculation_error` | 422 |

## Validation rules (from requirements)

- FR-001/FR-002/FR-003: integer-vs-scientific decision is gated on `magnitude <= 2^53`; never cast out-of-range float to int64.
- FR-004: scientific notation carries ≤12 significant digits, trailing zeros stripped.
- FR-005: very small non-zero values format via `%g`, never collapse to `0`.
- FR-007: ±Inf/NaN remain explicit errors, never formatted as a value.
