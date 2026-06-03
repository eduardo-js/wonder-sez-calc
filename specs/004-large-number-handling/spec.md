# Feature Specification: Large-Number & Edge-Case Handling

**Feature Branch**: `004-large-number-handling`

**Created**: 2026-06-02

**Status**: Draft

**Input**: User description: "handle more edge cases on backend and frontend. Example: `-9223372036854775808 - 999999999999999999999999999999999999999999` returns `-9223372036854775808` (wrong). How to handle this properly? Add 'e' (scientific notation) special handling for larger scenarios."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Large-magnitude results are never silently wrong (Priority: P1)

A user enters an expression whose true result is larger (in magnitude) than the largest plain whole number the calculator can represent exactly — for example subtracting a 42-digit number from a large negative number. Today the calculator returns a plausible-looking but **incorrect** number. The user must instead receive either the mathematically correct value (rendered in scientific notation) or a clear, explicit "too large" error — never a wrong number presented as if correct.

**Why this priority**: This is a correctness defect. A calculator that confidently returns wrong answers is worse than one that errors out; correctness is the product. Fixing this is the core of the feature.

**Independent Test**: Submit the reported expression `-9223372036854775808 - 999999999999999999999999999999999999999999` and confirm the response is either the correct value in scientific notation or an explicit error — and is in no case the corrupted value `-9223372036854775808`.

**Acceptance Scenarios**:

1. **Given** an expression whose exact whole-number result exceeds the exact-integer range, **When** it is evaluated, **Then** the result is returned in scientific notation carrying the supported number of significant digits, and never a truncated/wrapped integer.
2. **Given** the originally reported expression, **When** it is evaluated, **Then** the returned result is NOT `-9223372036854775808`.
3. **Given** any valid expression, **When** it is evaluated, **Then** the returned result is mathematically correct within the documented precision, OR an explicit error is returned — there is no third "silently wrong" outcome.

---

### User Story 2 - Readable scientific notation for very large and very small results (Priority: P2)

A user performs calculations that produce very large or very small (near-zero but non-zero) results. The calculator should present these in a compact, conventional scientific-notation form (e.g., `1.23e42`, `4.5e-18`) consistently in the API response and in the on-screen display, without breaking the UI layout.

**Why this priority**: Once correctness is guaranteed (P1), the result must also be human-readable. Long digit strings are unreadable and can overflow the display; scientific notation is the conventional fix.

**Independent Test**: Evaluate an expression producing a very large result and one producing a very small non-zero result; confirm both render in scientific notation in the API payload and appear fully and legibly in the UI.

**Acceptance Scenarios**:

1. **Given** a result whose magnitude is above the plain-integer threshold, **When** displayed, **Then** it uses scientific notation with at most the documented number of significant digits and no trailing-zero noise.
2. **Given** a very small non-zero fractional result, **When** displayed, **Then** it uses scientific notation rather than rounding to `0`.
3. **Given** a scientific-notation result, **When** shown in the UI, **Then** it is fully visible and does not break or overflow the calculator layout.

---

### User Story 3 - Explicit errors for truly unrepresentable results (Priority: P3)

A user enters an expression whose result is so large it cannot be represented at all (overflows to infinity) or is otherwise undefined (e.g., division by zero). The calculator returns a clear, user-facing error message and the UI surfaces it gracefully instead of showing a broken or blank result.

**Why this priority**: A small but real class of inputs cannot be represented even in scientific notation. These must fail loudly and clearly. Lower priority because the existing error path already partially covers it; this story ensures consistency and good messaging.

**Independent Test**: Evaluate an expression that overflows to infinity and one that divides by zero; confirm both return an explicit error that the UI displays clearly.

**Acceptance Scenarios**:

1. **Given** an expression whose result overflows the representable range, **When** evaluated, **Then** an explicit "result too large / not finite" error is returned and shown to the user.
2. **Given** an error response from evaluation, **When** received by the UI, **Then** the error message is displayed clearly and the prior valid display is not corrupted.

---

### Edge Cases

- Result exactly at the largest exact-integer boundary → rendered as a plain whole number (not scientific notation).
- Result one step beyond the exact-integer boundary → rendered in scientific notation.
- Large **negative** results → handled identically to large positive results (sign preserved, no wrap-around).
- Very small non-zero results (e.g., `1 / 1e20`) → scientific notation, not collapsed to `0`.
- Input literals already in scientific notation (e.g., `1e40 * 2`) → accepted and evaluated.
- Input literal so large it overflows on parse (e.g., `1e400`) → explicit error, not a wrong value.
- Chained operations that accumulate magnitude across several steps → correctness preserved at each step.
- Whole-number results within range → continue to render plain (no regression to existing behavior).

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST NOT return a numerically incorrect result for any valid expression. If the exact result cannot be represented, the system MUST return either a correct scientific-notation approximation within the documented precision, or an explicit error — never a wrong finite number.
- **FR-002**: For whole-number results whose magnitude exceeds the largest value representable **exactly** as a plain integer, the system MUST format the result in scientific notation rather than as a plain integer.
- **FR-003**: Whole-number results within the exact-integer range MUST continue to render as plain integers with no decimal point (no regression).
- **FR-004**: Scientific-notation results MUST carry at most the documented number of significant digits, with insignificant trailing zeros removed.
- **FR-005**: Very small non-zero results MUST render in scientific notation instead of being rounded to `0`.
- **FR-006**: System MUST accept input number literals written in scientific notation and evaluate them.
- **FR-007**: When a result cannot be represented at all (overflows to infinity) or is undefined, the system MUST return an explicit, user-facing error, not a placeholder or wrong value.
- **FR-008**: The on-screen display MUST render scientific-notation results fully and legibly without breaking or overflowing the calculator layout.
- **FR-009**: The on-screen display MUST surface evaluation errors clearly without corrupting the previously displayed valid value.
- **FR-010**: Result formatting MUST be consistent between the API response and what the user sees on screen (the display does not re-corrupt or re-round a correctly formatted result).

### Key Entities *(include if feature involves data)*

- **Calculation result**: The numeric outcome of an expression, plus the formatting decision (plain integer vs. scientific notation) and the precision applied.
- **Evaluation error**: An explicit failure outcome (too large / not finite / undefined) carrying a user-facing message.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: The originally reported expression no longer returns `-9223372036854775808`; it returns the correct value in scientific notation (or an explicit error).
- **SC-002**: 100% of valid expressions return a result that is either mathematically correct within the documented precision or an explicit error — zero silently-incorrect results across the regression test suite.
- **SC-003**: Every large-magnitude and very-small-magnitude result renders in scientific notation and is fully visible in the UI at supported viewport sizes (no clipping or layout break).
- **SC-004**: Whole-number results within the exact-integer range render identically to today (no regression), verified by existing test cases continuing to pass.
- **SC-005**: Truly unrepresentable inputs (overflow to infinity, division by zero) produce a clear error message visible to the user 100% of the time.

## Assumptions

- "Exact-integer range" means the largest magnitude at which consecutive whole numbers can each be represented exactly (≈ 9.0×10¹⁵, i.e. 2⁵³). Above this, plain-integer rendering is misleading, so scientific notation is used. Whole numbers between this threshold and the current integer-cast limit are also affected and are covered by the same rule.
- The supported significant-digit precision for scientific notation matches the existing non-integer formatting precision (12 significant digits) for consistency; no new precision target is introduced.
- The calculator continues to operate on standard double-precision floating-point magnitudes; results beyond that range are reported as errors rather than supporting arbitrary-precision big-number arithmetic (out of scope for this feature).
- Arbitrary-precision / exact big-integer arithmetic is explicitly out of scope; this feature corrects formatting/representation and error handling, not the underlying numeric domain.
- The existing API error envelope and result schema are reused; no breaking contract change is required (scientific-notation strings already satisfy the existing string result field).
