# Quickstart: Wire Frontend and Backend

**Feature**: 003-wire-frontend-backend

## Run

```bash
# Backend (terminal 1)
cd backend && go run ./cmd/server          # :8080

# Frontend (terminal 2)
cd frontend && npm run dev                  # :5173, proxies/calls :8080
```

Optional frontend env: `VITE_API_BASE_URL=http://localhost:8080` (default).

## Verify the contract (backend)

```bash
# valid → 200
curl -s localhost:8080/api/v1/calculate -H 'content-type: application/json' \
  -d '{"expression":"2 + 3 * 4"}'
# {"result":"14","expression":"2 + 3 * 4"}

# malformed → 400 validation_failed
curl -s localhost:8080/api/v1/calculate -H 'content-type: application/json' \
  -d '{"expression":"5 + * 2"}'

# divide by zero → 422 calculation_error
curl -s localhost:8080/api/v1/calculate -H 'content-type: application/json' \
  -d '{"expression":"1/0"}'

# empty → 400 validation_failed
curl -s localhost:8080/api/v1/calculate -H 'content-type: application/json' \
  -d '{"expression":""}'
```

## Verify in the UI

1. Enter `2 + 3`, press `=` → spinner shows on `=`, then `5` displays (from backend).
2. While loading, repeated `=` presses issue no extra requests.
3. Trigger `1 / 0` → red alert with backend message; display unchanged.
4. Stop the backend, press `=` → connectivity alert; `=` control restored (not stuck).

## Tests

```bash
cd backend && go test ./... -cover      # ≥95% (constitution); evaluator table-driven
cd frontend && npm test                 # api client + hook lifecycle + alert/spinner
```

## Acceptance trace

| Spec | Check |
|------|-------|
| FR-001 | `frontend/src/lib/calculator.ts` removed; no client arithmetic |
| FR-002..007 | `POST /api/v1/calculate` returns result JSON |
| FR-005/006/008 | error envelope → alert region |
| FR-009/010/011 | in-flight lock + `=` spinner restore |
| FR-012 | abort/timeout → alert + restore |
| FR-013 | matches `contracts/calculate-openapi.yaml` |
| FR-014 | structured request logs (existing middleware) |
