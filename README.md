[![go report card](https://goreportcard.com/badge/github.com/retr0h/mlb-sdk?style=for-the-badge)](https://goreportcard.com/report/github.com/retr0h/mlb-sdk)
[![license](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=for-the-badge)](LICENSE)
[![build](https://img.shields.io/github/actions/workflow/status/retr0h/mlb-sdk/go.yml?style=for-the-badge)](https://github.com/retr0h/mlb-sdk/actions/workflows/go.yml)
[![conventional commits](https://img.shields.io/badge/Conventional%20Commits-1.0.0-yellow.svg?style=for-the-badge)](https://conventionalcommits.org)
[![built with just](https://img.shields.io/badge/Built_with-Just-black?style=for-the-badge&logo=just&logoColor=white)](https://just.systems)
![github commit activity](https://img.shields.io/github/commit-activity/m/retr0h/mlb-sdk?style=for-the-badge)
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=for-the-badge)](https://pkg.go.dev/github.com/retr0h/mlb-sdk/pkg/mlb)
[![hovnokod](https://raw.githubusercontent.com/tekk/hovnokod-badge/main/assets/badges/hovnokod-for-the-badge.svg)](https://github.com/tekk/hovnokod-badge)
[![MLB](https://img.shields.io/badge/MLB-002D72?style=for-the-badge&logo=mlb&logoColor=white)](https://mlb.com)

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

| Endpoint                              | Docs                          | Example                                  |
| ------------------------------------- | ----------------------------- | ---------------------------------------- |
| `/api/v1/schedule`                    | [Client.Schedule][d-sched]    | [schedule.go](examples/schedule.go)      |
| `/api/v1/game/{gamePk}/boxscore`      | [Client.Boxscore][d-box]      | [boxscore.go](examples/boxscore.go)      |
| `/api/v1/game/{gamePk}/playByPlay`    | [Client.PlayByPlay][d-pbp]    | [playbyplay.go](examples/playbyplay.go)  |
| `/api/v1.1/game/{gamePk}/feed/live`   | [Client.LiveFeed][d-live]     | [livefeed.go](examples/livefeed.go)      |
| `/api/v1/standings`                   | [Client.Standings][d-stand]   | [standings.go](examples/standings.go)    |
| `/api/v1/teams/{teamId}/stats`        | [Client.TeamStats][d-stats]   | [teamstats.go](examples/teamstats.go)    |
| `/api/v1/venues/{venueId}`            | [Client.Venue][d-venue]       | [venue.go](examples/venue.go)            |
| `/api/v1/divisions`                   | [Client.Divisions][d-divs]    | [divisions.go](examples/divisions.go)    |
| `/api/v1/league`                      | [Client.Leagues][d-leagues]   | [leagues.go](examples/leagues.go)        |
| `/api/v1/seasons`                     | [Client.Seasons][d-seas]      | [seasons.go](examples/seasons.go)        |
| `/api/v1/seasons/{seasonId}`          | [Client.Season][d-season]     | [season.go](examples/season.go)          |
| `/api/v1/sports`                      | [Client.Sports][d-sports]     | [sports.go](examples/sports.go)          |
| `/api/v1/teams`                       | [Client.Teams][d-teams]       | [teams.go](examples/teams.go)            |
| `/api/v1/teams/{teamId}`              | [Client.Team][d-team]         | [team.go](examples/team.go)              |

Run any example with `go run examples/<name>.go`. Roadmap for additional
endpoints lives in [docs/roadmap.md][].

## ✨ Features

| Feature             | Description                                                |
| ------------------- | ---------------------------------------------------------- |
| OpenAPI-first       | Hand-authored spec + generated client (oapi-codegen)       |
| Idiomatic surface   | `time.Time`, typed `mlb.TeamID`, helper methods            |
| Hides API quirks    | e.g. `box.Team(mlb.LAD).DoublePlaysTurned()`               |
| Test-friendly       | `WithBaseURL` injects an `httptest.Server` for fixtures    |

## 💡 Inspiration

This module exists because the MLB Stats API is undocumented and the most
useful prior art is in Python and R:

- [toddrob99/MLB-StatsAPI](https://github.com/toddrob99/MLB-StatsAPI)
- [appac/mlb-data-api-docs](https://appac.github.io/mlb-data-api-docs/)
- [BillPetti/baseballr](https://github.com/BillPetti/baseballr)

## 📖 Documentation

See the [package documentation][] on pkg.go.dev for API details.

## 🤝 Contributing

See the [Development][] guide for prerequisites, setup, and conventions.
See the [Contributing][] guide before submitting a PR.

## ⚖️ Copyright notice

This package and its author are not affiliated with MLB or any MLB team. This
module is a typed Go client for MLB's public Stats API. Use of MLB data is
subject to the notice posted at
<http://gdx.mlb.com/components/copyright.txt>.

## 📄 License

The [MIT][] License.

[toddrob99/MLB-StatsAPI]: https://github.com/toddrob99/MLB-StatsAPI
[BillPetti/baseballr]: https://github.com/BillPetti/baseballr
[docs/roadmap.md]: docs/roadmap.md
[package documentation]: https://pkg.go.dev/github.com/retr0h/mlb-sdk/pkg/mlb
[Development]: docs/development.md
[Contributing]: docs/contributing.md
[MIT]: LICENSE
[d-sched]: https://pkg.go.dev/github.com/retr0h/mlb-sdk/pkg/mlb#Client.Schedule
[d-box]:   https://pkg.go.dev/github.com/retr0h/mlb-sdk/pkg/mlb#Client.Boxscore
[d-pbp]:   https://pkg.go.dev/github.com/retr0h/mlb-sdk/pkg/mlb#Client.PlayByPlay
[d-live]:  https://pkg.go.dev/github.com/retr0h/mlb-sdk/pkg/mlb#Client.LiveFeed
[d-stand]: https://pkg.go.dev/github.com/retr0h/mlb-sdk/pkg/mlb#Client.Standings
[d-stats]: https://pkg.go.dev/github.com/retr0h/mlb-sdk/pkg/mlb#Client.TeamStats
[d-venue]: https://pkg.go.dev/github.com/retr0h/mlb-sdk/pkg/mlb#Client.Venue
[d-divs]:  https://pkg.go.dev/github.com/retr0h/mlb-sdk/pkg/mlb#Client.Divisions
[d-leagues]: https://pkg.go.dev/github.com/retr0h/mlb-sdk/pkg/mlb#Client.Leagues
[d-seas]:  https://pkg.go.dev/github.com/retr0h/mlb-sdk/pkg/mlb#Client.Seasons
[d-season]: https://pkg.go.dev/github.com/retr0h/mlb-sdk/pkg/mlb#Client.Season
[d-sports]: https://pkg.go.dev/github.com/retr0h/mlb-sdk/pkg/mlb#Client.Sports
[d-team]:  https://pkg.go.dev/github.com/retr0h/mlb-sdk/pkg/mlb#Client.Team
[d-teams]: https://pkg.go.dev/github.com/retr0h/mlb-sdk/pkg/mlb#Client.Teams
