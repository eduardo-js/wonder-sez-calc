# Implementation Plan: Docker Compose Setup

**Branch**: `006-docker-compose-setup` | **Date**: 2026-06-03 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `/specs/006-docker-compose-setup/spec.md`

## Summary

Add a Docker Compose stack that builds and runs the existing React frontend and Go backend
as two services, startable with one command and no host toolchain. Configuration (listen
address, CORS origins, backend URL) is declared **inline** in the compose/Dockerfiles — no
`.env` files. New `make` targets wrap the container lifecycle (up, down, build, logs) and
appear in `make help`. The browser (on the host) calls the host-mapped backend port directly,
so backend CORS must permit the frontend's host origin; inter-container DNS is not required
for browser→backend traffic.

## Technical Context

**Language/Version**: Frontend — TypeScript 5.5 / Node 24 (Vite 5, React 18). Backend — Go 1.25.

**Primary Dependencies**: Docker Engine + Docker Compose v2 (`docker compose`). No new application
dependencies. Base images: `node:24-alpine` (frontend), `golang:1.25-alpine` (backend build).

**Storage**: N/A (stateless calculator).

**Testing**: Existing `make test` (Vitest + Go `testing`). Container layer verified by build
success + quickstart smoke test (health endpoint + a calculation round-trip).

**Target Platform**: Local developer machine with a container engine (Linux/macOS/WSL2).

**Project Type**: Web application (frontend + backend) in a monorepo.

**Performance Goals**: Clone → running stack < 10 min on broadband (SC-004); not a runtime perf feature.

**Constraints**: No host Node/Go required (FR-002); no `.env` files (FR-014); existing `make`
targets unchanged (FR-013); fail-fast on port collision (FR-011).

**Scale/Scope**: Single-developer local stack; 2 services. Out of scope: registries, k8s, TLS, auth, prod deploy.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- **I. Monorepo Cohesion & Boundaries** — PASS. Each service gets its own `Dockerfile` inside
  its workspace (`frontend/Dockerfile`, `backend/Dockerfile`); `docker-compose.yml` at root wires
  them. No cross-boundary imports; integration stays over HTTP/JSON as today.
- **II. Test-First (NON-NEGOTIABLE)** — PASS WITH NOTE. This feature adds infrastructure
  (Dockerfiles, compose, Make targets, one Vite `server` config line) with no unit-testable
  business logic. There is nothing to drive via a failing unit test; correctness is verified by
  (a) image builds succeeding and (b) the quickstart smoke test (healthz + calculate round-trip),
  which is the acceptance gate. All existing app tests remain green (FR-013). No production code
  logic is added, so no TDD violation.
- **III. Maintainability & Code Quality** — PASS. Pin base image tags; multi-stage builds; minimal
  `.dockerignore` per workspace; no dead config. The single source change (`vite.config.ts` `server`
  block) stays typed and lint-clean.
- **IV. Contract-Driven Integration** — PASS. No API contract change. The compose service contract
  (ports, env, healthcheck, dependency order) is documented in `contracts/compose-services.md`.
  CORS origin alignment preserves the existing contract boundary.
- **V. Observability & Operability** — PASS. Backend already emits JSON logs and exposes `/healthz`;
  compose uses `/healthz` as the backend healthcheck and `make docker-logs` surfaces logs. Fail-fast
  on bind errors is inherent to compose port mapping (FR-011).

No violations → Complexity Tracking left empty.

## Project Structure

### Documentation (this feature)

```text
specs/006-docker-compose-setup/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output (config/service model)
├── quickstart.md        # Phase 1 output (run/stop/rebuild + smoke test)
├── contracts/
│   └── compose-services.md   # Service interface contract (ports, env, healthcheck, deps)
└── checklists/
    └── requirements.md  # From /speckit-specify
```

### Source Code (repository root)

```text
docker-compose.yml          # NEW — root orchestration; 2 services, inline env, ports, healthcheck, depends_on
Makefile                    # EDIT — add docker-up / docker-down / docker-build / docker-logs / docker-rebuild

frontend/
├── Dockerfile              # NEW — node:24-alpine; install deps; run Vite dev server (host 0.0.0.0:5173)
├── .dockerignore           # NEW — exclude node_modules, dist, coverage
└── vite.config.ts          # EDIT — add server: { host: true, port: 5173 } for container reachability + HMR

backend/
├── Dockerfile              # NEW — multi-stage: golang:1.25-alpine build → minimal runtime; expose :8080
└── .dockerignore           # NEW — exclude bin, coverage, *.out
```

**Structure Decision**: Web-application monorepo (Option 2). Dockerfiles live in their respective
workspaces to honor Principle I boundaries; the single `docker-compose.yml` and the `Makefile`
edits live at repo root as the cross-cutting orchestration layer. Only one application source file
changes (`vite.config.ts`) to make the dev server reachable from the host and enable HMR for US3.

## Complexity Tracking

> No constitution violations. Section intentionally empty.
