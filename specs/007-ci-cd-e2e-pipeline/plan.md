# Implementation Plan: CI & E2E Pipeline

**Branch**: `007-ci-cd-e2e-pipeline` | **Date**: 2026-06-03 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `/specs/007-ci-cd-e2e-pipeline/spec.md`

## Summary

Add a GitHub Actions pipeline that gates every PR and push to `main` with the project's existing
quality checks (lint, type-check, unit tests, backend coverage gate, production build) for both
workspaces, plus an end-to-end job that boots the full stack via the existing Docker Compose setup
(feature 006) and drives the calculator through a real browser (Playwright) to verify a correct
result. CI invokes the same commands the local `make`/npm/go workflow runs, so "green locally" ==
"green in CI". **CD is out of scope** (deferred). Concurrency cancels superseded runs; e2e failures
upload Playwright traces/screenshots as artifacts.

Repo: `origin git@github.com:eduardo-js/wonder-sez-calc.git`. The branch will be pushed and become
the repository's default branch `main`; the pipeline gates `main` accordingly.

## Technical Context

**Language/Version**: Frontend TS 5.5 / Node 24 (Vite 5, React 18); Backend Go 1.25. Pipeline:
GitHub Actions (YAML).

**Primary Dependencies**: GitHub Actions (`actions/checkout`, `actions/setup-node`, `actions/setup-go`,
`actions/upload-artifact`). NEW dev dependency: **Playwright** (`@playwright/test`) for e2e. Reuses
Docker Compose stack from feature 006.

**Storage**: N/A.

**Testing**: Existing Vitest (frontend) + Go `testing` (backend) invoked via the same npm/go commands
the Makefile uses; NEW Playwright e2e suite under `e2e/` driving the running stack.

**Target Platform**: `ubuntu-latest` GitHub-hosted runners.

**Project Type**: Web application monorepo; this feature adds a CI layer (`.github/workflows/`) + an
e2e test layer (`e2e/`).

**Performance Goals**: Full pipeline < 15 min for a typical change (SC-005); cancel superseded runs.

**Constraints**: Mirror local commands (SC-006/FR-013); fail-fast on unhealthy stack (FR-008); no
secrets to forked-PR runs (FR-009); gates on `main` (FR-001).

**Scale/Scope**: 3 jobs (frontend, backend, e2e). Single modern browser (Chromium) for v1. No CD.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- **I. Monorepo Cohesion & Boundaries** — PASS. Workflows at repo root `.github/workflows/`; e2e tests
  in a dedicated top-level `e2e/` folder. No frontend/backend source coupling; e2e talks to the stack
  only over HTTP (browser → served UI → backend), honoring the contract boundary.
- **II. Test-First (NON-NEGOTIABLE)** — PASS WITH NOTE. The e2e suite *is* tests and encodes the US2
  acceptance scenarios; written to assert correct behavior and fail on contract breakage. The workflow
  YAML has no unit-testable logic — verified by executing on a PR (a deliberately failing change must
  block; a clean change must pass). Existing unit tests are unchanged and remain the test-first gate
  for app code.
- **III. Maintainability & Code Quality** — PASS. This feature *enforces* the constitution's CI mandate
  ("every push MUST run lint, type-check, unit tests, build for both workspaces"). Pin action versions;
  Playwright is a well-maintained, standard e2e tool with native Actions support and trace/screenshot
  diagnostics — justified new dependency.
- **IV. Contract-Driven Integration** — PASS. E2E exercises the REAL frontend↔backend contract against
  the running stack (not mocks), directly satisfying "integration tests MUST exercise the real contract
  boundary".
- **V. Observability & Operability** — PASS. E2E uses backend `/healthz` to gate startup; failures emit
  downloadable logs + Playwright traces/screenshots (FR-007), making CI failures diagnosable from
  artifacts alone (SC-004).

No violations → Complexity Tracking empty.

## Project Structure

### Documentation (this feature)

```text
specs/007-ci-cd-e2e-pipeline/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 — workflow/jobs config model
├── quickstart.md        # Phase 1 — run e2e locally + CI behavior
├── contracts/
│   ├── pipeline.md      # Triggers, jobs, required checks, artifacts, concurrency
│   └── e2e-scenarios.md # E2E acceptance scenarios (calculation journey)
└── checklists/
    └── requirements.md  # From /speckit-specify
```

### Source Code (repository root)

```text
.github/
└── workflows/
    └── ci.yml               # NEW — frontend, backend, e2e jobs; triggers on PR/push to main; concurrency cancel

e2e/                          # NEW — Playwright e2e (separate from Vitest scope)
├── package.json             # @playwright/test
├── playwright.config.ts     # baseURL http://localhost:5173, Chromium, trace/screenshot on-failure, retries
├── tests/
│   └── calculator.spec.ts   # US2: enter expression → assert correct result
└── README.md                # run locally against the compose stack

frontend/                     # scripts unchanged (lint, typecheck, test, build) — invoked by CI
backend/                      # unchanged — go vet, make test-backend (coverage gate), go build — invoked by CI
Makefile                      # EDIT — add `test-e2e` (+ `e2e`) convenience target reused by CI
```

**Structure Decision**: Web-application monorepo. CI definition is cross-cutting → repo root
`.github/workflows/`. E2E is a new, separate test layer in top-level `e2e/` so it does not pollute the
frontend Vitest scope (Principle I boundaries; unit vs e2e runners stay distinct). CI runs frontend and
backend as parallel jobs using the same npm/go commands the Makefile uses (parity → SC-006), with the
e2e job gated on both via `needs:`.

## Complexity Tracking

> No constitution violations. Section intentionally empty.
