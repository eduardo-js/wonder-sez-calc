# Contracts

Shared API contracts — the single source of truth for frontend/backend integration
(per the project constitution, Principle I & IV).

## Current contracts

- **Backend (current)**: [`backend-openapi.yaml`](../specs/002-expand-go-backend/contracts/backend-openapi.yaml)
  — v0.2.0: preserved `/healthz` + `/readyz`, standard `Error` envelope, versioned
  `/api/v1` group, and CORS policy. Supersedes the 001 scaffolding contract.
- **Backend (history)**: [`001 backend-openapi.yaml`](../specs/001-calculator-ui-barebones/contracts/backend-openapi.yaml)
  — v0.1.0 scaffolding (`/healthz`, `/readyz` only).
- **Frontend component contract**: [`frontend-component-contract.md`](../specs/001-calculator-ui-barebones/contracts/frontend-component-contract.md)

## Deferred

The frontend↔backend **calculation** contract is intentionally deferred to a later feature
(v1 computes client-side). New cross-boundary contracts land here when introduced.
