# Development Guide

## Prerequisites

Install tools using [mise](https://mise.jdx.dev/):

```bash
mise install
```

- **[Go](https://go.dev)** ≥ 1.25
- **[just](https://just.systems)** — task runner

## Quick Start

```bash
git clone https://github.com/retr0h/mlb-sdk.git
cd mlb-sdk
go test ./...
```

## Layout

```
api/openapi.yaml      Hand-authored OpenAPI 3.0 spec
internal/gen/         Generated typed HTTP client (oapi-codegen).
                      NOT importable externally.
pkg/mlb/              Public, idiomatic Go SDK — the only thing you import
```

This is a library only — no `main.go`, no `cmd/`, no published binary.

## Common tasks

```bash
go test ./...              # run tests
go vet ./...               # vet
gofmt -l .                 # find unformatted files
go generate ./internal/gen # regenerate client.gen.go from api/openapi.yaml
```

Or via just:

```bash
just deps
just test
just generate
just ready          # fmt + vet + lint
```

## Regenerating the client

The generated client (`internal/gen/client.gen.go`) is checked in. Regenerate
whenever `api/openapi.yaml` changes:

```bash
go generate ./internal/gen
```

The directive in `internal/gen/generate.go` calls `oapi-codegen` via a `go tool`
reference, so no separate install step is required.

## Adding a new endpoint

> **Porting from [toddrob99/MLB-StatsAPI][toddrob99]?** Use the
> [`ingest-from-toddrob99`](../.claude/skills/ingest-from-toddrob99/) skill — it
> wraps the recipe below with the Python→OpenAPI translation table and handles
> the manifest/commit cadence automatically.

When wrapping a new MLB Stats API path, touch every one of these files in this
order. Skipping any of them leaves the public surface incomplete.

1. **`api/openapi.yaml`** — add the path entry under `paths:` and any new
   component schemas. Set an explicit `operationId`. Every nested object must be
   a named `components/schemas/...` reference (no inline `type: object`); see
   the spec authoring rules below.
2. **`go generate ./internal/gen`** — regenerate `client.gen.go` from the spec.
   Commit the regenerated file.
3. **`pkg/mlb/<name>_types.go`** — declare the public types only. No methods, no
   logic. Use idiomatic Go (`time.Time` for dates, typed enum constants for
   status fields, etc.).
4. **`pkg/mlb/<name>.go`** — add the `Client.<Method>(ctx, ...)` implementation,
   plus private `<name>FromGen` converters that map the generated layer's
   pointer-heavy types onto the public type. Wrap errors as
   `fmt.Errorf("mlb: <method>: %w", err)`. Map 404 to `ErrNotFound`. Map other
   non-200 to a wrapped `unexpected status` error.
5. **`pkg/mlb/<name>_test.go`** — one table-driven test per public function,
   with rows covering: happy path, 200 with empty/missing fields, 404, 5xx,
   malformed JSON, network failure (closed server). Coverage must stay at 100.0%
   — run `just go::test` to confirm.
6. **`examples/<name>.go`** — a runnable example program (one file per endpoint,
   all under `examples/` which is its own Go submodule with a `replace`
   directive pointing at the parent). Run with `go run examples/<name>.go`.
7. **`README.md`** — add a row to the `## ⚙️ Endpoints` table. Three columns:
   endpoint path, link to the pkg.go.dev anchor for the new `Client.<Method>`,
   link to the `examples/<name>.go` file. Add the `[d-<name>]` reference-style
   link footer next to the existing ones.
8. **`just ready`** — final gate. fmt + vet + lint + 100% coverage all green
   before committing.

## OpenAPI spec authoring

When adding endpoints or types to `api/openapi.yaml`:

- **Every nested object gets its own named schema** under `components/schemas/`.
  Inline `type: object` with `properties:` is forbidden in path responses and
  request bodies — `oapi-codegen` produces ugly anonymous nested struct types
  for inline objects, and clean top-level Go types for named schemas.
- Reference shared shapes via `$ref`.
- Use `additionalProperties: true` on response schemas so forward-compatible
  fields the MLB API adds don't break unmarshalling.
- Set explicit `operationId` on every path — it becomes the generated function
  name in `internal/gen`.

## Public surface authoring

`internal/gen` is an implementation detail. Consumers must import only
`pkg/mlb`. The whole point of the SDK is to encapsulate the MLB API's awkward
bits behind idiomatic Go.

### File organization

Every endpoint or domain concept gets three files:

```
pkg/mlb/<resource>_types.go    Exported types (no methods, no logic)
pkg/mlb/<resource>.go          Methods on Client + private conversion helpers
pkg/mlb/<resource>_test.go     Table-driven tests
```

`pkg/mlb/client.go` holds the `Client` type and the functional-options
constructor; `pkg/mlb/teams.go` holds typed identifier constants. Add new
domain-wide enums in their own `<concept>.go` file.

### Naming conventions

- Public type for a top-level resource: singular noun (`Boxscore`, `Game`).
- Per-side / per-team subtype: prefix with the parent (`BoxscoreTeam`).
- Query parameters struct: `<Resource>Query` (e.g., `ScheduleQuery`).
- Status / category enums: `<Concept>Type` or `<Concept>Status` typed string
  with grouped constants (`StatusFinal`, `StatusLive`, `StatusPreview`).
- Sentinel errors: `Err<Reason>` at package level (`ErrNotFound`).
- Private gen→public conversion helpers: `<resource>FromGen` taking the `*gen.X`
  pointer and returning the public type.

### Functional options

Every component that needs configuration uses functional options:

```go
type Option func(*config)
type config struct { /* internal */ }

func WithThing(v T) Option { return func(c *config) { c.thing = v } }
```

Defaults live inside the constructor; options override.

### Error handling

- Every `Client` method returns `(*T, error)` or `([]T, error)` and takes
  `context.Context` as the first argument.
- Wrap errors with the method name: `fmt.Errorf("mlb: <method>: %w", err)`.
- HTTP 404 maps to `ErrNotFound` (callers can `errors.Is`).
- Other non-200 maps to
  `fmt.Errorf("mlb: <method>: unexpected status %d", code)`.
- Public methods never panic; they always return errors.

### Conversion pattern

The generated layer makes every field `*T` because the spec uses
`additionalProperties: true`. The handwritten layer must:

1. Nil-check every pointer before deref.
2. Copy fields onto a public type with non-pointer fields where the value
   semantics are clear (`int` for counts, `string` for names, `time.Time` for
   dates).
3. Retain the underlying `raw *gen.X` on the public type when downstream helpers
   may need fields we have not yet promoted (see `BoxscoreTeam.raw`).

### Generic helpers

`pkg/mlb/schedule.go` declares `ptr[T any](v T) *T { return &v }` — used to
build `*T` request parameter values without scratch variables. Reuse this helper
rather than redeclaring it per file.

## Testing conventions

**Every public function and method MUST have a table-driven test.** One table
per function, with rows covering both the happy path and every failure mode the
function can produce. Failure rows belong in the same table as the happy row —
not in a separate test.

> **Anti-pattern (do not do this):** writing a separate one-off test function
> for a failure scenario. If you find yourself drafting a
> `TestClient_BoxscoreReturnsErrorOn500`, stop — instead add a row to the
> existing `TestClient_Boxscore` table. Each public function gets exactly
> **one** `Test*` function in the codebase. Reviewers should reject PRs that
> introduce additional one-off tests for the same function.

> **Anti-pattern:** asserting on `gen.X` types from a public-package test. Tests
> should assert on the public types we wrap (`Boxscore`, `Game`, `Play`, etc.),
> proving the wrapping actually happens. The only place a `gen.X` reference is
> acceptable in a test is when _constructing_ fake input for a private
> conversion helper — `boxscoreFromGen`, `playFromGen`, etc.

### Table shape

For pure functions:

```go
func TestParseDPCount(t *testing.T) {
    cases := []struct {
        name  string
        input string
        want  int
    }{
        {"empty", "", 0},
        {"single DP, no leading number", "(Smith-Jones).", 1},
        {"two DPs", "2 (Smith; Jones).", 2},
    }
    for _, c := range cases {
        t.Run(c.name, func(t *testing.T) { ... })
    }
}
```

For `Client` methods (HTTP-level), each row configures a per-case
`httptest.Server` and asserts on either the parsed result or the error:

```go
func TestClient_Schedule(t *testing.T) {
    cases := []struct {
        name       string
        query      ScheduleQuery
        respStatus int    // http status the fake server returns
        respBody   string // raw response body (set "" with respStatus 0 for net failure)
        wantLen    int    // expected len(result) on success
        wantErr    string // substring match; "" means expect nil error
        wantIs     error  // optional: errors.Is target (e.g. ErrNotFound)
    }{
        {name: "happy path", respStatus: 200, respBody: `{...}`, wantLen: 2},
        {name: "404 returns ErrNotFound", respStatus: 404, wantIs: ErrNotFound},
        {name: "5xx is wrapped", respStatus: 500, wantErr: "unexpected status 500"},
        {name: "malformed JSON", respStatus: 200, respBody: `not json`, wantErr: "schedule"},
        {name: "network failure", respStatus: 0, wantErr: "schedule"},
    }
    for _, c := range cases {
        t.Run(c.name, func(t *testing.T) { ... })
    }
}
```

### Required failure rows for HTTP methods

Every `Client` method's table MUST include rows covering:

| Row                  | Setup                                                                                        |
| -------------------- | -------------------------------------------------------------------------------------------- |
| Happy path           | `respStatus: 200`, expected body                                                             |
| Empty/missing fields | `respStatus: 200`, body with omitted optional fields → expect graceful zero values, no error |
| 404                  | `respStatus: 404` → expect `errors.Is(err, ErrNotFound)`                                     |
| 5xx                  | `respStatus: 500` → expect wrapped "unexpected status" error                                 |
| Malformed JSON       | `respStatus: 200`, body that is not valid JSON → expect wrapped error                        |
| Network failure      | server closed before request → expect wrapped error                                          |

Pure helpers (`parseDPCount`, `doublePlaysTurned`, `<resource>FromGen`) need
only the failure rows that apply to them — typically empty input, nil fields,
malformed input.

### Server-per-case pattern

Use a fresh `httptest.NewServer` inside each `t.Run` so that test cases are
isolated. To simulate a network failure, close the server before calling the
client. Do not share servers across rows.

### Test helpers

`strPtr`, `intPtr`, etc. live at the bottom of the test file in question. Don't
promote them across files unless three or more files end up duplicating the same
helper.

### Test naming

- Pure functions: `TestFunctionName`.
- Methods on a type: `TestType_Method` (Go convention; `go test` displays it
  nicely).
- Helper / unexported: `Test<helperName>`, lowercase first letter preserved.

## Commit messages

Follow [Conventional Commits](https://www.conventionalcommits.org/):

- **Subject line**: max 50 characters, imperative mood, capitalized, no period
- **Body**: wrap at 72 characters, separated from subject by a blank line
- **Format**: `type(scope): description`
- **Types**: `feat`, `fix`, `docs`, `style`, `refactor`, `perf`, `test`, `chore`
- **Scopes**: `api` (OpenAPI spec), `gen` (generated client), `mlb` (public
  SDK), `cli`, `docs`

[toddrob99]: https://github.com/toddrob99/MLB-StatsAPI
