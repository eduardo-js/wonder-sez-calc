# Research: Wire Frontend and Backend (Server-Side Calculation)

**Feature**: 003-wire-frontend-backend | **Date**: 2026-06-02

## R1. Request shape: expression string vs. operand triple

- **Decision**: Frontend POSTs a single expression **string** (`{"expression": "..."}`); backend tokenizes and evaluates it.
- **Rationale**: Spec FR-004 and acceptance scenario 1.2 require standard operator precedence (`2 + 3 * 4 → 14`). A general expression parser satisfies this and any future multi-operand UI. The current UI only forms `a op b`, which is still a valid expression string — no contract change needed when the UI grows.
- **Alternatives considered**:
  - `{left, operator, right}` triple — rejected: cannot express precedence/parentheses, leaks the current UI's single-binary-op limitation into the contract (violates Contract-Driven Integration).

## R2. Backend evaluation algorithm

- **Decision**: Tokenizer + shunting-yard to RPN, evaluated with `*math/big*`-free `float64`. Supports `+ - * /`, parentheses, decimals, leading unary minus.
- **Rationale**: Shunting-yard is small, well-understood, table-testable (constitution II), and handles precedence/associativity without a heavyweight dependency (constitution: pin/justify deps — none added).
- **Validation/error rules**:
  - Empty/whitespace, unknown chars, unbalanced parens, malformed sequences (`5 + * 2`, trailing operator) → `validation_failed`.
  - Division by zero, non-finite result (overflow/`Inf`/`NaN`) → `calculation_error`.
- **Alternatives considered**: `go/parser`+`go/types` (overkill, allows non-math Go expressions — injection surface); third-party expr libs (unjustified dependency).

## R3. Result formatting

- **Decision**: Backend returns `result` as a **string**, formatted to drop float noise (parity with prior client formatting: integers plain, else ~12 significant digits).
- **Rationale**: Single source of formatting truth on the backend; frontend only displays. Avoids JS/Go float-rendering divergence.
- **Alternatives considered**: Return raw JSON number — rejected: pushes formatting back to the client, risks divergence.

## R4. Error envelope & codes

- **Decision**: Reuse the existing `httpx` error envelope (`{"error":{code,message,fields?}}`). Add one new code `calculation_error` for math-domain failures; keep `validation_failed` (with `fields.expression`) for malformed input.
- **Rationale**: One envelope keeps the contract consistent (002 already established it). Distinct code lets the frontend phrase math errors vs. input errors. HTTP: `400` for validation, `422 Unprocessable Entity` for `calculation_error`.
- **Alternatives considered**: New bespoke error shape — rejected: breaks contract consistency.

## R5. Frontend request lifecycle (lock + spinner)

- **Decision**: Add `status: "idle" | "loading" | "error"` + `requesting` guard to calculator state. `EVALUATE` triggers an async `calculate()` API call via a typed client (`src/lib/api.ts`); while `loading`, the `=` button renders a spinner and is disabled; duplicate presses are ignored.
- **Rationale**: Satisfies FR-009/010/011 with a single in-flight guard; reducer stays pure, side-effect (fetch) lives in the hook/effect.
- **Request control**: `AbortController` with a client timeout (default 5s). On network failure/timeout → alert + restore control (FR-012).
- **Alternatives considered**: Disable whole keypad during request — rejected: spec only requires blocking *calculation* requests; over-locking hurts UX.

## R6. Error surfacing ("alerts")

- **Decision**: Non-blocking, accessible alert region (`role="alert"`) showing the backend `error.message`; transient, dismiss on next input/clear.
- **Rationale**: Accessible, testable, and matches "display alerts for errors" without `window.alert` (untestable, blocking).
- **Alternatives considered**: `window.alert()` — rejected: blocks the event loop, untestable, poor UX.

## R7. API base URL / config

- **Decision**: Frontend reads base URL from `import.meta.env.VITE_API_BASE_URL` (default `http://localhost:8080`). Backend CORS already allows `http://localhost:5173` (002).
- **Rationale**: Keeps environments configurable; no hardcoded host (constitution V/operability).

## Remaining NEEDS CLARIFICATION

None. All spec ambiguities resolved via documented assumptions in spec.md and decisions above.
