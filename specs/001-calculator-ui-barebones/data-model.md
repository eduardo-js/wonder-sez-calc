# Data Model: Calculator UI Barebones

**Feature**: 001-calculator-ui-barebones | **Date**: 2026-06-02

The app is client-side in v1; these are **frontend in-memory models**, not persisted entities.

## Entity: CalculatorButton (config)

Describes one reusable button. Drives rendering of the keypad from a single data array.

| Field    | Type                                              | Notes |
|----------|---------------------------------------------------|-------|
| `id`     | string                                            | Stable key for rendering/tests |
| `label`  | string                                            | Visible text (`7`, `+`, `=`, `C`) |
| `value`  | string                                            | Logical token emitted on press |
| `kind`   | `'digit' \| 'operator' \| 'decimal' \| 'equals' \| 'clear'` | Drives behavior + variant |
| `variant`| `'number' \| 'operator' \| 'action'`              | Maps to shadcn Button styling |
| `span`   | number (optional, default 1)                      | Grid column span (e.g. `0` spans 2) |

**Validation/rules**: `kind` determines which reducer action fires; `variant` is presentation
only. One config array is the single source of truth for the keypad (supports FR-002/FR-008).

## Entity: CalculatorState (reducer state)

The complete in-progress calculator state.

| Field        | Type                          | Notes |
|--------------|-------------------------------|-------|
| `display`    | string                        | What the user sees; defaults to `"0"` |
| `current`    | string                        | Operand currently being typed |
| `previous`   | string \| null                | Stored left-hand operand |
| `operator`   | `'+' \| '-' \| '*' \| '/'` \| null | Pending operation |
| `overwrite`  | boolean                       | If true, next digit replaces display (after `=` or operator) |
| `error`      | string \| null                | Set on invalid op (e.g. divide-by-zero); blocks compute |

### State transitions (actions)

| Action          | Trigger        | Effect / Validation |
|-----------------|----------------|---------------------|
| `INPUT_DIGIT`   | digit button   | Append digit; if `overwrite`, start fresh. Normalize leading zeros (FR-004) |
| `INPUT_DECIMAL` | `.`            | Append `.` only if `current` has none (FR-004) |
| `CHOOSE_OP`     | operator       | Commit `current`→`previous`, store operator; ignore/replace if leading or repeated operator (FR-004) |
| `EVALUATE`      | `=`            | Compute `previous (op) current`; divide-by-zero → set `error` (FR-005). Result becomes `current`, `overwrite=true` |
| `CLEAR`         | `C`            | Reset to initial state, clear `error` (FR-006) |

**Initial state**: `{ display: "0", current: "0", previous: null, operator: null, overwrite: true, error: null }`

**Error recovery**: Any `CLEAR` from an `error` state returns to initial (FR-005 acceptance #2/#4).

## Backend (scaffolding only)

No domain entities this feature. The backend exposes operational endpoints only:

| Endpoint   | Response (JSON)                       | Purpose |
|------------|---------------------------------------|---------|
| `/healthz` | `{ "status": "ok" }`                  | Liveness |
| `/readyz`  | `{ "status": "ready" }`               | Readiness |

The frontend↔backend calculation contract is intentionally deferred to a later feature.
