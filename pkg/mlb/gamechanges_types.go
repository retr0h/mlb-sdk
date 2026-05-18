// Copyright (c) 2026 John Dewey
//
// SPDX-License-Identifier: MIT
//
package mlb
//
// GameChangesQuery filters a game-changes lookup. UpdatedSince is required.
type GameChangesQuery struct {
	UpdatedSince string // ISO-8601 timestamp (required)
	SportID      int
	GameType     string
	Season       int
	Fields       string
}
//
// GameChanges is the typed view of /api/v1/game/changes — a schedule-like
// response with games modified since UpdatedSince.
type GameChanges struct {
	TotalItems           int
	TotalEvents          int
	TotalGames           int
	TotalGamesInProgress int
	Dates                []GameChangesDate
}
//
// GameChangesDate is one date's worth of changed games.
type GameChangesDate struct {
	Date                 string
	TotalItems           int
	TotalEvents          int
	TotalGames           int
	TotalGamesInProgress int
	Games                []Game
}
