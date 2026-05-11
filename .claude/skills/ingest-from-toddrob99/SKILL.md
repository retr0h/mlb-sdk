---
name: ingest-from-toddrob99
description: Port MLB Stats API endpoints from the toddrob99/MLB-StatsAPI Python upstream into mlb-sdk's OpenAPI spec and pkg/mlb wrappers. Treats toddrob99 as the canonical source of truth so we never maintain our own endpoint catalog. Use when expanding pkg/mlb coverage.
---

# Ingest from toddrob99/MLB-StatsAPI

[`toddrob99/MLB-StatsAPI`](https://github.com/toddrob99/MLB-StatsAPI) has done
years of human reverse-engineering of `statsapi.mlb.com`. This skill ports
their work into our typed Go SDK without us re-doing the discovery.

> **Architectural reminders before you read anything else:**
>
> 1. `internal/gen/` is the generated implementation detail; `pkg/mlb/`
>    is the public surface that wraps it. **No `gen.X` type ever appears
>    in an exported signature.** Every public method on `*mlb.Client`
>    has a private `<name>FromGen` converter that translates the
>    pointer-heavy generated structs into clean idiomatic Go.
>
> 2. **Every field the upstream API returns is exposed as a public Go
>    field on the wrapping type.** Helper methods are additive — they
>    encode awkward parsing (`DoublePlaysTurned()` reads a free-text
>    info block) or domain shortcuts. They never replace direct field
>    access. If the API returns `runs`, the wrapping type has a
>    `Runs int` field, not just a `Runs()` method that hides it.
>
> The full rules are in AGENTS.md and docs/development.md (read them in
> step 0 below); these two paragraphs are the most important things to
> internalize before touching code.

## Source of truth

- Repo: `https://github.com/toddrob99/MLB-StatsAPI`
- Endpoint catalog: [`statsapi/endpoints.py`](https://github.com/toddrob99/MLB-StatsAPI/blob/master/statsapi/endpoints.py) — a single Python dict keyed by endpoint name, mapping to `{url, path_params, query_params, required_params, hydrate_options}`.
- Always fetch fresh — never assume a cached copy is current.

## Modes

Three ways to invoke this skill:

| Mode    | Trigger                                                       | What it does                                                                                              |
| ------- | ------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------- |
| Single  | "ingest the `<name>` endpoint" / `--endpoint X`               | Port one endpoint, run `just ready`, commit, push.                                                        |
| Batch   | "ingest everything new" / "port all from toddrob99" / `--all` | Iterate every endpoint not in `manifest.json` (sibling of this file). **One commit per endpoint.**         |
| Rescan  | "check toddrob99 for updates" / "rescan upstream" / `--rescan` | Re-fingerprint every ingested endpoint, diff against the manifest, surface drift. See Rescan procedure below. |

Default to single mode unless the user explicitly says "all" / "everything" / "batch" / "rescan".

## Rescan procedure

The MLB Stats API and toddrob99's wrappers evolve. Without rescan, our
ported endpoints silently rot when toddrob99 adds query params or fixes
path bugs. The manifest's `upstreamSha` per entry plus the fingerprint
helper detect this.

1. Clone toddrob99 fresh (`git clone --depth 1 ...`); record its commit SHA.
2. Run `python3 .claude/skills/ingest-from-toddrob99/fingerprint.py <clone>`
   — emits `{endpoint_name: sha256_hex}` JSON to stdout.
3. For every entry in `manifest.json#endpoints`:
   - If the upstream key is **missing** from the fingerprint output → 🔴
     endpoint disappeared upstream. Surface and let the user decide
     (deprecate, rename, leave).
   - If the upstream `sha` matches the manifest's `upstreamSha` → 🟢 in sync.
   - If the upstream `sha` differs → 🟡 drift. Inspect the diff between
     the upstream entry and what we last ingested, classify:
     - **Additive** (new optional query param, new hydrate option): apply
       the diff to `api/openapi.yaml`, regenerate, update the wrapper if
       needed, run `just ready`, bump the `upstreamSha` in manifest.json.
     - **Breaking** (path changed, required-combo changed, param removed):
       stop, surface to the user with the diff inline, do not auto-apply.
4. Also report endpoints in the upstream fingerprint that are **not** in
   our manifest — those are unported endpoints (candidates for batch
   ingest).
5. Bump the manifest's top-level `upstreamRef` to the new clone's SHA so
   the next rescan starts from the same baseline.

Rescan is **read-only by default**. Mutating actions (apply additive
diff, re-ingest a drifted endpoint) need explicit user confirmation
unless the user originally said "rescan and auto-apply additive drift".

## Step 0: load the universal rules into your context (REQUIRED)

This skill does **not** duplicate the universal rules — it would drift. They
live in two files you MUST `Read` before doing any work for this skill:

```
Read tool → AGENTS.md
Read tool → docs/development.md
```

You are not done with step 0 until both files are in your context. The skill's
instructions below assume you have internalized:

- the **Hard rules** in `AGENTS.md` (public-types-only at the boundary,
  one table-driven test per public function, 100.0% coverage gate, named
  components only in OpenAPI, `just ready` before commit),
- the **Adding a new endpoint** eight-file recipe in `docs/development.md`,
- the **Testing conventions** section in `docs/development.md` listing the
  required failure rows for `Client` methods (404, 5xx, malformed JSON,
  network failure, empty body) and the server-per-case `httptest.NewServer`
  pattern.

If you skip this step you will produce broken or non-conforming code. The
skill is the **toddrob99 wrapper** around the recipe in those two files — it
only tells you how to translate Python upstream into the recipe's inputs and
covers gotchas only relevant when ingesting from upstream.

## Per-endpoint procedure

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
