<!--
SYNC IMPACT REPORT
==================
Version change: [TEMPLATE/unversioned] → 1.0.0
Bump rationale: Initial ratification — all placeholders replaced with concrete principles.

Modified principles: N/A (initial adoption)
Added sections:
  - Core Principles (5): Monorepo Cohesion & Boundaries, Test-First (NON-NEGOTIABLE),
    Maintainability & Code Quality, Contract-Driven Integration, Observability & Operability
  - Technology & Structure Constraints
  - Development Workflow & Quality Gates
  - Governance
Removed sections: None

Templates requiring updates:
  - .specify/templates/plan-template.md ......... ✅ compatible (Constitution Check is generic)
  - .specify/templates/spec-template.md ......... ✅ compatible
  - .specify/templates/tasks-template.md ........ ✅ compatible
  - .claude/agents/*.md ......................... ✅ added (project-architect, frontend-react, backend-go)

Deferred TODOs: None
-->

# wonderlic-calc Constitution

## Core Principles

### I. Monorepo Cohesion & Clear Boundaries
The repository is a single monorepo hosting a React frontend and a Go backend. Each
side MUST live in its own top-level workspace (`frontend/`, `backend/`) with independent
build, test, and dependency management. Cross-cutting contracts (API schemas, shared
types) MUST live in a designated shared location and be the single source of truth.
Direct imports across the frontend/backend boundary are FORBIDDEN; integration happens
only through versioned, documented interfaces (HTTP/JSON, generated clients).

Rationale: Clear boundaries keep each stack independently buildable and testable while
preventing hidden coupling that erodes a monorepo over time.

### II. Test-First (NON-NEGOTIABLE)
TDD is mandatory. For every change: write the failing test → confirm it fails → implement
until green → refactor. Backend Go code MUST use the standard `testing` package with table-
driven tests; frontend MUST use a component/unit test runner (Vitest/Jest + Testing Library).
No feature merges without tests covering its acceptance criteria. Bug fixes MUST include a
regression test that fails before the fix.

Rationale: Tests written after the fact codify whatever was built, not what was intended.
Test-first is the only reliable guard for a calculator whose correctness is the product.

### III. Maintainability & Code Quality
Code MUST pass automated linting and formatting before merge: Go via `gofmt`/`goimports`
and `go vet` (plus `golangci-lint`); frontend via ESLint + Prettier and a strict
TypeScript config (`strict: true`, no implicit `any`). Functions and modules MUST have a
single clear responsibility. Public functions and exported types MUST be documented.
Dead code, commented-out blocks, and unjustified complexity MUST NOT be merged.

Rationale: Consistent, linted, typed code is the baseline that makes a long-lived monorepo
safe to change.

### IV. Contract-Driven Integration
The API contract between frontend and backend MUST be defined explicitly (OpenAPI or an
equivalent typed schema) and treated as the source of truth. Frontend clients and backend
handlers MUST be validated against this contract. Breaking changes to the contract MUST be
versioned and accompanied by a migration note. Integration tests MUST exercise the real
contract boundary, not mocked-only assumptions.

Rationale: A shared, enforced contract is what lets two independently developed stacks
evolve without silent breakage.

### V. Observability & Operability
Backend services MUST emit structured logs (JSON) and expose health/readiness endpoints.
Errors MUST be handled explicitly and surfaced with actionable context — no swallowed
errors, no silent failures. Frontend MUST surface failures to users gracefully and report
errors for diagnosis. Every externally visible behavior MUST be observable enough to debug
from logs/telemetry alone.

Rationale: You cannot operate or trust what you cannot observe; explicit error handling is
non-negotiable for reliability.

## Technology & Structure Constraints

- **Languages/Stacks**: React + TypeScript (frontend), Go (backend). No additional runtime
  stacks without an amendment.
- **Repository layout** (canonical):
  - `frontend/` — React app, self-contained build and tests.
  - `backend/` — Go modules, self-contained build and tests.
  - `shared/` (or `api/`) — API contracts and generated/shared types.
  - `.specify/` — Spec Kit artifacts; `specs/` — per-feature specs, plans, tasks.
- **Dependencies**: Pin versions. New third-party dependencies MUST be justified (need,
  maintenance status, license) in the relevant plan or PR.
- **CI**: Every push MUST run lint, type-check, unit tests, and build for both workspaces;
  a red pipeline blocks merge.

## Development Workflow & Quality Gates

- **Branching**: Feature work happens on feature branches created via the Spec Kit flow
  (`/speckit-specify` → `/speckit-plan` → `/speckit-tasks` → `/speckit-implement`).
- **Quality gates before merge**: tests pass, lint/format clean, type-check clean,
  contract validated, and the change traces to a spec/task.
- **Code review**: Every change requires review. Reviewers MUST verify constitution
  compliance (tests-first evidence, boundaries respected, contract honored).
- **Project subagents** (see `.claude/agents/`): use them for their domains —
  - `project-architect` — whole-repo architecture, boundaries, contracts, cross-cutting design.
  - `frontend-react` — React/TypeScript frontend implementation and review.
  - `backend-go` — Go backend implementation and review.

## Governance

This constitution supersedes all other practices. When guidance conflicts, the constitution
wins. Amendments MUST be proposed via PR, documented in the Sync Impact Report, approved by
the maintainer, and accompanied by any required migration steps.

Versioning policy (semantic):
- **MAJOR**: backward-incompatible governance/principle removal or redefinition.
- **MINOR**: new principle/section added or materially expanded guidance.
- **PATCH**: clarifications, wording, or non-semantic refinements.

Compliance is verified at code review and in CI quality gates. Complexity that violates a
principle MUST be justified in writing or removed. Runtime/agent guidance lives in
`CLAUDE.md` and `.claude/agents/`.

**Version**: 1.0.0 | **Ratified**: 2026-06-02 | **Last Amended**: 2026-06-02
