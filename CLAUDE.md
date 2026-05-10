# mlb-sdk

Idiomatic Go SDK and CLI for the public MLB Stats API hosted at
`statsapi.mlb.com`. MLB does not publish an OpenAPI spec — we author one in
this repo and generate a typed client from it via
[`oapi-codegen`](https://github.com/oapi-codegen/oapi-codegen).

## Layout

```
api/openapi.yaml      Hand-authored OpenAPI 3.0 spec
internal/gen/         Generated client. NEVER importable externally.
pkg/mlb/              Public, idiomatic SDK surface (the only thing users import)
```

This is a library only. There is no `main.go`, no `cmd/`, and no published
binary. Examples may live under `examples/` if useful for documentation.

## Commands

```sh
just deps       # mise + go modules
just generate   # regenerate internal/gen/client.gen.go from api/openapi.yaml
just test
just ready      # fmt + vet + lint
```

## Conventions (mandatory)

### OpenAPI spec — favor named components

The Go client is **only as ergonomic as the spec is**. `oapi-codegen` produces
named top-level Go types for `components/schemas/Foo` references and ugly
anonymous nested struct types for inline `type: object` properties. Therefore:

- **Every nested object gets its own named schema under `components/schemas/`.**
  Inline `type: object` with `properties:` is forbidden in path responses
  and request bodies.
- Reference shared shapes via `$ref` (e.g., `Team`, `SideScoreboard`,
  `DisplayLabel`).
- Use `additionalProperties: true` on response schemas so forward-compatible
  fields the MLB API adds don't break unmarshalling.
- Set explicit `operationId` on every path — that becomes the generated
  function name in `internal/gen`.
- The spec covers only the subset of endpoints the SDK exposes. It is **not**
  a comprehensive description of the MLB Stats API.

### Public surface — hide API quirks behind idiomatic methods

`internal/gen` is an implementation detail. Consumers must import only
`pkg/mlb`. The whole point of this SDK is to encapsulate the MLB API's
awkward bits behind idiomatic Go:

- Convert API string dates to `time.Time` at the boundary.
- Expose typed team identifiers (`mlb.LAD`, not `119`).
- Surface methods like `Boxscore.Team(LAD).DoublePlaysTurned()` that hide
  awkward parsing (e.g., per-game team double plays only appear in the
  free-text `info.FIELDING.DP` block — the SDK parses it).
- Public types live in `pkg/mlb/<resource>_types.go`; methods on those types
  live in `pkg/mlb/<resource>.go`.
- The generated `gen.Client` is wrapped by `pkg/mlb.Client` — never returned
  to callers.

### Branching

Feature branches off `main` with `type/short-description` (e.g.,
`feat/add-live-feed-wrapper`, `fix/dp-parser-edge-case`). PRs only — no
direct commits to `main`.

### Commit messages

Conventional Commits (`feat`, `fix`, `docs`, `refactor`, `chore`, `test`,
`perf`). Scopes: `api` (OpenAPI spec), `gen` (generated client), `mlb`
(public SDK), `cli`, `docs`.

When committing via Claude Code, end with:

- `🤖 Generated with [Claude Code](https://claude.ai/code)`
- `Co-Authored-By: Claude <noreply@anthropic.com>`
