# wonderlic-calc

This project is built with [GitHub Spec Kit](https://github.com/github/spec-kit) — a toolkit for **Spec-Driven Development**, where you describe *what* and *why* before *how*, and executable specs drive implementation.

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
