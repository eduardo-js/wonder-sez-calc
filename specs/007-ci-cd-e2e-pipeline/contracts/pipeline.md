# Contract: CI Pipeline

The verifiable behavior the pipeline MUST satisfy.

## Triggers

- `pull_request` targeting `main` → full pipeline runs.
- `push` to `main` → full pipeline runs.
- New commit on an open PR/branch → in-progress run for that ref is cancelled (concurrency).

## Jobs & required checks

| Job        | Must run                                                              | Fails when                              |
|------------|----------------------------------------------------------------------|-----------------------------------------|
| `frontend` | lint, typecheck, unit tests, production build (in `frontend/`)       | any step non-zero                       |
| `backend`  | `go vet`, unit tests + coverage gate, `go build`                     | any step non-zero or coverage < 95%     |
| `e2e`      | start stack → wait healthy → Playwright calculation journey → teardown | stack unhealthy in time budget, or any test fails |

Overall workflow status = AND of all jobs. This status is the required merge gate (FR-004).

## Acceptance

- PR with a failing lint rule / unit test / coverage drop / build error → workflow **fails**, merge blocked (US1, SC-002).
- Clean PR → workflow **passes**, mergeable.
- `e2e` does not start unless `frontend` and `backend` both pass.

## Artifacts (on e2e failure)

- Playwright HTML report + traces + screenshots.
- `docker compose logs` dump.
- Uploaded with `if: always()` so they exist even when the job fails (FR-007, SC-004).

## Non-functional contract

- Permissions: `contents: read`; no secrets referenced → forked-PR runs safe (FR-009).
- Re-running a run on the same commit yields the same pass/fail decision (FR-010).
- Typical run completes < 15 min (SC-005).
- CI commands == local commands (npm scripts / go / `make test-backend`) → local-CI parity (SC-006).
