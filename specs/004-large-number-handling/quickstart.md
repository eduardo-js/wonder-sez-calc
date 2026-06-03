# Quickstart: Large-Number & Edge-Case Handling

## Prereqs

- Toolchain via `make` (Go at `~/.local/go/bin`, Node 24 via `.nvmrc`).

## The change (one function)

`backend/internal/calc/evaluator.go` → `formatResult`:

```go
const maxExactInt = 1 << 53 // 9007199254740992: largest exactly-representable integer in float64

func formatResult(v float64) string {
    if v == math.Trunc(v) && !math.IsInf(v, 0) && math.Abs(v) <= maxExactInt {
        return strconv.FormatInt(int64(v), 10)
    }
    // Large or non-integer (incl. very small) → scientific notation, ≤12 sig digits.
    return strconv.FormatFloat(v, 'g', 12, 64)
}
```

## Verify (TDD order)

1. **Red** — add the regression case to `evaluator_test.go` and run it; it must fail today:

   ```sh
   make test-backend   # or: go test ./internal/calc/...
   ```

   Reported input: `-9223372036854775808 - 999999999999999999999999999999999999999999`
   - Before fix: result `-9223372036854775808` ❌
   - After fix: result `-1e+42` (scientific notation) ✅, and asserted `!= "-9223372036854775808"`.

2. **Green** — apply the `formatResult` change; rerun tests.

3. **Coverage / lint**:

   ```sh
   make test-backend     # keep ≥95%
   make lint     # gofmt/vet/golangci-lint
   ```

4. **Frontend** — confirm the display renders an `e`-notation result and error string:

   ```sh
   make test-frontend
   ```

## Manual check (optional)

```sh
curl -s localhost:8080/api/v1/calculate \
  -H 'content-type: application/json' \
  -d '{"expression":"-9223372036854775808 - 999999999999999999999999999999999999999999"}'
# => {"result":"-1e+42","expression":"..."}
```

## Test cases to add (table-driven)

| Expression | Expected result | Why |
|-----------|-----------------|-----|
| `-9223372036854775808 - 999999999999999999999999999999999999999999` | `-1e+42` (and `!= -9223372036854775808`) | Reported regression (SC-001) |
| `9007199254740992` (2⁵³) | `9007199254740992` (plain) | Boundary stays plain (FR-003) |
| `9007199254740992 + 9007199254740992` | scientific (`1.80143985095e+16`) | Just past boundary → scientific (FR-002) |
| `1e40 * 2` | `2e+40` | Large positive scientific |
| `1 / 1e20` | `1e-20` | Very small non-zero, not `0` (FR-005) |
| `1e400` | error `result is not a finite number` | Over-large literal → explicit error (FR-007) |
| `1 / 0` | error `division by zero` | Unchanged error path |
| `12 + 7 * 3` | `33` | No regression for normal values (SC-004) |
