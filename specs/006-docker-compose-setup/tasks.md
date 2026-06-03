---

description: "Task list for 006-docker-compose-setup"
---

# Tasks: Docker Compose Setup

**Input**: Design documents from `/specs/006-docker-compose-setup/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/compose-services.md, quickstart.md

**Tests**: No automated test tasks. This feature is infrastructure (Dockerfiles, compose, Make targets, one
Vite config line) with no unit-testable business logic. Verification is the quickstart smoke test
(healthz + calculation round-trip) and existing `make test` staying green — per the plan's Constitution
Check note. Existing app tests are NOT modified.

**Organization**: Tasks grouped by user story (US1 → US2 → US3) for independent, incremental delivery.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: US1 / US2 / US3
- Exact file paths included

## Path Conventions

Monorepo web app: `frontend/` (Node 24 / Vite 5), `backend/` (Go 1.25), root for `docker-compose.yml` and `Makefile`.

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Build-context hygiene and the one source change needed for container reachability.

- [X] T001 [P] Create `backend/.dockerignore` excluding `bin/`, `coverage*.out`, `*.test`, `.git`, build artifacts
- [X] T002 [P] Create `frontend/.dockerignore` excluding `node_modules/`, `dist/`, `coverage/`, `.git`
- [X] T003 Add `server: { host: true, port: 5173 }` to `frontend/vite.config.ts` so the dev server binds `0.0.0.0` and is reachable from the host (no other config changed)

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Both service images must build before any compose-based story can run.

**⚠️ CRITICAL**: No user story work can begin until both images build successfully.

- [X] T004 [P] Create `backend/Dockerfile`: multi-stage — `golang:1.25-alpine` builds `./cmd/server`, final slim stage runs the static binary, `EXPOSE 8080`
- [X] T005 [P] Create `frontend/Dockerfile`: `node:24-alpine`, install deps, `CMD` runs Vite dev server on `0.0.0.0:5173`, `EXPOSE 5173`
- [X] T006 Verify both images build standalone: `docker build ./backend` and `docker build ./frontend` succeed

**Checkpoint**: Both images build — compose wiring can begin.

---

## Phase 3: User Story 1 - Run the full stack with one command (Priority: P1) 🎯 MVP

**Goal**: One command brings up frontend + backend; browser performs a calculation that round-trips to the backend.

**Independent Test**: From a clean checkout with only a container engine, run `docker compose up`, open
http://localhost:5173, perform a calculation, and confirm the result returns (no CORS error); `curl localhost:8080/healthz` returns `{"status":"ok"}`.

### Implementation for User Story 1

- [X] T007 [US1] Create root `docker-compose.yml` with `backend` and `frontend` services, `build.context` `./backend` and `./frontend`, port maps `8080:8080` and `5173:5173`
- [X] T008 [US1] Add inline `environment` to the `backend` service in `docker-compose.yml`: `ADDR=:8080`, `CORS_ALLOWED_ORIGINS=http://localhost:5173` (no `.env` file)
- [X] T009 [US1] Add inline `environment` to the `frontend` service in `docker-compose.yml`: `VITE_API_BASE_URL=http://localhost:8080` (no `.env` file)
- [X] T010 [US1] Add a `healthcheck` to the `backend` service hitting `/healthz`, and `depends_on: { backend: { condition: service_healthy } }` to the `frontend` service in `docker-compose.yml`
- [X] T011 [US1] Add `restart: unless-stopped` to both services in `docker-compose.yml`
- [X] T012 [US1] Verify (quickstart smoke): `docker compose up` brings both up, backend reports `healthy`, `curl localhost:8080/healthz` is OK, and a browser calculation at http://localhost:5173 succeeds with no CORS error

**Checkpoint**: MVP — full stack runs and works with a single command.

---

## Phase 4: User Story 2 - Manage the stack through Make (Priority: P2)

**Goal**: Container lifecycle is driven by `make` targets, listed in `make help`, consistent with existing targets.

**Independent Test**: Run each new `make` target and confirm the expected lifecycle action; `make help` lists them with descriptions; existing targets behave identically.

### Implementation for User Story 2

- [X] T013 [US2] In `Makefile`, add `docker-up` (build if needed + start detached), `docker-down` (stop + remove containers/anonymous volumes), `docker-build` (build images only), and `docker-logs` (follow combined logs) targets, each with a `## ` description; add all four to `.PHONY`
- [X] T014 [US2] Verify `make help` lists the new targets with descriptions and that existing targets (`install build test run-frontend run-backend dev` …) still run unchanged (FR-013)

**Checkpoint**: Stack fully manageable via `make`.

---

## Phase 5: User Story 3 - Iterate on code while containerized (Priority: P3)

**Goal**: Frontend changes hot-reload live; backend changes apply via a documented rebuild path.

**Independent Test**: With the stack running, edit a `frontend/src` file and see it reflected live; edit a `backend` file, run the refresh target, and see the change reflected.

### Implementation for User Story 3

- [X] T015 [US3] In `docker-compose.yml`, add bind-mount volumes to the `frontend` service: `./frontend:/app` plus an anonymous `/app/node_modules` volume to preserve container-installed deps (enables Vite HMR)
- [X] T016 [US3] In `Makefile`, add a `docker-rebuild` target (rebuild images + restart stack) with a `## ` description and add it to `.PHONY` (backend refresh path)
- [X] T017 [US3] Verify: editing `frontend/src/**` hot-reloads in the browser; editing `backend/**` then `make docker-rebuild` reflects the change

**Checkpoint**: Inner dev loop works for both services.

---

## Phase 6: Polish & Cross-Cutting Concerns

- [X] T018 [P] Add a "Run with Docker" section to `README.md` documenting `make docker-*` targets and URLs (FR-012)
- [X] T019 [P] Confirm `make test` (frontend + backend) is still green and no existing `make` target regressed
- [X] T020 Run full `quickstart.md` validation end-to-end on a clean checkout (clone → `make docker-up` → smoke test)

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies — start immediately.
- **Foundational (Phase 2)**: Depends on Setup (T003 affects frontend image behavior) — BLOCKS all user stories.
- **US1 (Phase 3)**: Depends on Foundational (images must build). MVP.
- **US2 (Phase 4)**: Depends on US1 (compose file must exist for targets to act on). Independently testable once US1 done.
- **US3 (Phase 5)**: Depends on US1 (extends the same compose file) and benefits from US2's `make` targets. Independently testable.
- **Polish (Phase 6)**: After desired stories complete.

### Within Each User Story

- US1: T007 → T008/T009/T010/T011 (all edit `docker-compose.yml`, sequential) → T012 verify.
- US2: T013 (edit `Makefile`) → T014 verify.
- US3: T015 (`docker-compose.yml`) and T016 (`Makefile`) can proceed in parallel → T017 verify.

### Parallel Opportunities

- T001, T002 (different `.dockerignore` files) — parallel.
- T004, T005 (different Dockerfiles) — parallel.
- T015 and T016 (different files) — parallel.
- T018, T019 (docs vs. test run) — parallel.
- Note: all `docker-compose.yml` edits (T007–T011, T015) touch one file — keep sequential to avoid conflicts.

---

## Parallel Example: Phase 2 Foundational

```bash
# Build both images concurrently (different Dockerfiles):
Task: "Create backend/Dockerfile (multi-stage Go build)"
Task: "Create frontend/Dockerfile (node:24-alpine Vite dev server)"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Phase 1: Setup (`.dockerignore` ×2, vite server config).
2. Phase 2: Foundational (both Dockerfiles build).
3. Phase 3: US1 (compose file + inline env + healthcheck).
4. **STOP and VALIDATE**: `docker compose up`, smoke test. This is a demoable MVP.

### Incremental Delivery

1. Setup + Foundational → images build.
2. US1 → one-command stack works → demo (MVP).
3. US2 → `make` lifecycle targets → demo.
4. US3 → live frontend HMR + backend rebuild path → demo.

---

## Notes

- All configuration is inline in `docker-compose.yml` — **no `.env` files** (FR-014/015/016).
- Browser calls the host-mapped backend (`localhost:8080`); CORS must allow `http://localhost:5173` (research R1).
- Only one application source file changes: `frontend/vite.config.ts` (T003).
- Pin base image tags (`node:24-alpine`, `golang:1.25-alpine`) per Constitution III.
- Commit after each task or logical group; stop at any checkpoint to validate independently.
