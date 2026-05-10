# mlb-sdk

Idiomatic Go SDK and CLI for the MLB Stats API.

> ⚠️ Pre-release. v0.1 covers the subset of endpoints needed by
> [freebies](https://github.com/bitsugar-io/freebies) (schedule, boxscore,
> play-by-play, live feed, team stats). Surface area will grow.

## Why

MLB does not publish an OpenAPI specification for `statsapi.mlb.com`. Existing
tooling is mostly Python or R. This module provides:

- A hand-authored OpenAPI spec for the endpoints we use.
- A typed HTTP client generated from that spec via
  [`oapi-codegen`](https://github.com/oapi-codegen/oapi-codegen).
- An idiomatic Go SDK on top that papers over the API's quirks (e.g., per-game
  team double-plays only appear in a free-text `info.FIELDING.DP` block — the
  SDK exposes them as a typed method).
- A small CLI for ad-hoc lookups.

## Install

```sh
go get github.com/retr0h/mlb-sdk/pkg/mlb
```

The CLI:

```sh
go install github.com/retr0h/mlb-sdk@latest
```

## Quick start

```go
import "github.com/retr0h/mlb-sdk/pkg/mlb"

c := mlb.New()
games, _ := c.Schedule(ctx, mlb.ScheduleQuery{Team: mlb.LAD, On: time.Now()})
box, _   := c.Boxscore(ctx, games[0].GamePk)
fmt.Println(box.Team(mlb.LAD).DoublePlaysTurned())
```

## Layout

```
api/openapi.yaml      hand-authored OpenAPI 3.0 spec
internal/gen/         generated client (not importable externally)
pkg/mlb/              public, idiomatic SDK surface
cmd/                  cobra CLI tree
```

## Development

```sh
just deps       # mise + go modules
just generate   # regenerate internal/gen/client.gen.go
just test
just ready      # fmt + vet + lint
```

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

## License

MIT
