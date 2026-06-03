# Prompts

This project was built with AI tooling (Claude Code + [GitHub Spec Kit](https://github.com/github/spec-kit)).
Work proceeded feature-by-feature: each feature began with a natural-language **prompt**, which
Spec Kit captured verbatim as the **`**Input**` line** at the top of that feature's `spec.md`.
That spec then drove the downstream `/speckit-*` artifacts (`plan.md`, `tasks.md`, `research.md`,
`contracts/`, `quickstart.md`) and finally the implementation.

**The prompts are the spec inputs.** This file is a convenience index — the source of truth for
each prompt is the `**Input**` line inside the linked `spec.md`.

| # | Feature | Spec (prompt source) |
|---|---------|----------------------|
| 001 | Calculator UI barebones | [`specs/001-calculator-ui-barebones/spec.md`](specs/001-calculator-ui-barebones/spec.md) |
| 002 | Expand Go backend | [`specs/002-expand-go-backend/spec.md`](specs/002-expand-go-backend/spec.md) |
| 003 | Wire frontend ↔ backend | [`specs/003-wire-frontend-backend/spec.md`](specs/003-wire-frontend-backend/spec.md) |
| 004 | Large-number handling | [`specs/004-large-number-handling/spec.md`](specs/004-large-number-handling/spec.md) |
| 005 | Expression hints | [`specs/005-expression-hints/spec.md`](specs/005-expression-hints/spec.md) |
| 006 | Docker Compose setup | [`specs/006-docker-compose-setup/spec.md`](specs/006-docker-compose-setup/spec.md) |
| 007 | CI/CD + e2e pipeline | [`specs/007-ci-cd-e2e-pipeline/spec.md`](specs/007-ci-cd-e2e-pipeline/spec.md) |

---

## Verbatim feature prompts

### 001 — Calculator UI barebones
> create the barebones for a react project: it should have a clear and intuitive UI for a
> calculator app — (1) reusable components for each button, (2) tests for each component and how
> they render, (3) input validation, (4) responsive support; use shadcn; this project is a
> monorepo containing a react and golang project orchestrated via a Makefile, prepared in advance.

### 002 — Expand Go backend
> Let's expand the go backend: (1) add CORS for the frontend url, (2) add a framework for handling
> HTTP, gin preferred, (3) input validation — package or reflect?, (4) define project standards:
> stretchr/testify, prefer mocks, aim for >95% coverage, always prefer table-driven test
> structure, (5) ensure clear routing.

### 003 — Wire frontend ↔ backend
> let's wire the frontend and backend; 1. frontend should not make any calculations, these should
> be handled by backend 2. backend should receive user input, validate it and handle it correctly
> 2.1 if its a valid mathematical equation, calculate it and return as json 2.2 if its invalid
> input return an error with details 2.2.1 frontend should display alerts for errors 2.2.2 block
> multiple requests from frontend, add a spinner in the '=' sign while backend is processing
> information

### 004 — Large-number handling
> handle more edge cases on backend and frontend. Example:
> `-9223372036854775808 - 999999999999999999999999999999999999999999` returns
> `-9223372036854775808` (wrong). How to handle this properly? Add 'e' (scientific notation)
> special handling for larger scenarios.

### 005 — Expression hints
> let's improve frontend UI / lets add hints for current expression being evalued
>
> *(clarified: the hint shows the in-progress expression text only; the frontend does not
> evaluate it — evaluation is the backend's job, on equals.)*

### 006 — Docker Compose setup
> add docker compose for frontend and backend project, update make with docker commands

### 007 — CI/CD + e2e pipeline
> lets add CI CD and e2e tests to run on pipeline using github actions

---

## Follow-up prompts (post-007, no separate spec)

These were small iterations made directly against the running repo rather than new Spec Kit
features:

- Fix the local `make e2e` setup (Playwright `--with-deps` required root/apt and broke locally).
- Add e2e coverage for all calculator scenarios (12 Playwright tests).
- Update the 007 plan/contracts to match the implemented e2e scenarios; merge & push.
- Add frontend test coverage (Vitest v8) with a CI gate.
- Add README API examples, design-decisions, and this prompts index.
