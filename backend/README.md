# wonderlic-calc backend

Go HTTP service (Gin). Module: `github.com/wonderlic-calc/backend`.

## Layout

```
cmd/server/main.go        thin entrypoint: logger → config.Load() → server.Run()
internal/config/          env config (Addr, AllowedOrigins, timeouts)
internal/server/          central router (NewRouter) + serve/shutdown (Run)
internal/middleware/      Recovery, RequestLogger (slog JSON), CORS
internal/httpx/           Error envelope, WriteError, BindJSON, NotFound/MethodNotAllowed
internal/health/          /healthz, /readyz probes
```

## Run

```bash
make run-backend                       # :8080
ADDR=":8080" CORS_ALLOWED_ORIGINS="http://localhost:5173" make run-backend
```

| Env var | Default | Meaning |
|---------|---------|---------|
| `ADDR` | `:8080` | listen address |
| `CORS_ALLOWED_ORIGINS` | `http://localhost:5173` | comma-separated allowed browser origins |

## Routing conventions

- **All routes are registered centrally** in `internal/server/router.go` (`NewRouter`). There is
  no route registration scattered across packages.
- **Middleware chain** (applied to every route): `Recovery → RequestLogger → CORS`.
- **Operational probes** (`/healthz`, `/readyz`) live at the top level.
- **Application routes** go under the versioned group `/api/v1` (`r.Group("/api/v1")`).
- **Unknown route** → `404` and **wrong method** → `405`, both via `NoRoute`/`NoMethod` using the
  standard error envelope.

### Adding a route

1. Define a request DTO with `binding:"..."` tags (go-playground/validator):
   ```go
   type CreateThingRequest struct {
       Name   string  `json:"name"   binding:"required"`
       Amount float64 `json:"amount" binding:"required,gt=0"`
   }
   ```
2. In the handler, bind + validate with `httpx.BindJSON` (it writes the error envelope and
   returns `false` on failure):
   ```go
   func createThing(c *gin.Context) {
       var req CreateThingRequest
       if !httpx.BindJSON(c, &req) { return }
       // ... business logic ...
   }
   ```
3. Register it inside the `/api/v1` group in `NewRouter`.
4. Return errors via `httpx.WriteError(c, status, code, msg, fields)` so every response shares
   one shape.

## Error envelope

All client/server errors return:

```json
{ "error": { "code": "validation_failed", "message": "...", "fields": { "amount": "gt" } } }
```

Codes: `validation_failed` (400), `bad_request` (400), `not_found` (404),
`method_not_allowed` (405), `internal` (500). `fields` appears only for validation errors.

## Testing standards

- **Framework**: `stretchr/testify` (`require`/`assert`, `mock` for collaborators).
- **Structure**: table-driven by default; `gin.SetMode(gin.TestMode)` + `net/http/httptest`.
- **Coverage**: **≥95%** over `internal/...`, enforced by `make test-backend`.

```bash
make test-backend     # tests + coverage gate
make cover-backend     # HTML coverage report
```

`cmd/server/main.go` is intentionally trivial (signal wiring) and excluded from the coverage
gate; the serve/shutdown logic lives in `server.Run`, which is tested.
