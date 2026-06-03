# Contracts

Shared API contracts — the single source of truth for frontend/backend integration
(per the project constitution, Principle I & IV).

## Current contracts

- **Backend (operational)**: [`backend-openapi.yaml`](../specs/001-calculator-ui-barebones/contracts/backend-openapi.yaml)
  — `/healthz` and `/readyz` only. The backend is scaffolding-only in feature
  `001-calculator-ui-barebones`.
- **Frontend component contract**: [`frontend-component-contract.md`](../specs/001-calculator-ui-barebones/contracts/frontend-component-contract.md)

## Deferred

The frontend↔backend **calculation** contract is intentionally deferred to a later feature
(v1 computes client-side). New cross-boundary contracts land here when introduced.
