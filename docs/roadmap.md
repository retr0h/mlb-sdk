# Roadmap

Tracking near-term scope. Phases 2 and 3 are concrete; the AI tasks below are
research items — useful enough to scope, fuzzy enough that the acceptance
criteria will sharpen as we go.

## Phase 2 — Remaining v0.1 endpoints

Wrap the four endpoints already in `pkg/api/openapi.yaml` that don't yet have
public surfaces, plus matching CLI subcommands.

- [ ] `pkg/mlb.PlayByPlay(ctx, gamePk)` — return `[]Play` with typed
      `EventType`, `result.event` constants, and helpers like
      `Play.IsDoublePlay()`.
- [ ] `pkg/mlb.LiveFeed(ctx, gamePk)` — same shape as PlayByPlay, sourced from
      the v1.1 `feed/live` endpoint. Decide whether to dedupe with PlayByPlay or
      expose distinctly.
- [ ] `pkg/mlb.TeamStats(ctx, q)` — query season / by-date-range stats with
      typed `StatGroup` constants (Hitting, Pitching, Fielding).
- [ ] Migrate `freebies` to consume `pkg/mlb` once v0.1 is tagged.

## Phase 3 — Polish layer from kvlt

Mechanical port of files we already know we want. Mostly copy-and-adapt.

- [ ] `.github/workflows/` — go.yml, release.yml, dep-review, commit-lint,
      labeler, stale, greetings, report-card.
- [ ] `install.sh` adapted from kvlt's installer.
- [ ] `AI_POLICY.md`, `CODE_OF_CONDUCT.md`.
- [ ] `docs/development.md`, `docs/contributing.md`, `docs/recipes.md`.
- [ ] `asset/` directory for any banner / logo art.
- [ ] Themed CLI banner via `internal/cli` (lipgloss) — optional.

## AI reverse-engineering tasks

Three escalating bets on using AI to expand the OpenAPI spec coverage past what
we've hand-authored. Most of the MLB Stats API surface is undocumented or
documented only by community Python wrappers
([toddrob99/MLB-StatsAPI](https://github.com/toddrob99/MLB-StatsAPI),
[BillPetti/baseballr](https://github.com/BillPetti/baseballr)). Always import
their endpoint catalog first as the human-curated baseline before running any of
the below.

### 1. LLM-as-spec-author

**Goal:** quick coverage expansion. Feed an LLM a small set of sample JSON
responses for a new endpoint and have it draft `components/schemas/...` entries
that match `CLAUDE.md`'s spec authoring rules (named components, no inline
nesting).

**Acceptance criteria:**

- A `tools/spec-author/` directory containing a script that takes
  `(endpoint, sampleJSON)` and emits a candidate YAML fragment for
  `pkg/api/openapi.yaml`.
- Output passes `oapi-codegen` without errors.
- Documented limitation: human review required — LLM hallucinates field types
  and misses optionality.

**Effort:** small. **Risk:** low. **Value:** modest — saves typing, but output
quality matches what one careful afternoon of hand-coding produces.

### 2. Differential probing

**Goal:** discover undocumented query parameters and `hydrate` flag combinations
by firing every documented and adjacent variant against a real endpoint,
capturing all responses, and clustering the resulting schemas.

**Acceptance criteria:**

- A `tools/probe/` script that, given an endpoint, expands a parameter matrix
  (dates, seasons, statGroup × statType, hydrate=…), fires it, and stores
  `(params, response_schema_fingerprint)`.
- An LLM pass that clusters fingerprints into schema variants and describes
  which params toggle which fields.
- Output written to `docs/probes/<endpoint>.md` for human review.
- Stretch: feed the discovered variants back into the OpenAPI spec as `oneOf` /
  discriminator branches.

**Effort:** medium. **Risk:** medium (API rate limits, request budget).
**Value:** high — finds parameters and fields humans miss.

### 3. Fixture-driven spec discovery

**Goal:** build a continually self-correcting OpenAPI spec by capturing a large
corpus of real responses, having an LLM analyze cross-corpus schema variation,
and emitting a spec that matches actual production shape — not just the first
response you happened to look at.

**Procedure (rough):**

- Build a corpus: every Dodgers game in 2025 × 5 endpoints = ~810 raw JSON
  captures. Persist under `tools/corpus/` (gitignored, regen-able).
- Per endpoint, compute per-field presence rates: always, sometimes, rare.
  Surface conditional fields (e.g., `Result.IsScoringPlay` only appears on RBI
  plays).
- LLM pass that proposes `oneOf` discriminators and `additionalProperties`
  decisions based on observed variation.
- Output:
  - A diff between the corpus-derived spec and `pkg/api/openapi.yaml`.
  - A `mlb spec verify` CLI command that runs the corpus-diff against a fresh
    capture and reports drift.
- Run quarterly to catch silent MLB schema changes before they break production.

**Acceptance criteria:**

- `tools/corpus-diff/` script reproducible end-to-end.
- A `docs/spec-drift-<date>.md` report from at least one full corpus pass, even
  if minor.
- The `mlb spec verify` CLI works.

**Effort:** large. **Risk:** medium-high (corpus storage, MLB rate limits, LLM
context budget for cross-cluster diffing). **Value:** very high — solves a class
of bugs humans are bad at catching, and the same machinery applies to every
other undocumented sports API we wrap later.

## AI-native consumer hooks

Nice-to-have once the SDK surface is stable. As a library-only repo, an MCP
server would ship as a separate companion module (`mlb-sdk-mcp`) that imports
`pkg/mlb` rather than as a subcommand.

- [ ] `mlb-sdk-mcp` — separate module exposing `pkg/mlb` as an MCP server over
      stdio so Claude Desktop / Cursor / agentic clients can call methods
      directly.
- [ ] `tools.json` schema bundle for tool-use APIs that don't speak MCP.
- [ ] LLM-friendly docstrings on every public method (concrete examples,
      embedded team-ID table, common-mistake warnings).
