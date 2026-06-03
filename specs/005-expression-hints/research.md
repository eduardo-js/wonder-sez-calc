# Phase 0 Research: Expression Hints & UI Polish

> **Design correction (2026-06-02):** An earlier version of this feature proposed a live
> *result preview* (debounced backend `calculate()` calls showing `12 + 3 = 15`). That was
> rejected by the user: the hint must show the **in-progress expression text only**; the
> frontend must **not** evaluate. Evaluation is the backend's responsibility and runs only
> on equals. The decisions below reflect the corrected design.

## R1 — What the hint shows

**Decision**: The hint renders the **in-progress expression text** (`12`, then `12 +`,
then `12 + 3`) — never a computed result.

**Rationale**: Matches the user's explicit intent and keeps evaluation a backend
responsibility (separation of concerns). The result appears only on equals.

**Alternatives considered**:
- *Live result preview via debounced backend calls* — rejected by the user; adds request
  load and blurs the frontend/backend boundary.
- *Local JS evaluation for a preview* — rejected: there is no local evaluator and
  duplicating backend semantics risks divergence.

## R2 — How the hint is produced

**Decision**: A pure helper `buildHintExpression(previous, operator, current, overwrite)`
in `frontend/src/lib/expression.ts` derives the text from existing state. No new state,
no reducer actions, no effects, no async.

**Rationale**: The hint is a pure function of current operands/operator. Keeping it
derived (computed in the component render path) is the simplest correct design and keeps
the reducer untouched (constitution II/III).

**Alternatives considered**:
- *Store `hint`/`hintStatus` in state + dispatch on change* — rejected: unnecessary state
  and lifecycle for a pure derivation.

## R3 — Operator/overwrite handling

**Decision**: After `CHOOSE_OP`, `current` still holds the left operand while `overwrite`
is true, so the helper shows `previous operator` only; once the RHS is typed
(`overwrite === false`) it appends `current`. Fresh/post-equals state → `""` (hidden).

**Rationale**: Prevents a stale duplicated operand (e.g. `12 + 12`) and avoids a redundant
hint on the initial `0`.

## R4 — Presentation & accessibility

**Decision**: Render the hint as a smaller, muted secondary line above the main value with
`role="status"` + `aria-live="polite"`; suppress it while an error alert is shown. Add
`focus-visible` rings to keypad buttons for keyboard a11y.

**Rationale**: `polite` avoids fighting the assertive committed-error alert. Distinct
size/position satisfies FR-007. Focus rings satisfy FR-009.

## R5 — No backend or contract change

**Decision**: `POST /api/v1/calculate` and `CalculateOptions` are unchanged; the hint never
calls the backend. Evaluation remains on the existing equals (EVALUATE) path.

**Rationale**: The hint is display-only; no integration surface changes (constitution IV).
