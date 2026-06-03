# Phase 1 Data Model: Expand Go Backend

This feature is infrastructure-level; "entities" are configuration and HTTP-contract types,
not persisted domain data.

## Config

Loaded once at startup from environment (FR-002).

| Field | Env var | Type | Default | Rules |
|-------|---------|------|---------|-------|
| Addr | `ADDR` | string | `:8080` | non-empty; `host:port` form |
| AllowedOrigins | `CORS_ALLOWED_ORIGINS` | []string | `["http://localhost:5173"]` | comma-split; each trimmed; empty entries dropped; no `*` when credentials enabled |
| ReadTimeout / WriteTimeout / IdleTimeout | (constants, retained) | duration | 5s / 10s / 120s | > 0 |

**Validation**: empty `CORS_ALLOWED_ORIGINS` → fall back to default (never an empty allow-list
that silently blocks dev). Malformed `ADDR` surfaces at server start.

## CORS Policy (derived from Config)

| Field | Value |
|-------|-------|
| AllowOrigins | `Config.AllowedOrigins` |
| AllowMethods | `GET, POST, PUT, PATCH, DELETE, OPTIONS` |
| AllowHeaders | `Origin, Content-Type, Accept, Authorization` |
| AllowCredentials | false (v1; no wildcard-with-credentials hazard) |
| MaxAge | 12h (preflight cache) |

Rule (FR-003): an origin not in `AllowOrigins` receives no `Access-Control-Allow-Origin` for it.

## Error Envelope (HTTP contract type)

Single shape for all error responses (FR-007).

```go
type Error struct {
    Code    string            `json:"code"`              // machine-readable
    Message string            `json:"message"`           // human-readable
    Fields  map[string]string `json:"fields,omitempty"`  // per-field validation detail
}
// wire: { "error": Error }
```

| Code | HTTP status | When |
|------|-------------|------|
| `validation_failed` | 400 | body/params fail binding or validation rules (FR-006) |
| `bad_request` | 400 | malformed/empty/oversized/non-JSON body |
| `not_found` | 404 | unknown route (`NoRoute`) (FR-008) |
| `method_not_allowed` | 405 | unsupported method on known path (`NoMethod`) (FR-008) |
| `internal` | 500 | recovered panic / unexpected error; no internals leaked |

## Validation Rule Set (pattern, not a concrete endpoint)

Bound to request DTOs via struct tags; evaluated by `httpx.BindJSON` before any handler logic
(FR-006). No business DTOs ship in this feature — the pattern is proven with an in-test DTO.

Example pattern:

```go
type exampleReq struct {
    Name   string  `json:"name"   binding:"required"`
    Amount float64 `json:"amount" binding:"required,gt=0"`
    Mode   string  `json:"mode"   binding:"omitempty,oneof=a b c"`
}
```

| Failure | Maps to |
|---------|---------|
| missing required field | `validation_failed` + `fields[name]="required"` |
| wrong type / unparsable JSON | `bad_request` |
| out-of-range (`gt`, `min`, `max`) | `validation_failed` + field rule |
| disallowed enum (`oneof`) | `validation_failed` + field rule |
| valid | passes to handler |

## Migrated Probe Responses (unchanged contract — FR-009)

| Endpoint | Status | Body |
|----------|--------|------|
| `GET /healthz` | 200 | `{"status":"ok"}` |
| `GET /readyz` | 200 | `{"status":"ready"}` |

## Route Registry (central — FR-005)

| Method | Path | Group | Handler |
|--------|------|-------|---------|
| GET | `/healthz` | top-level | health.Healthz |
| GET | `/readyz` | top-level | health.Readyz |
| (future) | `/api/v1/...` | `/api/v1` | later features |
| ANY | unknown path | — | NoRoute → `not_found` |
| (bad method) | known path | — | NoMethod → `method_not_allowed` |
