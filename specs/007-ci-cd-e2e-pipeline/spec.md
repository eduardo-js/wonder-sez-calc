# Feature Specification: CI & E2E Pipeline

**Feature Branch**: `007-ci-cd-e2e-pipeline`

**Created**: 2026-06-03

**Status**: Draft

**Input**: User description: "lets add CI CD and e2e tests to run on pipeline using github actions"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Quality gate on every change (Priority: P1)

A contributor opens a pull request. An automated pipeline runs the project's quality gates —
lint, type-check, unit tests for both the frontend and backend, the backend coverage threshold,
and a production build — and reports a clear pass/fail status back on the pull request. A red
pipeline blocks the merge; a green one clears it.

**Why this priority**: This is the core promise of CI and is mandated by the project constitution
("Every push MUST run lint, type-check, unit tests, and build for both workspaces; a red pipeline
blocks merge"). It delivers value immediately and independently of e2e or deployment.

**Independent Test**: Open a PR with a deliberately failing lint/test change and confirm the
pipeline fails and blocks merge; open a clean PR and confirm it passes and clears merge.

**Acceptance Scenarios**:

1. **Given** a PR targeting the default branch, **When** the pipeline runs, **Then** it executes lint, type-check, unit tests, coverage gate, and build for both frontend and backend, and reports a single overall status.
2. **Given** a change that breaks a unit test or lint rule, **When** the pipeline runs, **Then** the pipeline fails and the failure is visible on the PR.
3. **Given** the backend coverage falls below the project threshold, **When** the pipeline runs, **Then** the pipeline fails.
4. **Given** a green pipeline, **When** a reviewer approves, **Then** the change is mergeable.

---

### User Story 2 - End-to-end verification of the running stack (Priority: P2)

The pipeline starts the full application (frontend + backend together) and runs end-to-end tests
that drive the calculator through a real browser against the running services, confirming a user
can perform a calculation and see the correct result. This catches integration/contract breakage
that unit tests miss.

**Why this priority**: E2E is explicitly requested and validates the real frontend↔backend contract
boundary (constitution Principle IV). It depends on the CI foundation (US1) existing first.

**Independent Test**: Run the e2e job locally/in CI against the started stack; confirm a calculation
scenario passes, and that an intentionally broken contract (e.g. wrong response shape) fails the job.

**Acceptance Scenarios**:

1. **Given** the pipeline has started the full stack, **When** the e2e suite runs, **Then** it performs at least one real calculation through the UI and asserts the displayed result is correct.
2. **Given** the backend is unreachable or returns an unexpected response, **When** the e2e suite runs, **Then** the e2e job fails with a diagnostic artifact (logs/screenshot/trace).
3. **Given** an e2e run completes, **When** a maintainer inspects the run, **Then** test results and failure artifacts are downloadable from the pipeline.

---

### Edge Cases

- A flaky e2e run: the pipeline must distinguish an infrastructure flake from a real failure (retry policy and artifacts), and must not silently pass on flake.
- The full stack fails to become healthy within a time budget: the e2e job must time out with a clear message rather than hang.
- A forked-repository PR: secrets must not be exposed to untrusted PR runs, while still running the no-secret quality gates.
- A pipeline run on a change that touches only one workspace: behavior must be defined (run both, or only the affected one) and consistent.
- Concurrent pushes to the same branch: superseded runs should be cancellable to save resources.
- A dependency or base-image fetch fails transiently: the run should fail clearly, and be safely re-runnable.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The pipeline MUST run automatically on every pull request targeting the `main` branch and on every push to `main`.
- **FR-002**: The pipeline MUST run lint, type-check, unit tests, and a production build for BOTH the frontend and backend workspaces.
- **FR-003**: The pipeline MUST enforce the existing backend coverage threshold and fail when coverage is below it.
- **FR-004**: The pipeline MUST report an overall pass/fail status visible on the pull request, suitable for use as a required merge gate.
- **FR-005**: The pipeline MUST run an end-to-end suite that exercises the running full stack (frontend served + backend reachable) through a real browser and verifies a correct calculation result.
- **FR-006**: The e2e job MUST start the full application stack and wait until it is healthy before running tests, and MUST tear it down afterward.
- **FR-007**: On e2e failure, the pipeline MUST capture and expose diagnostic artifacts (logs, and browser screenshots/traces).
- **FR-008**: The pipeline MUST fail fast and clearly on infrastructure errors (stack not healthy within a time budget, dependency fetch failure) rather than hang indefinitely.
- **FR-009**: Secrets MUST NOT be exposed to untrusted (forked-PR) runs; quality gates that need no secrets MUST still run for such PRs.
- **FR-010**: The pipeline MUST be re-runnable and idempotent — re-running a run on the same commit yields the same gating decision.
- **FR-011**: Redundant in-progress runs for the same branch/PR SHOULD be cancelled when a newer commit arrives.
- **FR-012**: Pipeline definitions and any added test tooling MUST live in the repository and be reviewable like any other code.
- **FR-013**: The local developer workflow (existing `make` targets) MUST remain the source of truth the pipeline invokes, so "green locally" and "green in CI" stay consistent.

### Key Entities

- **Pipeline run**: One automated execution triggered by a PR or push; has an overall status, per-job results, logs, and downloadable artifacts.
- **Quality gate job**: A unit of the pipeline (lint, type-check, unit tests, coverage, build) per workspace, each pass/fail.
- **E2E job**: A job that brings up the full stack, runs browser-driven scenarios, and produces result + failure artifacts.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 100% of pull requests targeting the default branch automatically trigger the pipeline with no manual action.
- **SC-002**: A change that breaks lint, a unit test, the coverage threshold, or the build is blocked from merging in 100% of cases.
- **SC-003**: The e2e suite verifies at least one full calculation user journey end-to-end against the running stack on every pipeline run.
- **SC-004**: When an e2e test fails, a maintainer can diagnose the cause from pipeline artifacts alone (no local reproduction needed) in at least 90% of cases.
- **SC-005**: The full pipeline (gates + e2e) completes within 15 minutes for a typical change, keeping feedback fast.
- **SC-006**: "Passes locally" and "passes in CI" agree — the pipeline runs the same checks the documented local workflow runs, with zero known divergences.

## Assumptions

- The repository is (or will be) hosted on GitHub and GitHub Actions is the pipeline platform (per the request).
- The existing `make` targets (`lint`, `test`, `build`, and the Docker compose targets) are the canonical commands the pipeline invokes, so CI mirrors local behavior.
- The full stack can be started in the pipeline via the existing Docker Compose setup (feature 006), and the backend `/healthz` endpoint signals readiness.
- E2E covers the primary calculator journey (enter expression → see correct result); exhaustive cross-browser/device matrices are out of scope for v1 (single modern browser).
- Unit/integration tests already exist for both workspaces; this feature orchestrates and gates them, and adds the e2e layer — it does not rewrite existing tests.
- **CD is out of scope for v1** (deferred by decision): the pipeline is CI + e2e only. No artifact publishing, registry push, or deployment is included until a delivery target is defined. The branch name retains "cd" for continuity, but no continuous-delivery requirements are in this spec.
- The protected default branch is `main`.
- Performance, security scanning, and release versioning/changelogs are out of scope for v1 unless added later.
