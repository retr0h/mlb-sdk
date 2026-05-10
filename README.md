[![go report card](https://goreportcard.com/badge/github.com/retr0h/mlb-sdk?style=for-the-badge)](https://goreportcard.com/report/github.com/retr0h/mlb-sdk)
[![license](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=for-the-badge)](LICENSE)
[![build](https://img.shields.io/github/actions/workflow/status/retr0h/mlb-sdk/go.yml?style=for-the-badge)](https://github.com/retr0h/mlb-sdk/actions/workflows/go.yml)
[![just](https://img.shields.io/badge/just-command%20runner-blue?style=for-the-badge)](https://github.com/casey/just)
[![conventional commits](https://img.shields.io/badge/Conventional%20Commits-1.0.0-yellow.svg?style=for-the-badge)](https://conventionalcommits.org)
[![go reference](https://img.shields.io/badge/go-reference-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://pkg.go.dev/github.com/retr0h/mlb-sdk/pkg/mlb)

<h1 align="center">
<pre>
█▀▄▀█ █░░ █▄▄
█░▀░█ █▄▄ █▄█
</pre>
</h1>

<p align="center">⚾ Idiomatic Go library for the public MLB Stats API.</p>

A typed Go client for `statsapi.mlb.com` — the same JSON feed that powers
MLB.com Gameday, Baseball-Reference, and most of the community Python and R
wrappers. MLB doesn't publish an OpenAPI specification, so this repo authors
one and generates a client from it via
[`oapi-codegen`](https://github.com/oapi-codegen/oapi-codegen). The public
surface in `pkg/mlb` papers over the API's quirks (e.g., per-game team
double-plays only appear in a free-text `info.FIELDING.DP` block — the SDK
exposes them as a typed method).

## ✨ Features

- 🧩 **OpenAPI-first** — hand-authored spec at `api/openapi.yaml`, generated
  client under `internal/gen`. The generated client is hidden behind an
  `internal/` boundary; consumers only see idiomatic Go.
- 🧠 **Idiomatic surface** — typed `mlb.TeamID` constants (`mlb.LAD`),
  `time.Time` for dates, helper methods that hide ugly upstream encodings.
- 🪄 **Hides API quirks** — `box.Team(mlb.LAD).DoublePlaysTurned()` parses
  the free-text `info` block so callers don't have to.
- 🧪 **Test-friendly** — `WithBaseURL` lets you point the client at an
  `httptest.Server` for hermetic tests.

## 📦 Install

```bash
go get github.com/retr0h/mlb-sdk/pkg/mlb
```

## 🚀 Quick start

```go
package main

import (
    "context"
    "fmt"
    "time"

    "github.com/retr0h/mlb-sdk/pkg/mlb"
)

func main() {
    c, _ := mlb.New()

    games, _ := c.Schedule(context.Background(), mlb.ScheduleQuery{
        Team: mlb.LAD,
        On:   time.Now().AddDate(0, 0, -1), // yesterday
    })
    for _, g := range games {
        box, _ := c.Boxscore(context.Background(), g.GamePk)
        fmt.Printf("%s vs %s — Dodgers turned %d DPs\n",
            g.Away.Name, g.Home.Name,
            box.Team(mlb.LAD).DoublePlaysTurned())
    }
}
```

## ⚙️ How it works

```
api/openapi.yaml          Hand-authored OpenAPI 3.0 spec
        │
        │  oapi-codegen (via `just generate`)
        ▼
internal/gen/             Generated typed HTTP client. Not importable externally.
        │
        │  wrapped, parsed, normalized
        ▼
pkg/mlb/                  Public, idiomatic Go SDK (the only thing you import)
```

The generated layer handles HTTP, query-string assembly, and JSON parsing.
The handwritten layer hides the rough edges the MLB API exposes (free-text
fielding annotations, optional fields everywhere, mixed v1 / v1.1 paths).

## 💡 Inspiration

This module exists because the MLB Stats API is undocumented and the most
useful references are Python and R libraries. We owe a debt to:

- [appac/mlb-data-api-docs](https://appac.github.io/mlb-data-api-docs/) — the
  most-cited community reference for the MLB Stats API endpoint shapes.
- [toddrob99/MLB-StatsAPI](https://github.com/toddrob99/MLB-StatsAPI) — the
  Python wrapper whose source code is the closest thing to a reference manual
  for per-endpoint behavior.
- [BillPetti/baseballr](https://github.com/BillPetti/baseballr) — R package
  with extensive coverage of MLB / college / minor league stats.
- The MLB.com Gameday viewer, which is what `statsapi.mlb.com` actually
  powers — when in doubt, what Gameday displays is the authoritative answer.

## 🗺️ Roadmap

See [docs/roadmap.md](docs/roadmap.md) for the phased plan, including the
AI-driven reverse-engineering tasks (LLM-as-spec-author, differential
probing, and fixture-driven spec discovery) and the planned `mlb mcp`
subcommand.

## 📖 Documentation

See the [package documentation][] on pkg.go.dev for API details.

## 🤝 Contributing

See the [Development](docs/development.md) guide for prerequisites, setup,
and conventions. See the [Contributing](docs/contributing.md) guide before
submitting a PR.

## 📄 License

The [MIT][] License.

[package documentation]: https://pkg.go.dev/github.com/retr0h/mlb-sdk/pkg/mlb
[MIT]: LICENSE
