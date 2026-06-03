# Quickstart: Expand Go Backend

Prereqs: Go 1.23+, Make. All commands from repo root unless noted.

## Add dependencies (backend)

```bash
cd backend
go get github.com/gin-gonic/gin
go get github.com/gin-contrib/cors
go get github.com/go-playground/validator/v10
go get github.com/stretchr/testify
go mod tidy
```

## Run

```bash
make run-backend                     # :8080 by default
# configure CORS + addr per environment:
ADDR=":8080" CORS_ALLOWED_ORIGINS="http://localhost:5173,https://app.example.com" make run-backend
```

## Verify CORS (US1)

```bash
# Preflight from allowed origin → headers permit it
curl -i -X OPTIONS http://localhost:8080/api/v1/ping \
  -H "Origin: http://localhost:5173" \
  -H "Access-Control-Request-Method: POST"
# expect: Access-Control-Allow-Origin: http://localhost:5173

# Disallowed origin → not granted
curl -i -X OPTIONS http://localhost:8080/api/v1/ping \
  -H "Origin: http://evil.example" \
  -H "Access-Control-Request-Method: POST"
# expect: no Access-Control-Allow-Origin for that origin
```

## Verify routing + error envelope (US3)

```bash
curl -s http://localhost:8080/healthz   # {"status":"ok"}   (unchanged)
curl -s http://localhost:8080/readyz    # {"status":"ready"} (unchanged)
curl -i http://localhost:8080/nope      # 404 {"error":{"code":"not_found",...}}
curl -i -X DELETE http://localhost:8080/healthz  # 405 {"error":{"code":"method_not_allowed",...}}
```

## Tests + coverage gate (US3, FR-011/012)

```bash
make test-backend          # runs go test with coverage; fails if total < 95%
cd backend && go test ./... -cover            # quick local total
cd backend && go test ./... -coverprofile=coverage.out && go tool cover -func=coverage.out | tail -1
cd backend && go tool cover -html=coverage.out   # inspect gaps
```

## Conventions for adding a route (US3)

1. Register it in `internal/server/router.go` — top-level for ops, `/api/v1` group for app routes.
2. Define a request DTO with `binding:"..."` tags; bind via `httpx.BindJSON`.
3. On error, return through `httpx.WriteError` (or let validation/recovery middleware do it) so
   every response uses the standard envelope.
4. Write **table-driven** tests with `testify`; mock collaborators via `testify/mock`; keep
   total coverage ≥95%.

## Definition of done (maps to Success Criteria)

- [ ] Frontend origin can call the API in-browser; non-allowed origin denied (SC-001)
- [ ] All invalid-input classes rejected with standard envelope, no crash (SC-002)
- [ ] Validation, 404, and 405 share one error shape (SC-003)
- [ ] All routes registered centrally; new route added by convention (SC-004)
- [ ] `make test-backend` green, table-driven, coverage ≥95% (SC-005)
- [ ] `/healthz` + `/readyz` behavior unchanged (SC-006)
```
