# Frontend UI Contract: Calculator Components

**Feature**: 001-calculator-ui-barebones | **Date**: 2026-06-02

The UI contract = the public props/behavior each reusable component exposes. Tests assert these.

## `<CalculatorButton>`

Reusable button wrapping the shadcn `Button` primitive.

**Props**
| Prop      | Type                                   | Required | Notes |
|-----------|----------------------------------------|----------|-------|
| `label`   | string                                 | yes      | Rendered text content |
| `value`   | string                                 | yes      | Token passed to `onPress` |
| `variant` | `'number' \| 'operator' \| 'action'`   | no       | Defaults to `'number'`; maps to shadcn variant/classes |
| `onPress` | `(value: string) => void`              | yes      | Fired on click/tap |
| `ariaLabel` | string                               | no       | Falls back to `label` |

**Behavior contract (tested)**
- Renders an element with role `button` whose accessible name = `ariaLabel ?? label`.
- Calls `onPress(value)` exactly once per activation.
- Applies the variant's styling class so number/operator/action are visually distinct.

## `<CalculatorDisplay>`

**Props**
| Prop      | Type    | Required | Notes |
|-----------|---------|----------|-------|
| `value`   | string  | yes      | Current display string |
| `error`   | string \| null | no | When set, shows error state styling/text |

**Behavior contract (tested)**
- Renders `value` as visible text.
- When `error` is set, exposes an error indication (e.g. `role="alert"` / distinct text).

## `<CalculatorKeypad>`

**Props**
| Prop      | Type                       | Required | Notes |
|-----------|----------------------------|----------|-------|
| `onPress` | `(value: string) => void`  | yes      | Forwarded to each button |

**Behavior contract (tested)**
- Renders one `<CalculatorButton>` per entry in the keypad config array.
- Uses a CSS Grid that reflows responsively (no overlap/clipping at sm/md/lg).

## `<Calculator>` (container)

**Behavior contract (tested — integration)**
- Wires `useCalculator` reducer to `CalculatorDisplay` + `CalculatorKeypad`.
- Pressing `1 2 + 7 =` shows `19` (FR-003).
- Pressing `.` twice in one operand adds a single decimal (FR-004).
- Dividing by zero shows the error state and recovers on `C` (FR-005/FR-006).
