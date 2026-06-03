# wonderlic-calc e2e

Playwright end-to-end tests. Expects the full stack already running.

## Prerequisites

- Docker
- Node 24

## Run locally

```sh
make docker-up           # from repo root — starts frontend:5173 + backend:8080
cd e2e && npm ci
npx playwright install chromium
npx playwright test
npx playwright show-report
make docker-down
```

`baseURL` is `http://localhost:5173`. Do not run `npx playwright test` before `make docker-up` or tests will fail to connect.
