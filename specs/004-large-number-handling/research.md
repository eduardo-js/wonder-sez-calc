# Phase 0 Research: Large-Number & Edge-Case Handling

## Decision 1 — Root cause of the corruption

- **Decision**: The defect is in `formatResult` (`backend/internal/calc/evaluator.go:322`), not in evaluation. The integer-valued branch does `strconv.FormatInt(int64(v), 10)`. Converting a `float64` whose magnitude exceeds the int64 range to `int64` is undefined/saturating in Go and yields `math.MinInt64` (`-9223372036854775808`). Evaluation itself already produced the correct `float64` (`≈ -1e42`).
- **Rationale**: `evalRPN` returns a finite `float64`; the value is only mangled at formatting time. Confirmed by the reported payload: the result equals exactly `MinInt64`.
- **Alternatives considered**: Arbitrary-precision arithmetic (`math/big`) — rejected; out of scope (spec assumption), large dependency-free but large behavioral change, and the contract is a double-precision calculator.

## Decision 2 — Formatting threshold

- **Decision**: Render a plain integer only when `v == math.Trunc(v)`, `v` is finite, and `math.Abs(v) <= 2^53` (`9007199254740992`). Otherwise format with `strconv.FormatFloat(v, 'g', 12, 64)` (scientific notation for large/small magnitudes).
- **Rationale**: 2⁵³ is the largest magnitude at which every consecutive integer is exactly representable in `float64`. Below it, plain digits are meaningful and exact; above it, plain digits imply false unit-precision, and beyond int64 they corrupt. A single threshold both removes the misleading output and eliminates the int64-overflow bug (1e42 ≫ 2⁵³). Matches the spec assumption.
- **Alternatives considered**:
  - Threshold at int64 max (2⁶³): fixes only the crash, still prints misleading exact-looking digits for 2⁵³–2⁶³. Rejected.
  - `math/big.Int` for exact large integers: out of scope (Decision 1).

## Decision 3 — Scientific-notation representation

- **Decision**: Use Go's `%g` with 12 significant digits (`strconv.FormatFloat(v, 'g', 12, 64)`), which emits `e+NN` / `e-NN` form (e.g., `1e+42`, `1.23457e+10`, `1e-20`) and strips insignificant trailing zeros automatically.
- **Rationale**: Reuses the precision already used for non-integer results (consistency, FR-004), no new formatting code, no new dependency. The `e+NN` form is conventional scientific notation and the frontend renders it verbatim.
- **Alternatives considered**: Normalizing to bare `e42` (strip `+`/leading zeros) — cosmetic only; rejected to keep the change minimal and avoid hand-rolled formatting. Can revisit if product wants a specific display style.

## Decision 4 — Very small results (already correct)

- **Decision**: No code change needed for very small non-zero results; they already flow through the `%g` branch (e.g., `1/1e20` → `1e-20`) because they are not integer-valued.
- **Rationale**: Confirmed by reading the current branch logic. Covered with an explicit regression test to lock the behavior (FR-005).

## Decision 5 — Truly unrepresentable results (no change)

- **Decision**: Keep existing behavior: ±Inf / NaN → `ErrNonFinite`; divide-by-zero → `ErrDivideByZero`; over-large literals (e.g., `1e400`) parse to ±Inf and surface as `ErrNonFinite`. These map to the `calculation_error` wire code (HTTP 422).
- **Rationale**: Satisfies FR-007 already; the feature only needs to ensure these stay explicit and are covered by tests.
- **Alternatives considered**: A distinct "too large" error code — rejected; adds a contract change for no user benefit over the existing clear message.

## Decision 6 — Frontend scope

- **Decision**: No production frontend change. The display renders the backend `result` string verbatim and already truncates + auto-sizes by length; short `e`-notation strings fit. Add a Vitest case confirming an `e`-notation result and an error string render fully and legibly (FR-008, FR-009).
- **Rationale**: The keypad cannot emit `e`, and all arithmetic is server-side (feature 003), so there is no client formatting to correct — only verification.
- **Alternatives considered**: Client-side reformatting of scientific notation — rejected; would duplicate/contradict the backend source of truth (FR-010, constitution IV).
