[![go report card](https://goreportcard.com/badge/github.com/retr0h/mlb-sdk?style=for-the-badge)](https://goreportcard.com/report/github.com/retr0h/mlb-sdk)
[![license](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=for-the-badge)](LICENSE)
[![build](https://img.shields.io/github/actions/workflow/status/retr0h/mlb-sdk/go.yml?style=for-the-badge)](https://github.com/retr0h/mlb-sdk/actions/workflows/go.yml)
[![just](https://img.shields.io/badge/just-command%20runner-blue?style=for-the-badge)](https://github.com/casey/just)
[![conventional commits](https://img.shields.io/badge/Conventional%20Commits-1.0.0-yellow.svg?style=for-the-badge)](https://conventionalcommits.org)
[![go reference](https://img.shields.io/badge/go-reference-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://pkg.go.dev/github.com/retr0h/mlb-sdk/pkg/mlb)

# mlb-sdk

⚾ Idiomatic Go library for the public MLB Stats API.

A typed Go client for `statsapi.mlb.com`. MLB does not publish an OpenAPI
specification, so this repo authors one and generates the underlying
client from it. Builds on years of community reverse engineering by
[toddrob99/MLB-StatsAPI][] (Python) and [BillPetti/baseballr][] (R); the
public surface in `pkg/mlb` hides the API's quirks behind idiomatic
helpers.

## 📦 Install

```bash
go get github.com/retr0h/mlb-sdk/pkg/mlb
```

## ⚙️ Endpoints

| Endpoint                         | SDK method                     |
| -------------------------------- | ------------------------------ |
| `/api/v1/schedule`               | `Client.Schedule`              |
| `/api/v1/game/{gamePk}/boxscore` | `Client.Boxscore`              |

Additional endpoints in flight — see [docs/roadmap.md][].

## ✨ Features

| Feature             | Description                                                |
| ------------------- | ---------------------------------------------------------- |
| OpenAPI-first       | Hand-authored spec + generated client (oapi-codegen)       |
| Idiomatic surface   | `time.Time`, typed `mlb.TeamID`, helper methods            |
| Hides API quirks    | e.g. `box.Team(mlb.LAD).DoublePlaysTurned()`               |
| Test-friendly       | `WithBaseURL` injects an `httptest.Server` for fixtures    |

## 📖 Documentation

See the [package documentation][] on pkg.go.dev for API details.

## 🤝 Contributing

See the [Development][] guide for prerequisites, setup, and conventions.
See the [Contributing][] guide before submitting a PR.

## 📄 License

The [MIT][] License.

[toddrob99/MLB-StatsAPI]: https://github.com/toddrob99/MLB-StatsAPI
[BillPetti/baseballr]: https://github.com/BillPetti/baseballr
[docs/roadmap.md]: docs/roadmap.md
[package documentation]: https://pkg.go.dev/github.com/retr0h/mlb-sdk/pkg/mlb
[Development]: docs/development.md
[Contributing]: docs/contributing.md
[MIT]: LICENSE
