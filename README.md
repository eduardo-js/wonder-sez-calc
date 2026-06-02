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

## Layout

- `.specify/` — templates, scripts, workflows, and project memory used by Spec Kit.
- `.claude/skills/` — the Spec Kit skills exposed as slash commands.
- `CLAUDE.md` — instructions for the coding agent.

## Reference

- Spec Kit: https://github.com/github/spec-kit
