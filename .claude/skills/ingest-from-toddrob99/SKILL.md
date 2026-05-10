---
name: ingest-from-toddrob99
description: Port MLB Stats API endpoints from the toddrob99/MLB-StatsAPI Python upstream into mlb-sdk's OpenAPI spec and pkg/mlb wrappers. Treats toddrob99 as the canonical source of truth so we never maintain our own endpoint catalog. Use when expanding pkg/mlb coverage.
---

# Ingest from toddrob99/MLB-StatsAPI

[`toddrob99/MLB-StatsAPI`](https://github.com/toddrob99/MLB-StatsAPI) has done
years of human reverse-engineering of `statsapi.mlb.com`. This skill ports
their work into our typed Go SDK without us re-doing the discovery.

## Source of truth

- Repo: `https://github.com/toddrob99/MLB-StatsAPI`
- Endpoint catalog: [`statsapi/endpoints.py`](https://github.com/toddrob99/MLB-StatsAPI/blob/master/statsapi/endpoints.py) — a single Python dict keyed by endpoint name, mapping to `{url, path_params, query_params, required_params, hydrate_options}`.
- Always fetch fresh — never assume a cached copy is current.

## Modes

Two ways to invoke this skill:

| Mode    | Trigger                                          | What it does                                                                                |
| ------- | ------------------------------------------------ | ------------------------------------------------------------------------------------------- |
| Single  | "ingest the `<name>` endpoint" / `--endpoint X`  | Port one endpoint, run `just ready`, commit, push.                                          |
| Batch   | "ingest everything new" / `--all`                | Iterate every endpoint not yet in `manifest.json (sibling of this file)`. **One commit per endpoint.** |

Default to single mode unless the user explicitly says "all" / "everything" / "batch".

## Per-endpoint procedure

**Before doing anything**, the universal rules live elsewhere:

- [`AGENTS.md` → Hard rules](../../../AGENTS.md#hard-rules) — public-types-only,
  one-table-per-public-function, 100% coverage, named components, `just ready`
  gate.
- [`docs/development.md` → Adding a new endpoint](../../../docs/development.md#adding-a-new-endpoint) —
  the eight-file recipe (`api/openapi.yaml`, `go generate`, `pkg/mlb/<name>_types.go`,
  `pkg/mlb/<name>.go`, `pkg/mlb/<name>_test.go`, `examples/<name>.go`, `README.md`,
  `just ready`).
- [`docs/development.md` → Testing conventions](../../../docs/development.md#testing-conventions) —
  the required failure rows for `Client` methods (404, 5xx, malformed JSON,
  network failure, empty body) and the server-per-case pattern.

This skill is the **toddrob99 wrapper** around that recipe — it tells you how
to translate their Python catalog into the inputs the recipe expects, and
covers the gotchas only relevant to ingesting from upstream.

### 1. Fetch the upstream definition

```bash
TMP=$(mktemp -d)
git clone --depth 1 https://github.com/toddrob99/MLB-StatsAPI "$TMP/upstream"
ENDPOINTS_PY="$TMP/upstream/statsapi/endpoints.py"
```

Inspect the dict entry for the endpoint you're porting. Each entry looks roughly like:

```python
"team_stats": {
    "url": "https://statsapi.mlb.com/api/{ver}/teams/{teamId}/stats",
    "path_params": {
        "ver": {"type": "str", "default": "v1", "leading_slash": False, "trailing_slash": False, "required": True},
        "teamId": {"type": "str", "default": "", "leading_slash": False, "trailing_slash": False, "required": True},
    },
    "query_params": ["season", "stats", "group", "gameType"],
    "required_params": [["stats", "group"]],
    "hydrate_options": [...]
}
```

### 2. Translate to OpenAPI

| toddrob99 field     | OpenAPI equivalent                                                            |
| ------------------- | ----------------------------------------------------------------------------- |
| `url`               | Replace `{ver}` with the resolved version (`v1` / `v1.1`); becomes `paths.<route>`. |
| `path_params`       | `parameters: [{ in: path, name, required: true, schema: { type } }]`.        |
| `query_params`      | `parameters: [{ in: query, name, schema: { type } }]`.                       |
| `required_params`   | OpenAPI cannot express *one-of* required combos. Document in the operation `description` and validate in the Go wrapper with a returned `ErrInvalidQuery`. |
| `hydrate_options`   | A single `hydrate` query param of type string; document allowed values in description, do **not** enum-restrict (the list grows). |
| Response schema     | toddrob99 doesn't model responses — fetch a sample with `curl <url>` and translate per the named-component rule in `docs/development.md`. |

### 3. Run the eight-file recipe

Execute every step in
[`docs/development.md` → Adding a new endpoint](../../../docs/development.md#adding-a-new-endpoint).
The toddrob99 translation hints below are the only thing this skill
adds beyond that recipe — gotchas only relevant when the upstream
input is a Python dict from `endpoints.py`.

### 4. Update the manifest

Append the endpoint name to `manifest.json (sibling of this file)` once verified. The
manifest is how batch mode knows what's already done.

### 5. Verify

- `just ready` passes (fmt + vet + lint).
- `go test -coverprofile=/tmp/c.out ./pkg/mlb/...` reports 100.0%.
- `go run examples/<name>.go` runs cleanly against the live API (or the user's
  fixture).

## Translation gotchas

- **Path version variants** — toddrob99 uses `{ver}` for the API version
  segment. Most endpoints are v1; the live game feed is v1.1. The OpenAPI spec
  treats them as separate paths, not a parameterized `{ver}`.
- **Required combos** — `required_params: [["a", "b"]]` means *both* `a` and
  `b` must be set. Encode as a runtime check in the Go wrapper that returns
  `ErrInvalidQuery`. Do not try to encode in OpenAPI.
- **Hydrate** — toddrob99 documents the full hydrate vocabulary, but it grows
  every season. Expose `hydrate` as a free-form string and let users compose;
  do not turn it into typed constants.
- **Response shapes** — toddrob99 does not document response shapes. Fetch a
  sample with `curl https://statsapi.mlb.com/<path>?<params>` and translate.
  Use `additionalProperties: true` on every response schema so the spec
  remains forward-compatible.

## Output

Each successful invocation produces:

- One commit per endpoint, message `feat(api,mlb): Wrap <endpoint> from toddrob99 catalog`.
- An updated `manifest.json (sibling of this file)`.
- A push to `origin/main` (this repo is skunkworks; commits land directly on main).

If `just ready` or coverage fails, **stop the batch immediately** and surface the failure. Do not commit broken state.
