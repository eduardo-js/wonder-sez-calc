# Feature Specification: Expression Hints & UI Polish

**Feature Branch**: `005-expression-hints`

**Created**: 2026-06-02

**Status**: Draft

**Input**: User description: "let's improve frontend UI / lets add hints for current expression being evalued" — clarified: the hint shows the **in-progress expression text only**; the frontend does **not** evaluate it (evaluation is the backend's job, on equals).

## User Scenarios & Testing *(mandatory)*

### User Story 1 - See the expression being built (Priority: P1)

As a person entering a calculation, I want a secondary line that shows the expression
I am building — the current number, then the running expression (e.g. `12 +`, then
`12 + 3`) — so I can see the full context of what I am about to evaluate, not just the
last operand on the main line.

**Why this priority**: This is the core of the request. It delivers immediate,
standalone value: the user always sees what expression will be sent on equals. Shipping
only this story already improves the product. It does **not** evaluate anything — the
result still appears only after pressing equals.

**Independent Test**: Press `1 2 + 3`; confirm a secondary line shows `12 + 3` (updating
through `12`, `12 +`, `12 + 3`) while the main line shows the current operand, and that
no result is computed until equals is pressed.

**Acceptance Scenarios**:

1. **Given** an empty calculator, **When** the user types a number, **Then** the hint line shows that number as it is typed.
2. **Given** a number has been entered, **When** the user chooses an operator, **Then** the hint line shows `number operator` (e.g. `12 +`).
3. **Given** an operator has been chosen, **When** the user types the second operand, **Then** the hint line shows the full expression (e.g. `12 + 3`) — still without a result.
4. **Given** a visible expression hint, **When** the user presses equals, **Then** the backend result appears on the main line and the hint resets for the next entry.

---

### User Story 2 - Graceful, non-distracting hint behavior (Priority: P2)

As a person entering a calculation, I want the expression hint to stay out of the way
when it has nothing useful to show — hidden on a fresh entry, cleared after equals, and
not competing with error messages — so it informs without adding noise.

**Why this priority**: Edge behaviors keep the hint trustworthy. The feature still
delivers value (Story 1) without this polish, but these rules prevent a redundant or
confusing hint.

**Independent Test**: On a fresh calculator the hint line is empty; after an error the
error message is shown without a competing hint; after equals the hint resets.

**Acceptance Scenarios**:

1. **Given** a fresh calculator (nothing typed yet), **When** the user has not entered a number, **Then** no hint line is shown.
2. **Given** the backend returns an error on equals, **When** the error is displayed, **Then** the hint line does not compete with or obscure the error.
3. **Given** the user pressed equals, **When** the result is shown, **Then** the hint line is cleared until new input begins.

---

### User Story 3 - Clearer, more polished calculator interface (Priority: P3)

As a user, I want the calculator to look and feel polished — readable display, clear
separation between the expression hint and the entry/result line, comfortable button
targets, and visible feedback on interaction (including keyboard focus) — so it is
pleasant and easy to use.

**Why this priority**: General UI improvement requested alongside hints. It raises
overall quality but is not required for the hint to work.

**Independent Test**: Tab through the keypad and confirm each button shows a visible
focus ring; confirm the expression hint line is visually distinct (smaller, muted,
positioned above) from the main entry/result line on desktop and a narrow viewport.

**Acceptance Scenarios**:

1. **Given** the calculator is open, **When** the user views it, **Then** the expression hint and the main entry/result line are visually distinguished (size, weight, position) so neither is mistaken for the other.
2. **Given** a long expression that exceeds the display width, **When** it is entered, **Then** the display handles overflow gracefully (truncation) without breaking layout.
3. **Given** any interactive control, **When** the user hovers, presses, or focuses it (including via keyboard), **Then** the control gives visible feedback and remains operable.

---

### Edge Cases

- Long expressions: the hint line truncates without breaking layout.
- Fresh/initial state (`0`, nothing typed): no redundant hint line is shown.
- After an operator but before the second operand: the hint shows `number operator` (e.g. `12 +`), never a stale duplicated operand.
- Error on equals: the error message takes precedence; the hint is suppressed.
- After equals: the hint clears until the next input begins.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The system MUST display a secondary hint line, separate from the main entry/result line, showing the **in-progress expression text** as the user types.
- **FR-002**: The hint MUST update as the entry changes: the number being typed, then `number operator`, then `number operator number`.
- **FR-003**: The system MUST NOT evaluate the expression for the hint; no result is shown in the hint, and no calculation is performed until the user presses equals.
- **FR-004**: When the user presses equals, the backend-computed result MUST appear on the main line and the hint MUST reset for the next entry.
- **FR-005**: The hint expression text MUST be consistent with the expression actually sent to the backend on equals (same operands and operator).
- **FR-006**: The hint MUST be hidden when there is nothing meaningful to show (fresh/empty entry) and MUST be suppressed while an error message is displayed.
- **FR-007**: The hint line and the main entry/result line MUST be visually distinguishable (size, weight, position).
- **FR-008**: The display MUST handle expressions and results that exceed the visible width without breaking layout (truncation).
- **FR-009**: Interactive controls MUST provide visible feedback for hover, press, and keyboard-focus states and remain fully operable via keyboard.

### Key Entities

- **Current Expression**: The in-progress input (left operand, operator, right operand) the user is building.
- **Expression Hint**: A derived, read-only text rendering of the current expression — never a computed result.
- **Committed Result**: The value produced by the backend when the user presses equals; shown on the main line.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: While typing, the hint text always equals the expression that would be sent to the backend on equals (100% consistency for the same input).
- **SC-002**: No calculation is performed before equals (0 evaluation requests are triggered by typing digits or operators).
- **SC-003**: The hint reflects the current entry immediately (the next keystroke never shows a stale expression).
- **SC-004**: Users correctly identify the hint as "the expression I'm building" vs the main entry/result line in at least 90% of usability checks.
- **SC-005**: The layout remains intact (no overflow, clipping, or broken controls) across desktop and narrow-mobile viewports for long expressions.

## Assumptions

- Evaluation is the backend's responsibility and happens only on equals (`POST /api/v1/calculate`); the frontend never evaluates for the hint.
- The hint is read-only derived state from the current operands/operator; it never triggers a request and never auto-commits.
- "Improve frontend UI" is scoped to the calculator screen (display, hint line, keypad, interactive feedback), not a full redesign or theming system.
- The calculator is a binary-operation entry model (`previous operator current`); the hint reflects that model.
- Accessibility basics (keyboard operability, visible focus, readable contrast) are in scope; a full WCAG audit is out of scope for this iteration.
