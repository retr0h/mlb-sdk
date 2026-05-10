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

The directive in `internal/gen/generate.go` calls `oapi-codegen` via a
`go tool` reference, so no separate install step is required.

## OpenAPI spec authoring

When adding endpoints or types to `api/openapi.yaml`:

- **Every nested object gets its own named schema** under
  `components/schemas/`. Inline `type: object` with `properties:` is
  forbidden in path responses and request bodies — `oapi-codegen` produces
  ugly anonymous nested struct types for inline objects, and clean top-level
  Go types for named schemas.
- Reference shared shapes via `$ref`.
- Use `additionalProperties: true` on response schemas so forward-compatible
  fields the MLB API adds don't break unmarshalling.
- Set explicit `operationId` on every path — it becomes the generated
  function name in `internal/gen`.

## Public surface authoring

`internal/gen` is an implementation detail. Consumers must import only
`pkg/mlb`. The whole point of the SDK is to encapsulate the MLB API's
awkward bits behind idiomatic Go:

- Convert API string dates to `time.Time` at the boundary.
- Expose typed identifiers (`mlb.LAD`, not `119`).
- Surface helper methods like `Boxscore.Team(LAD).DoublePlaysTurned()` that
  hide awkward parsing (e.g., per-game team double plays only appear in the
  free-text `info.FIELDING.DP` block — the SDK parses it).
- Public types live in `pkg/mlb/<resource>_types.go`; methods on those types
  live in `pkg/mlb/<resource>.go`.
- The generated `gen.Client` is wrapped by `pkg/mlb.Client` — never returned
  to callers.

## Testing pattern

Use `mlb.WithBaseURL(srv.URL)` in tests against an `httptest.Server` to
avoid hitting the live MLB API. See `pkg/mlb/boxscore_test.go` for an
example.

## Commit messages

Follow [Conventional Commits](https://www.conventionalcommits.org/):

- **Subject line**: max 50 characters, imperative mood, capitalized, no period
- **Body**: wrap at 72 characters, separated from subject by a blank line
- **Format**: `type(scope): description`
- **Types**: `feat`, `fix`, `docs`, `style`, `refactor`, `perf`, `test`, `chore`
- **Scopes**: `api` (OpenAPI spec), `gen` (generated client), `mlb` (public
  SDK), `cli`, `docs`
