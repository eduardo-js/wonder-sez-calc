# Feature Specification: Docker Compose Setup

**Feature Branch**: `006-docker-compose-setup`

**Created**: 2026-06-03

**Status**: Draft

**Input**: User description: "add docker compose for frontend and backend project, update make with docker commands"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Run the full stack with one command (Priority: P1)

A developer who has just cloned the repository wants to run the entire application —
frontend and backend together — without installing Node, Go, or any language toolchain
on their host machine. They issue a single command and, after the containers build and
start, the calculator is reachable in a browser and successfully talks to the backend.

**Why this priority**: This is the core value of the feature. A one-command, toolchain-free
startup is what makes containerization worth doing; everything else builds on it.

**Independent Test**: From a clean checkout on a machine with only a container runtime
installed, run the documented startup command, open the frontend URL in a browser, perform
a calculation, and confirm the result returns from the backend.

**Acceptance Scenarios**:

1. **Given** a clean checkout and a running container engine, **When** the developer runs the documented "start everything" command, **Then** both the frontend and backend services build and reach a running state.
2. **Given** both services are running, **When** the developer opens the frontend URL in a browser and performs a calculation, **Then** the frontend reaches the backend and displays the computed result.
3. **Given** the running stack, **When** the developer queries the backend health endpoint, **Then** it reports a healthy status.

---

### User Story 2 - Manage the stack through Make (Priority: P2)

A developer who already uses the project's `make` targets wants container operations to be
available the same way. They can start, stop, rebuild, and view logs for the containerized
stack using `make` targets, consistent with the existing command vocabulary.

**Why this priority**: The repository standardizes developer workflow on `make`. Exposing
container operations through `make` keeps one consistent entry point and lowers the learning
curve, but it depends on the compose definition from US1 existing first.

**Independent Test**: With the compose definition in place, run each new `make` container
target and confirm it performs the expected lifecycle action (start, stop, rebuild, logs).

**Acceptance Scenarios**:

1. **Given** the repository, **When** the developer runs the `make` target that starts the containerized stack, **Then** both services start.
2. **Given** a running containerized stack, **When** the developer runs the `make` target that stops it, **Then** both services stop and their containers are removed.
3. **Given** the repository, **When** the developer runs `make help`, **Then** the new container targets are listed with descriptions alongside existing targets.
4. **Given** source changes, **When** the developer runs the `make` target that rebuilds images, **Then** images rebuild and the stack runs the updated code.

---

### User Story 3 - Iterate on code while containerized (Priority: P3)

A developer running the stack in containers makes a change to frontend or backend source
and wants to see it reflected without a slow, fully manual teardown-and-rebuild cycle.

**Why this priority**: Improves the inner development loop but is a convenience on top of a
working containerized stack; the feature delivers value without it.

**Independent Test**: With the stack running, change a source file and confirm the change is
reflected through the documented rebuild/refresh path within the stated time target.

**Acceptance Scenarios**:

1. **Given** the running stack, **When** the developer changes backend source and follows the documented refresh path, **Then** the running backend reflects the change.
2. **Given** the running stack, **When** the developer changes frontend source and follows the documented refresh path, **Then** the served frontend reflects the change.

---

### Edge Cases

- What happens when a required host port is already in use by another process? The startup must fail with a clear, actionable message rather than hang or silently bind elsewhere.
- How does the system behave when the backend is not yet ready but the frontend is already serving? The frontend must degrade gracefully and recover once the backend becomes healthy.
- What happens when the developer runs the start command twice? The second run must not create duplicate or conflicting containers.
- How are CORS-allowed origins kept consistent between the containerized frontend URL and the backend configuration?
- What happens to the stack on container engine restart or host reboot — do services stay down until explicitly restarted, or come back automatically?

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The repository MUST provide a single container-orchestration definition that runs both the frontend and backend as separate services.
- **FR-002**: A developer MUST be able to start the entire stack with one command, with no language toolchain (Node, Go) installed on the host.
- **FR-003**: The frontend service MUST be reachable from the host browser at a documented, stable URL/port.
- **FR-004**: The backend service MUST be reachable for the frontend to call it, and the frontend MUST be configured to reach the backend within the orchestrated environment.
- **FR-005**: Backend CORS-allowed origins MUST be configured to permit the containerized frontend's origin so cross-service calls succeed.
- **FR-006**: The orchestration MUST surface the backend health/readiness state, and the frontend SHOULD only be considered ready once the backend is reachable.
- **FR-007**: The `Makefile` MUST provide targets to start, stop, rebuild, and view logs for the containerized stack.
- **FR-008**: New `make` container targets MUST appear in `make help` output with descriptions, consistent with the existing self-documenting help format.
- **FR-009**: The container build for each service MUST produce a runnable image from that workspace's source without relying on host-installed dependencies.
- **FR-010**: Host port assignments for each service MUST be documented and configurable to avoid collisions.
- **FR-011**: The startup MUST fail fast with a clear, actionable message when a required port is unavailable or a service fails to build/start.
- **FR-012**: Documentation (README and/or quickstart) MUST describe how to start, stop, and rebuild the containerized stack and what URLs to use.
- **FR-013**: Existing non-container `make` targets (install, build, test, run, dev) MUST continue to work unchanged.
- **FR-014**: Each service's environment configuration MUST be declared inline within the container definitions (compose service `environment` entries and/or `Dockerfile` directives). No separate `.env` files are introduced or committed.
- **FR-015**: Starting the stack MUST apply this inline configuration automatically — a developer passes no environment values by hand.
- **FR-016**: Because this stack carries no real secrets (only the listen address, allowed origins, and backend URL), declaring these values inline in the committed container definitions is acceptable. If a true secret is ever introduced, it MUST NOT be hardcoded inline and MUST be supplied out-of-band.

### Key Entities

- **Compose definition**: The declarative description of the multi-service stack — which services exist, how they build, what ports they expose, their environment, and their dependency/health relationships.
- **Frontend service**: The containerized React application served to the browser.
- **Backend service**: The containerized Go HTTP server exposing the calculator API and health endpoint.
- **Inline service environment**: The non-secret configuration values declared directly in each service's container definition (compose `environment` / `Dockerfile`), applied automatically on startup.
- **Make container targets**: The `make` commands that wrap the orchestration lifecycle (start, stop, rebuild, logs).

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: From a clean checkout on a machine with only a container engine installed (no Node, no Go), a developer can bring the full stack to a working state with a single command.
- **SC-002**: After startup, a developer can open the frontend in a browser, perform a calculation, and receive a correct result from the backend with no manual configuration.
- **SC-003**: 100% of the documented container lifecycle actions (start, stop, rebuild, logs) are available as `make` targets and listed in `make help`.
- **SC-004**: A new contributor can go from clone to a running containerized stack in under 10 minutes on a typical broadband connection, following only the documentation.
- **SC-005**: All pre-existing `make` targets continue to pass and behave identically after the change.
- **SC-006**: A port collision or service failure on startup produces a clear error message identifying the cause within the startup output.

## Assumptions

- A standards-compliant container engine with multi-service orchestration support is available on the developer's machine; installing it is out of scope.
- The backend listens on its existing default address (`:8080`) and the frontend dev/served port aligns with the existing default (`5173`) unless reconfigured; these are the documented host ports.
- The backend's existing `ADDR` and `CORS_ALLOWED_ORIGINS` environment configuration is the mechanism used to wire the services together; these values, and the frontend's backend URL, are declared inline in the container definitions — no `.env` files.
- The stack carries no real secrets, so declaring configuration inline in the committed compose/Dockerfiles is acceptable and is the explicitly chosen approach for zero-config startup.
- Production deployment, image registries/publishing, and orchestration platforms (e.g., Kubernetes) are out of scope; this feature targets local development and evaluation.
- TLS/HTTPS, authentication, and multi-environment promotion are out of scope.
- The existing toolchain-based workflow (`make dev`, etc.) remains the primary path for contributors who have local toolchains; containers are an additive option.
