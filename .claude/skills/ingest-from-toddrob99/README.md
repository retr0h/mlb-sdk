# 🥷 ingest-from-toddrob99

A Claude skill that ports MLB Stats API endpoints from
[toddrob99/MLB-StatsAPI](https://github.com/toddrob99/MLB-StatsAPI) into this
repo's typed Go SDK without us maintaining a parallel endpoint catalog.

## 🎯 What it does

Treats toddrob99's `endpoints.py` dict as upstream truth. Translates one entry
(or all unported entries) into:

- a `paths:` entry plus named components in `api/openapi.yaml`,
- the regenerated `internal/gen/client.gen.go`,
- handwritten public types and a `Client` method in `pkg/mlb/`,
- a table-driven test holding 100.0% statement coverage,
- a runnable `examples/<name>.go`,
- a row in the README endpoints table.

## 🚀 Usage

The skill activates when you ask Claude (in this repo) to port endpoints from
toddrob99. Use any of these phrasings:

| You say                                                  | What runs                                                          |
| -------------------------------------------------------- | ------------------------------------------------------------------ |
| `ingest the team_leaders endpoint from toddrob99`        | Single mode — port `team_leaders`, one commit, push.               |
| `port the standings endpoint`                            | Single mode — same flow, alternate phrasing.                       |
| `port everything new from toddrob99`                     | Batch mode — every endpoint not in `tools/ingest/manifest.json`.   |
| `ingest the next 5 unported endpoints`                   | Bounded batch — same as above, capped to 5.                        |
| `regenerate the manifest from current SDK state`         | Manifest sync — re-derive `manifest.json` from `pkg/mlb/`.         |

The keywords Claude looks for are **`ingest` / `port` / `from toddrob99` /
`unported`**. Mention any of those alongside an endpoint name (or `everything`
/ `all` / `unported`) and the skill takes over.

You don't need to know toddrob99's exact endpoint key — just paste the URL or
say "the team-leaders one" and Claude will resolve it from
`statsapi/endpoints.py`.

Skill output: one commit per endpoint to `main`, pushed immediately.

## 🛑 How to stop / scope

- `stop after the current endpoint` — finish the in-flight commit, then halt.
- `dry run — show me the OpenAPI diff first` — produce the spec change but
  don't write generated code or commit.
- `ingest <name> but skip the example` — same flow, omit `examples/<name>.go`
  (you'll add it later by hand).

## 📋 Two modes

| Mode    | Default? | Behavior                                                                |
| ------- | -------- | ----------------------------------------------------------------------- |
| Single  | ✅       | One endpoint → one commit → push.                                       |
| Batch   |          | Loop over every endpoint not already in `tools/ingest/manifest.json`. **Still one commit per endpoint** so each diff is reviewable. |

## 🔁 What happens per endpoint

1. Fetch fresh from toddrob99 (`git clone --depth 1`).
2. Read the dict entry for the endpoint.
3. Run the 8-step recipe from
   [`docs/development.md`](../../../docs/development.md#adding-a-new-endpoint).
4. Verify with `just ready` + 100% coverage.
5. Append to `tools/ingest/manifest.json`.
6. Commit + push.

## ✅ Done criteria

- `just ready` passes.
- Coverage stays at 100.0%.
- Example runs cleanly against the live API.
- Manifest updated.

If any check fails the skill **stops batch mode immediately** so we never
commit broken state.

## 🧠 Translation cheat sheet

See [`SKILL.md`](SKILL.md) for the full Python→OpenAPI translation table,
including how to handle `path_params`, `query_params`, `required_params`,
`hydrate_options`, and response shapes (which toddrob99 does not document —
sample the live API and translate).

## 📦 Manifest

`tools/ingest/manifest.json` tracks which endpoints we've already ported, keyed
by toddrob99's name (e.g. `team_stats`). The skill consults this in batch mode
to skip work already done.
