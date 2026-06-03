---

description: "Task list for 007-ci-cd-e2e-pipeline"
---

# Tasks: CI & E2E Pipeline

**Input**: Design documents from `/specs/007-ci-cd-e2e-pipeline/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/{pipeline,e2e-scenarios}.md, quickstart.md

**Tests**: The e2e Playwright suite IS a requested test artifact (US2) and is included below. No unit
tests are added for the workflow YAML (no unit-testable logic); the pipeline is verified by execution.
Existing app unit tests are unchanged.

**Organization**: Tasks grouped by user story (US1 CI gate → US2 e2e) for independent delivery.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: US1 / US2
- Exact file paths included

## Path Conventions

Monorepo: `.github/workflows/` (CI), `e2e/` (Playwright), `frontend/`, `backend/`, root `Makefile`.

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Workflow skeleton shared by all jobs.

- [ ] T001 Create `.github/workflows/ci.yml` skeleton: `name: CI`; triggers `pull_request` + `push` on `main`; `permissions: { contents: read }`; `concurrency: { group: ci-${{ github.ref }}, cancel-in-progress: true }`; empty `jobs:` map
- [ ] T002 [P] Add CI artifact ignores to `.gitignore`: `e2e/node_modules/`, `e2e/playwright-report/`, `e2e/test-results/`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: None beyond Setup — both stories build on the `ci.yml` skeleton (T001). No additional blocking work.

(No foundational tasks; proceed to user stories once T001 exists.)

---

## Phase 3: User Story 1 - Quality gate on every change (Priority: P1) 🎯 MVP

**Goal**: Every PR/push to `main` runs lint, type-check, unit tests, coverage gate, and build for both workspaces and reports a blocking pass/fail status.

**Independent Test**: Open a PR with a deliberately failing lint/test change → workflow fails, merge blocked. Clean PR → passes.

### Implementation for User Story 1

- [ ] T003 [US1] In `.github/workflows/ci.yml`, add the `frontend` job (`runs-on: ubuntu-latest`): `actions/checkout@v4`; `actions/setup-node@v4` with `node-version-file: .nvmrc` + `cache: npm` (cache-dependency-path `frontend/package-lock.json`); steps `npm ci`, `npm run lint`, `npm run typecheck`, `npm run test`, `npm run build` — all run in `frontend/`
- [ ] T004 [US1] In `.github/workflows/ci.yml`, add the `backend` job (`runs-on: ubuntu-latest`): `actions/checkout@v4`; `actions/setup-go@v5` with `go-version-file: backend/go.mod`; steps `go vet ./...`, `make test-backend` (coverage gate ≥ 95%), `go build ./cmd/server` — backend steps run in `backend/` (use repo-root `make` for the coverage target)
- [ ] T005 [US1] Verify gating: push a throwaway commit that breaks a frontend lint rule and a backend test on a PR branch → confirm `frontend`/`backend` jobs fail and the PR shows a red required check; revert → confirm green

**Checkpoint**: MVP — CI blocks bad changes on `main` PRs.

---

## Phase 4: User Story 2 - End-to-end verification of the running stack (Priority: P2)

**Goal**: Pipeline boots the full stack and drives the calculator through a real browser, asserting a correct result; failures upload diagnostic artifacts.

**Independent Test**: `make e2e` locally passes the calculation journey; in CI the `e2e` job runs only after both gates pass and uploads a report/trace on failure.

### Implementation for User Story 2

- [ ] T006 [P] [US2] Create `e2e/package.json`: name `wonderlic-calc-e2e`, private, `devDependencies` `@playwright/test`, script `"test": "playwright test"`
- [ ] T007 [P] [US2] Create `e2e/playwright.config.ts`: `use.baseURL=http://localhost:5173`, Chromium project, `trace: 'on-first-retry'`, `screenshot: 'only-on-failure'`, `retries: process.env.CI ? 2 : 0`, `reporter: 'html'`
- [ ] T008 [US2] Create `e2e/tests/calculator.spec.ts`: 12 scenarios via `getByRole('button', { name })` against the live stack — core round trip `2 + 2 =` → `4`, all four operators, decimals, chained ops, clear/reset, expression hint (`role="status"`), division-by-zero alert (`role="alert"`) + recovery, and large-number scientific-notation round trip (see `contracts/e2e-scenarios.md`)
- [ ] T009 [P] [US2] Create `e2e/README.md`: how to run locally against the compose stack (`make docker-up` → `npm ci` → `npx playwright install chromium` → `npx playwright test`)
- [ ] T010 [US2] In `.github/workflows/ci.yml`, add the `e2e` job: `needs: [frontend, backend]`; checkout; `setup-node` (`.nvmrc`, cache `e2e/package-lock.json`); `npm ci` in `e2e/`; `npx playwright install --with-deps chromium`; `docker compose up -d --build`; poll `http://localhost:8080/healthz` with bounded timeout (fail-fast + print `docker compose logs` on timeout); `npx playwright test` in `e2e/`; `actions/upload-artifact@v4` with `if: always()` for `e2e/playwright-report` + compose logs; `docker compose down -v` with `if: always()`
- [ ] T011 [US2] In `Makefile`, add `test-e2e` (`cd e2e && npm ci && npx playwright test`) and `e2e` (`docker compose up -d --build` → wait healthz → `$(MAKE) test-e2e` → `docker compose down -v`) targets, each `## `-documented, added to `.PHONY`
- [ ] T012 [US2] Verify locally: `make e2e` brings up the stack, runs Playwright, the `2+2=4` scenario passes, stack torn down

**Checkpoint**: E2E validates the real frontend↔backend contract in CI and locally.

---

## Phase 5: Polish & Cross-Cutting Concerns

- [ ] T013 [P] Add a CI status badge and a short "Continuous Integration" section to `README.md` (jobs, what gates merge, how to run e2e)
- [ ] T014 [P] Confirm `make test` (frontend + backend) and the full local flow still pass; no existing target regressed
- [ ] T015 Verify the full pipeline end-to-end on a real PR to `main`: all three jobs run, e2e gated on the two quality jobs, concurrency cancels superseded runs, total < 15 min

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: T001 (ci.yml skeleton) blocks all job tasks; T002 independent.
- **US1 (Phase 3)**: needs T001. MVP.
- **US2 (Phase 4)**: needs T001; the `e2e` job (T010) references `needs: [frontend, backend]` so it depends on US1's jobs existing. E2E package tasks (T006–T009) are independent of US1.
- **Polish (Phase 5)**: after US1/US2.

### Within Each User Story

- US1: T003 and T004 edit the same `ci.yml` → sequential; then T005 verify.
- US2: T006/T007/T009 (separate files) parallel; T008 needs T007; T010 edits `ci.yml` (sequential w/ other ci.yml edits); T011 edits `Makefile`; T012 verify.

### Parallel Opportunities

- T002 ∥ T001-dependent work (different file).
- T006 ∥ T007 ∥ T009 (distinct e2e files).
- T013 ∥ T014 (docs vs. test run).
- **Cross-story fan-out**: the e2e package (T006–T009, in `e2e/`) and the CI jobs (T003/T004, in `ci.yml`) touch disjoint files and can be built by separate agents concurrently; only the `ci.yml` edits among themselves (T001, T003, T004, T010) must serialize.

---

## Parallel Example: subagent fan-out

```text
Agent A (frontend-react): e2e/ Playwright package — T006, T007, T008, T009
Agent B (general/devops):  .github/workflows/ci.yml — T001, T003, T004, T010 (sequential within agent)
Main:                      Makefile T011, .gitignore T002, README T013
```

---

## Implementation Strategy

### MVP First (US1)

1. T001 skeleton → T003/T004 quality jobs → T005 verify. CI gates merges. Demoable MVP.

### Incremental Delivery

1. Setup → US1 (CI gate) → demo.
2. US2 (e2e) → demo.
3. Polish (badge, parity check, full PR validation).

---

## Notes

- CI runs the same commands as local (npm scripts / go / `make test-backend`) → parity (SC-006).
- No secrets used; `permissions: contents: read` → forked-PR runs safe (FR-009).
- E2E reuses the feature-006 Docker Compose stack and `/healthz` readiness.
- Gate branch is `main` (this branch will be pushed as the default `main`).
- Pin all action versions (Constitution III).
