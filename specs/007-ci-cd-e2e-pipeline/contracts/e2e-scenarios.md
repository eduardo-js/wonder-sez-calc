# Contract: E2E Scenarios (v1)

Browser-driven scenarios run by Playwright against the running Docker Compose stack
(frontend `http://localhost:5173`, backend `http://localhost:8080`). Chromium only.

## Scenario 1 — Calculation round trip (P1 of US2)

**Given** the full stack is up and the backend reports healthy,
**When** the user opens the calculator, enters `2 + 2`, and activates equals,
**Then** the displayed result is `4`.

- Proves: UI renders, frontend reaches backend, backend computes, result returns and renders.
- This is the minimum bar required by FR-005 / SC-003.

### Implemented suite (`e2e/tests/calculator.spec.ts`)

The round trip above is the core bar; the suite extends it to cover the full
keypad surface and the error path. Each test drives the real browser → frontend
→ backend → display chain. Operator buttons carry the Unicode glyphs `÷ × − +`;
the in-progress hint line renders the logical tokens `/ * - +`.

| # | Scenario | Input | Expected |
|---|----------|-------|----------|
| 1 | Initial load | (none) | display `0` |
| 2 | Addition | `2 + 2 =` | `4` |
| 3 | Subtraction | `9 − 4 =` | `5` |
| 4 | Multiplication | `6 × 7 =` | `42` |
| 5 | Division | `10 ÷ 4 =` | `2.5` |
| 6 | Decimal operands | `1.5 + 2.25 =` | `3.75` |
| 7 | Chained operations | `2 + 3 =` then `× 4 =` | `20` (result reused) |
| 8 | Clear resets | `55` then `C` | `0` |
| 9 | Expression hint line | `12 +` then `3` | `role="status"` shows `12 +` → `12 + 3` |
| 10 | Division by zero | `5 ÷ 0 =` | `role="alert"` matches `/division by zero/i` |
| 11 | Error recovery | `5 ÷ 0 =` then `7` | alert clears, display `7` |
| 12 | Large-number round trip | `999999999 × 999999999 =` | scientific notation (contains `e`) |

## Scenario 2 — Failure produces diagnostics (edge case)

**Given** a run where a test assertion fails (or the stack misbehaves),
**When** the e2e job finishes,
**Then** a Playwright trace + screenshot and `docker compose logs` are available as artifacts.

- Proves: FR-007 / SC-004 (diagnose from artifacts alone).

## Assertions style

- Prefer user-visible assertions (text in the result display), not internal calls.
- Existing stable selectors to use (verified in `frontend/src/components/calculator/`):
  - result display: `[data-testid="calculator-display"]` (read `textContent`)
  - keypad container: `[data-testid="calculator-keypad"]`
  - buttons: `aria-label` (e.g. accessible name per digit/operator) → use `getByRole('button', { name })`
  - error state: `role="alert"`; in-progress: `role="status"` / `data-testid="spinner"`
- These resilient selectors mean styling changes won't break the e2e suite.

## Out of scope (v1)

- Multiple browsers/devices, visual regression, responsive/layout assertions.
- Expand in a later feature if needed.
