// Copyright (c) 2026 John Dewey
//
// SPDX-License-Identifier: MIT

package mlb

import "strconv"

// LeagueID is the MLB Stats API numeric league identifier. Use the constants
// below for the major leagues; pass via `?leagueId=...` query parameters.
type LeagueID int

// MLB league identifiers as the API uses them.
const (
	AL LeagueID = 103 // American League
	NL LeagueID = 104 // National League
)

// Int returns the numeric form expected by the MLB Stats API.
func (l LeagueID) Int() int { return int(l) }

// String returns the MLB-API-formatted league id string. Used to assemble
// `leagueId` query parameters; multiple leagues are joined with commas at
// the call site.
func (l LeagueID) String() string { return strconv.Itoa(int(l)) }
