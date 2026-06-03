# wonderlic-calc

This project is built with [GitHub Spec Kit](https://github.com/github/spec-kit) — a toolkit for **Spec-Driven Development**, where you describe *what* and *why* before *how*, and executable specs drive implementation.

## Development

Monorepo: `frontend/` (React + TypeScript + Vite + Tailwind + shadcn/ui) and `backend/`
(Go, stdlib `net/http`), orchestrated by a root `Makefile`.

**Prerequisites**: Node **20+** (a `.nvmrc` pins the latest LTS — run `nvm use`) and Go **1.23+**.

```bash
make install        # install frontend + backend deps
make dev            # run frontend (Vite) and backend (:8080) together
make test           # run frontend (Vitest) + backend (go test) suites
make lint           # eslint + go vet
make fmt            # prettier + gofmt
make build          # production build of both workspaces
make help           # list all targets
```

Quick check: `curl localhost:8080/healthz` → `{"status":"ok"}`. The frontend never
evaluates expressions itself — it sends them to the backend and renders the result (see
[API](#api)). Full verification steps: `specs/001-calculator-ui-barebones/quickstart.md`.

### Run with Docker

No Node or Go required — only Docker Engine + Compose v2:

```bash
make docker-up        # build + start frontend (:5173) and backend (:8080), detached
make docker-logs      # follow combined logs
make docker-rebuild   # rebuild images + restart (apply code changes)
make docker-down      # stop + remove containers
make docker-build     # build images without starting
```

Then open http://localhost:5173. Configuration is inline in `docker-compose.yml` (no `.env`
file). Editing `frontend/src/**` hot-reloads live; for backend changes run `make docker-rebuild`.
See `specs/006-docker-compose-setup/quickstart.md` for details.

### Continuous Integration

GitHub Actions (`.github/workflows/ci.yml`) gates every PR and push to `main`:

| Job        | Runs                                                            |
|------------|-----------------------------------------------------------------|
| `frontend` | lint, typecheck, unit tests, build                              |
| `backend`  | `go vet`, tests + 95% coverage gate, build                      |
| `e2e`      | boots the Docker stack, runs Playwright against the live UI (gated on the two above) |

The combined status is the required merge check. E2E uses the same commands as local — run them with:

```bash
make e2e              # start stack, run Playwright, tear down
make test-e2e         # run e2e against an already-running stack
```

E2E lives in `e2e/` (Playwright, Chromium). See `specs/007-ci-cd-e2e-pipeline/quickstart.md`.

## API

The backend exposes one calculation endpoint plus a health probe.

| Method | Path | Purpose |
|--------|------|---------|
| `POST` | `/api/v1/calculate` | Evaluate an arithmetic expression |
| `GET`  | `/healthz` | Liveness probe → `{"status":"ok"}` |

Request body: `{ "expression": string }` (max 256 chars; digits, `. + - * / ( )` and spaces).
Results are returned as **strings** to preserve precision and scientific notation for large values.

**Success** — `200 OK`:

```bash
curl -s localhost:8080/api/v1/calculate \
  -H 'Content-Type: application/json' \
  -d '{"expression":"2 + 2 * 3"}'
# → {"result":"8","expression":"2 + 2 * 3"}
```

**Division by zero** — `422 Unprocessable Entity`:

```bash
curl -s localhost:8080/api/v1/calculate \
  -H 'Content-Type: application/json' \
  -d '{"expression":"5 / 0"}'
# → {"error":{"code":"calculation_error","message":"division by zero"}}
```

**Invalid expression** — `400 Bad Request`:

```bash
curl -s localhost:8080/api/v1/calculate \
  -H 'Content-Type: application/json' \
  -d '{"expression":"2 +"}'
# → {"error":{"code":"validation_failed","message":"invalid expression",
#      "fields":{"expression":"invalid expression: not enough operands for \"+\""}}}
```

All errors share the envelope `{"error":{"code","message","fields?"}}`. Codes:
`validation_failed` (400), `calculation_error` (422, e.g. divide-by-zero / non-finite),
`bad_request`, `not_found`, `method_not_allowed`, `internal`.

## Design decisions & assumptions

- **Backend owns all evaluation.** The frontend sends the raw expression and renders the
  returned string; it never computes (no client-side `eval`). This keeps one source of
  arithmetic truth, avoids JS float drift, and makes the API the contract both sides test against.
- **Results are strings, not numbers.** Large results use scientific notation (`1e+42`) and
  precision is formatted server-side, so a JSON number would lose fidelity.
- **Custom evaluator, no `eval`.** The backend tokenizes and evaluates with explicit operator
  precedence (`backend/internal/calc/evaluator.go`) — safe against injection and fully testable.
- **Typed error envelope.** Sentinel errors (`ErrDivideByZero`, `ErrInvalidExpression`,
  `ErrNonFinite`) map to stable wire `code`s so the UI can react without string-matching messages.
- **Scope.** Four operations (`+ − × ÷`) only; exponentiation / square root / percentage were
  left out as optional per the brief (correctness & clarity over extra features).
- **Spec-driven.** Every feature began as a spec under `specs/` before code — see [Prompts](#prompts).
- **Assumptions.** Single-user local/demo use; no auth, persistence, or rate limiting; expressions
  capped at 256 chars; CORS allows the dev frontend origin only.

## Prompts

This project was built with AI tooling (Claude Code + GitHub Spec Kit). Every feature started
from a natural-language prompt captured verbatim as the **`**Input**` line** at the top of each
feature's `spec.md`. The prompts, in order, live in:

| Feature | Prompt file |
|---------|-------------|
| 001 calculator UI barebones | `specs/001-calculator-ui-barebones/spec.md` |
| 002 expand Go backend | `specs/002-expand-go-backend/spec.md` |
| 003 wire frontend ↔ backend | `specs/003-wire-frontend-backend/spec.md` |
| 004 large-number handling | `specs/004-large-number-handling/spec.md` |
| 005 expression hints | `specs/005-expression-hints/spec.md` |
| 006 docker compose setup | `specs/006-docker-compose-setup/spec.md` |
| 007 CI/CD + e2e pipeline | `specs/007-ci-cd-e2e-pipeline/spec.md` |

`PROMPTS.md` collects the same inputs in one place. Each `spec.md` also records the downstream
`/speckit-*` artifacts (`plan.md`, `tasks.md`, …) generated from that prompt.

## How it works

Spec Kit organizes work into a sequence of artifacts per feature:

```
constitution → specify → clarify → plan → tasks → implement
```

Each step produces a document (`spec.md`, `plan.md`, `tasks.md`, …) that feeds the next.

## Interacting with it

In Claude Code, run the `speckit-*` skills (slash commands) in order:

| Command | Purpose |
|---|---|
| `/speckit-constitution` | Define project principles & constraints |
| `/speckit-specify` | Create a feature spec from a natural-language description |
| `/speckit-clarify` | Answer targeted questions to remove ambiguity in the spec |
| `/speckit-plan` | Generate the technical implementation plan |
| `/speckit-tasks` | Break the plan into ordered, actionable tasks |
| `/speckit-analyze` | Cross-check spec ↔ plan ↔ tasks for consistency |
| `/speckit-checklist` | Generate a custom review checklist |
| `/speckit-implement` | Execute the tasks |
| `/speckit-taskstoissues` | Export tasks as GitHub issues |

### Typical flow

1. `/speckit-constitution` — set the ground rules (once per project).
2. `/speckit-specify "Calculate Wonderlic score from raw answers"` — describe the feature.
3. `/speckit-clarify` — resolve open questions.
4. `/speckit-plan` — pick the stack and design.
5. `/speckit-tasks` — produce the task list.
6. `/speckit-implement` — build it.

## Project subagents

Three specialized subagents (in `.claude/agents/`) own different parts of the monorepo:

| Agent | Domain | Edits code? |
|---|---|---|
| `project-architect` | Whole-repo architecture, frontend/backend boundaries, API contract, cross-cutting design | No (read-only) |
| `frontend-react` | React + TypeScript frontend (`frontend/`) | Yes |
| `backend-go` | Go backend (`backend/`) | Yes |

### Triggering agents during the spec workflow

Invoke an agent with the `@` mention or by asking Claude Code to "use the `<agent>` subagent".
Match the agent to the phase of the spec:

| Spec phase | Trigger | Why |
|---|---|---|
| After `/speckit-specify` | `@project-architect` review the spec | Validate boundaries & contract surface before planning |
| During `/speckit-plan` | `@project-architect` design the plan | Owns cross-cutting design and the API contract |
| During `/speckit-implement` (frontend tasks) | `@frontend-react` implement these tasks | Frontend specialist, test-first React/TS |
| During `/speckit-implement` (backend tasks) | `@backend-go` implement these tasks | Backend specialist, table-driven Go tests |
| Code review (either stack) | `@frontend-react` / `@backend-go` review the diff | Stack-specific quality + constitution checks |

Per-spec rule of thumb:

1. `/speckit-specify` → then `@project-architect` to sanity-check scope and boundaries.
2. `/speckit-plan` → drive with `@project-architect`; it sets the contract the other two implement against.
3. `/speckit-tasks` → split tasks by stack (`frontend/` vs `backend/`).
4. `/speckit-implement` → route each task: frontend tasks to `@frontend-react`, backend tasks to `@backend-go`.
5. Review → the matching stack agent reviews its own diff against the [constitution](.specify/memory/constitution.md).

All agents enforce the project [constitution](.specify/memory/constitution.md): test-first, clear boundaries, contract-driven integration.

## Layout

- `.specify/` — templates, scripts, workflows, and project memory (incl. the constitution).
- `.claude/skills/` — the Spec Kit skills exposed as slash commands.
- `.claude/agents/` — the project subagents (`project-architect`, `frontend-react`, `backend-go`).
- `CLAUDE.md` — instructions for the coding agent.

## Reference

- Spec Kit: https://github.com/github/spec-kit
