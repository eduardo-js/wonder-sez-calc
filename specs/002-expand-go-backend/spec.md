# Feature Specification: Expand Go Backend

**Feature Branch**: `002-expand-go-backend`

**Created**: 2026-06-02

**Status**: Draft

**Input**: User description: "Let's expand the go backend: (1) add CORS for the frontend url, (2) add a framework for handling HTTP, gin preferred, (3) input validation — package or reflect?, (4) define project standards: stretchr/testify, prefer mocks, aim for >95% coverage, always prefer table-driven test structure, (5) ensure clear routing."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Frontend can call the backend across origins (Priority: P1)

The browser-based calculator frontend, served from its own origin, makes HTTP requests to
the backend API and the browser permits the responses instead of blocking them as
cross-origin violations.

**Why this priority**: Without cross-origin access the frontend cannot consume any backend
endpoint from the browser. This unblocks all future frontend/backend integration and is the
first thing that must work.

**Independent Test**: From the frontend dev origin, issue a request (including the browser's
preflight `OPTIONS`) to a backend endpoint and confirm the response carries the headers that
allow the configured frontend origin, while a disallowed origin is not granted access.

**Acceptance Scenarios**:

1. **Given** a request from the configured frontend origin, **When** the browser sends a preflight `OPTIONS`, **Then** the backend responds with headers permitting that origin, the required methods, and the required headers.
2. **Given** a request from the configured frontend origin, **When** an actual request is made, **Then** the response includes the header that allows that origin so the browser exposes the body.
3. **Given** a request from an origin that is not configured, **When** the request is made, **Then** the response does not grant cross-origin access to that origin.
4. **Given** the deployment environment changes, **When** the allowed frontend origin is configured, **Then** the backend honors the configured value without code changes.

---

### User Story 2 - Malformed requests are rejected with clear errors (Priority: P2)

A client (frontend or any caller) sends a request whose body or parameters are missing
required fields, have the wrong type, or fall outside allowed ranges, and the backend
rejects it with a clear, structured error instead of processing bad data or crashing.

**Why this priority**: Trustworthy input handling protects correctness and reliability. It
builds on the connectivity established in P1 and is required before any real endpoint can
safely accept data.

**Independent Test**: Send requests with valid and invalid payloads to a validated endpoint
and confirm valid input is accepted while each class of invalid input is rejected with a
consistent error shape and an appropriate client-error status.

**Acceptance Scenarios**:

1. **Given** a request missing a required field, **When** it is received, **Then** the backend returns a client-error status with a structured message identifying the problem.
2. **Given** a request with a field of the wrong type or out-of-range value, **When** it is received, **Then** the backend rejects it with a structured validation error and does not process the request.
3. **Given** a well-formed, valid request, **When** it is received, **Then** the backend accepts it and proceeds normally.
4. **Given** any validation failure, **When** the error is returned, **Then** the response body follows a single consistent error format across all endpoints.

---

### User Story 3 - Consistent, discoverable routing and quality standards (Priority: P3)

A developer adding a new endpoint finds a clear, central place where routes are registered,
follows documented conventions for grouping and naming, and writes tests that meet the
project's agreed quality bar.

**Why this priority**: Maintainability and a consistent contributor experience matter, but a
working, validated API on the established routes already delivers value. Standards broaden
that value across the team and over time.

**Independent Test**: Inspect the routing setup and confirm all routes are registered in one
discoverable place with consistent grouping; run the test suite and confirm it meets the
coverage and structure standards from a single command.

**Acceptance Scenarios**:

1. **Given** the running backend, **When** a developer inspects the route registration, **Then** all routes (including health/readiness and new API routes) are declared in one discoverable, consistently grouped location.
2. **Given** a new endpoint is added, **When** the developer follows the documented routing convention, **Then** it is registered consistently with existing routes (grouping, versioning prefix, naming).
3. **Given** the test suite, **When** it is run from the orchestration command, **Then** it passes, follows the table-driven structure, and reports coverage at or above the agreed threshold.
4. **Given** an unknown route or unsupported method, **When** a request is made, **Then** the backend returns the appropriate not-found / method-not-allowed response in the standard error format.

---

### Edge Cases

- A preflight `OPTIONS` request for an allowed origin but a disallowed method or header — access is not granted for the unsupported method/header.
- A request body that is empty, oversized, or not valid JSON — rejected with a clear client-error, no crash.
- Multiple allowed origins (e.g. local dev plus a deployed frontend) — each configured origin is honored, others are not.
- A validation error and a not-found error returned from different endpoints — both share the same error response shape.
- Concurrent requests — routing, CORS, and validation behave consistently under parallel load.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The backend MUST permit cross-origin requests from the configured frontend origin(s), including correct handling of browser preflight (`OPTIONS`) requests, allowed methods, and allowed headers.
- **FR-002**: The set of allowed frontend origin(s) MUST be configurable per environment (e.g. local development vs. deployed) without source code changes.
- **FR-003**: The backend MUST NOT grant cross-origin access to origins that are not in the configured allow-list.
- **FR-004**: The backend MUST route HTTP requests through a single, consistent HTTP-handling layer that supports route grouping, path parameters, and middleware (CORS, validation, logging) applied uniformly.
- **FR-005**: All routes — existing health/readiness probes and any new API routes — MUST be registered in one central, discoverable location with consistent grouping and a versioned API prefix for application routes.
- **FR-006**: The backend MUST validate incoming request payloads and parameters against declared rules (required fields, types, ranges/formats) and reject invalid input with a client-error status before any business logic runs.
- **FR-007**: Validation and routing errors MUST be returned in a single consistent, structured response format used across all endpoints.
- **FR-008**: The backend MUST return appropriate responses for unknown routes (not found) and unsupported methods (method not allowed) in the standard error format.
- **FR-009**: The backend MUST preserve existing behavior and contracts of the current health (`/healthz`) and readiness (`/readyz`) endpoints after migration to the new routing layer.
- **FR-010**: The backend MUST continue to emit structured logs and handle errors explicitly (no swallowed errors), per the project constitution, including for CORS rejections and validation failures.
- **FR-011**: The project MUST adopt documented backend testing standards: a shared assertion/mocking toolkit, preference for mocking external collaborators, table-driven test structure as the default, and a target test coverage of **at least 95%**.
- **FR-012**: The test suite MUST be runnable from the existing single orchestration entry point (Makefile) and MUST report coverage so the threshold can be verified.
- **FR-013**: New third-party dependencies (HTTP framework, validation, test tooling) MUST be pinned and justified per the constitution's dependency policy.

### Key Entities *(include if data involved)*

- **Route**: A registered HTTP endpoint defined by its method, path (under a versioned group where applicable), the middleware applied to it, and its handler.
- **CORS Policy**: The configured set of allowed origin(s), methods, and headers governing which browser origins may consume the API.
- **Validation Rule Set**: The declared constraints (required, type, range, format) bound to a request's fields, used to accept or reject input.
- **Error Response**: The single structured shape returned for validation, not-found, method-not-allowed, and other client/server errors.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: The frontend, running on its configured origin, can successfully call a backend endpoint from the browser (including preflight) with no CORS errors; a non-allowed origin is denied.
- **SC-002**: 100% of defined invalid-input scenarios (missing required field, wrong type, out-of-range, malformed body) are rejected with a client-error status and the standard error shape — none cause a crash or process invalid data.
- **SC-003**: Every endpoint returns errors in one consistent structured format, verified across at least the validation, not-found, and method-not-allowed cases.
- **SC-004**: All routes are registered in a single discoverable location, and a developer can add a new route following the documented convention without touching unrelated code.
- **SC-005**: The backend test suite runs green from the single orchestration command, uses table-driven tests as the default structure, and reports total coverage of **at least 95%**.
- **SC-006**: Existing `/healthz` and `/readyz` behavior is unchanged after the routing migration, verified by passing tests.

## Assumptions

- **HTTP framework**: Gin is adopted as the HTTP-handling framework (per explicit user preference), replacing the raw `net/http` `ServeMux` while preserving current endpoint behavior.
- **CORS**: The allowed origin defaults to the local frontend dev origin (Vite default `http://localhost:5173`) and is configurable via environment/config for other environments.
- **Input validation — package vs. reflect**: A declarative validation **package** (`go-playground/validator`, which Gin integrates natively via struct tags) is used rather than hand-rolled reflection. Rationale: it is the de-facto standard for Gin, less error-prone, and avoids maintaining custom reflection code. This resolves the open question in the user input in favor of a maintained package.
- **Test tooling**: `stretchr/testify` (assert/require + mock) is the standard assertion and mocking toolkit; mocks are preferred for external collaborators; table-driven tests are the default structure; coverage target is **≥95%**.
- **Routing convention**: Application routes live under a versioned prefix (e.g. `/api/v1/...`); health/readiness probes remain at their current top-level paths; all registration is centralized.
- **Scope**: This feature establishes the backend HTTP/validation/testing foundation and migrates existing health/readiness endpoints onto it. Defining concrete calculator/business endpoints and the full frontend↔backend API contract is **out of scope** here and handled in a later feature, consistent with the barebones backend established in feature 001.
- **Constitution alignment**: Structured JSON logging, explicit error handling, pinned/justified dependencies, monorepo boundaries, and test-first discipline continue to apply.
