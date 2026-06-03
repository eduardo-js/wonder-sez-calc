# Feature Specification: Wire Frontend and Backend (Server-Side Calculation)

**Feature Branch**: `003-wire-frontend-backend`

**Created**: 2026-06-02

**Status**: Draft

**Input**: User description: "let's wire the frontend and backend; 1. frontend should not make any calculations, these should be handled by backend 2. backend should receive user input, validate it and handle it correctly 2.1 if its a valid mathematical equation, calculate it and return as json 2.2 if its invalid input return an error with details 2.2.1 frontend should display alerts for errors 2.2.2 block multiple requests from frontend, add a spinner in the '=' sign while backend is processing information"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Compute a valid expression via the backend (Priority: P1)

A user enters a mathematical expression on the calculator and presses the equals control. The frontend sends the expression to the backend, which validates and evaluates it, and the frontend displays the returned result. The frontend performs no arithmetic itself.

**Why this priority**: This is the core of the feature — moving all calculation authority to the backend. Without it, nothing else is meaningful. It alone delivers a working server-computed calculator.

**Independent Test**: Enter a valid expression (e.g., `12 + 7 * 3`), press equals, and confirm the displayed result matches the correct value and that the result originated from the backend (no local computation path is exercised).

**Acceptance Scenarios**:

1. **Given** a valid expression is entered, **When** the user presses equals, **Then** the frontend sends the raw expression to the backend and displays the numeric result returned in the backend's JSON response.
2. **Given** a valid expression with operator precedence (e.g., `2 + 3 * 4`), **When** evaluated by the backend, **Then** the returned result respects standard mathematical precedence (`14`, not `20`).
3. **Given** the result is displayed, **When** the user inspects behavior, **Then** no arithmetic was computed on the client — the displayed value comes solely from the backend response.

---

### User Story 2 - Receive and surface validation errors (Priority: P2)

When a user submits input that is not a valid mathematical expression, the backend rejects it with a structured error describing the problem, and the frontend surfaces that error to the user as an alert.

**Why this priority**: Real input is messy; without clear error handling the calculator appears broken or silently fails. Depends on the P1 request path existing.

**Independent Test**: Submit malformed input (e.g., `5 + * 2`, `(3 + )`, division by zero) and confirm the backend returns an error with details and the frontend shows an alert containing a human-readable message.

**Acceptance Scenarios**:

1. **Given** a syntactically invalid expression, **When** the user presses equals, **Then** the backend returns an error response with a message describing why it is invalid, and the display is not updated with a result.
2. **Given** an error response is received, **When** the frontend processes it, **Then** an alert is shown to the user with the error details.
3. **Given** division by zero or another mathematically undefined operation, **When** evaluated, **Then** the backend returns an error (not a result), and the frontend alerts the user.
4. **Given** empty or whitespace-only input, **When** the user presses equals, **Then** the request is rejected with a validation error (or not sent), and the user is informed.

---

### User Story 3 - Prevent concurrent requests with in-progress feedback (Priority: P3)

While a calculation request is being processed by the backend, the frontend visually indicates progress by showing a spinner in place of the equals control and prevents additional submissions until the in-flight request resolves.

**Why this priority**: Improves correctness and UX under latency by preventing duplicate/overlapping requests and giving clear feedback. It is an enhancement on top of the working request/response path (P1/P2).

**Independent Test**: Trigger a calculation, and while the response is pending, confirm a spinner replaces the equals symbol and repeated presses do not issue additional requests; once the response arrives, the equals control returns to normal.

**Acceptance Scenarios**:

1. **Given** a calculation request is in flight, **When** the user presses equals again, **Then** no additional request is sent.
2. **Given** a calculation request is in flight, **When** the user views the equals control, **Then** a spinner is shown in place of the equals symbol.
3. **Given** the backend responds (success or error), **When** the response is handled, **Then** the spinner is removed, the equals control is restored, and further submissions are allowed again.

---

### Edge Cases

- **Backend unreachable / network failure**: The frontend must surface a failure alert and restore the equals control (no permanent lock).
- **Slow backend response**: The spinner persists and the input remains locked until the response or a timeout; on timeout the user is informed and the control is restored.
- **Very long or deeply nested expression**: Backend enforces input limits and returns a clear validation error rather than hanging.
- **Result that is not finite** (e.g., overflow): Backend returns an error rather than a non-numeric result.
- **Non-numeric or unsupported characters / operators**: Rejected with a descriptive validation error.
- **Rapid repeated equals presses**: Only one request is in flight; extras are ignored while locked.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The frontend MUST NOT perform any mathematical evaluation of user expressions; all calculation MUST be delegated to the backend.
- **FR-002**: The frontend MUST send the user's expression to the backend when the user requests a calculation (presses equals).
- **FR-003**: The backend MUST receive the submitted input, validate it, and determine whether it is a well-formed mathematical expression.
- **FR-004**: When the input is a valid expression, the backend MUST evaluate it using standard operator precedence and return the result in a JSON response.
- **FR-005**: When the input is invalid, the backend MUST return an error response containing details describing why the input was rejected, and MUST NOT return a computed result.
- **FR-006**: The backend MUST reject mathematically undefined operations (e.g., division by zero) and non-finite results with an error response.
- **FR-007**: The frontend MUST display the backend-returned result on the calculator display when the response is successful.
- **FR-008**: The frontend MUST display an alert to the user containing the error details when the backend returns an error.
- **FR-009**: The frontend MUST prevent additional calculation requests from being submitted while a request is already in progress.
- **FR-010**: The frontend MUST replace the equals control with a spinner indicator while a calculation request is in progress.
- **FR-011**: The frontend MUST restore the equals control and re-enable submissions once the request resolves (success, error, or failure).
- **FR-012**: The frontend MUST handle backend unreachability and request timeouts gracefully by alerting the user and restoring the equals control.
- **FR-013**: The integration MUST follow the documented API contract for the request and response shapes (success and error), as the single source of truth between frontend and backend.
- **FR-014**: The backend MUST emit structured logs for calculation requests, including validation failures, sufficient to diagnose errors from logs alone.

### Key Entities *(include if feature involves data)*

- **Calculation Request**: The user-submitted input to be evaluated — the raw mathematical expression as text.
- **Calculation Result**: The successful outcome — the numeric value produced by evaluating a valid expression, returned as JSON.
- **Calculation Error**: The failure outcome — a structured description of why input was rejected or could not be evaluated (e.g., reason/message, and where applicable the offending part of the input).

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 100% of displayed results for valid expressions originate from the backend; zero arithmetic is computed on the client.
- **SC-002**: For valid expressions, the displayed result is mathematically correct (respecting precedence) in 100% of tested cases.
- **SC-003**: 100% of invalid inputs result in a user-visible error alert and no result being displayed.
- **SC-004**: While a request is in progress, at most one request is ever in flight; duplicate submissions during that window are 0.
- **SC-005**: The in-progress spinner appears within 100 ms of the user pressing equals and is removed within 100 ms of the response resolving.
- **SC-006**: Network failures and timeouts always restore the equals control (no permanent lock) and inform the user in 100% of failure cases.
- **SC-007**: Users can complete a typical calculate→result cycle without any unhandled error or stuck UI state.

## Assumptions

- The backend exposes a single calculation endpoint that accepts the raw expression as text and returns JSON for both success and error outcomes.
- Standard mathematical conventions apply: `+`, `-`, `*`, `/`, parentheses, decimals, and standard operator precedence; advanced/scientific functions are out of scope for this feature unless already supported.
- "Display alerts for errors" means a user-visible, non-blocking notification surfacing the backend's error message; exact visual treatment is an implementation detail.
- "Spinner in the '=' sign" means the equals control visually swaps to a loading indicator while the request is pending.
- A reasonable client-side request timeout is applied (default assumption: a few seconds) after which the user is informed and the control is restored.
- The existing client-side calculation logic is removed or bypassed so it can no longer be the source of displayed results.
- The frontend and backend run as separate workspaces communicating over HTTP/JSON per the project's contract-driven integration principle; CORS is already handled by the backend.
