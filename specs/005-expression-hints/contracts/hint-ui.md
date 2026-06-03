# UI Contract: Expression Hint

This feature exposes **no new network interface** and makes **no change** to
`POST /api/v1/calculate` (still used only on equals). The hint never calls the backend.
The contract below is the internal UI component contract.

## `CalculatorDisplay` props (extended)

```ts
export interface CalculatorDisplayProps {
  value: string;                 // existing — primary entry/result line
  error: string | null;          // existing — local validation error (alert)
  errorMsg?: string | null;      // existing — backend/network error (alert)

  // NEW (hint):
  expression?: string;           // in-progress expression text, e.g. "12 +" or "12 + 3"; "" / omitted => hide
}
```

### Rendering rules

| Condition | Render |
|-----------|--------|
| `error`/`errorMsg` set | Existing alert path (unchanged); expression hint suppressed. |
| `expression` is `""`, null, or omitted | No hint line. |
| `expression` non-empty (and no error) | Muted secondary line **above** the main value showing `expression`. |

### Accessibility / hierarchy

- Hint line is visually secondary: smaller, muted text (e.g. `text-sm text-slate-400`),
  positioned above the main `value`, right-aligned, truncating on overflow. (FR-007, FR-008)
- Hint line uses `role="status"` + `aria-live="polite"` (must not preempt the assertive
  committed-error `role="alert"`).
- Main `value` remains the dominant element; its sizing/truncation behavior is unchanged.

## No change to `calculate`

`CalculateOptions` is unchanged (timeout only). The hint does not call `calculate`; the
expression is evaluated by the backend **only** when the user presses equals (existing
EVALUATE flow).

## Behavioral guarantees

- **No evaluation while typing**: the hint is derived text; no request is triggered by
  digits or operators. (FR-003, SC-002)
- **Consistency**: the hint text equals the expression sent on equals for the same
  operands/operator. (FR-005, SC-001)
- **Non-competing**: the hint is suppressed while an error alert is shown. (FR-006)
