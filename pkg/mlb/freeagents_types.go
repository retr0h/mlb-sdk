// Copyright (c) 2026 John Dewey
//
// SPDX-License-Identifier: MIT
//
package mlb
//
// FreeAgentsQuery filters a free-agents listing. Every field is optional.
type FreeAgentsQuery struct {
	Season  int
	Order   string
	Hydrate string
	Fields  string
}
//
// FreeAgents is the typed view of /api/v1/people/freeAgents.
type FreeAgents struct {
	FreeAgents []FreeAgent
}
//
// FreeAgent is one free-agent row — the player, original/new team, and
// signing metadata.
type FreeAgent struct {
	Player       Person
	OriginalTeam TeamRef
	NewTeam      TeamRef
	Notes        string
	DateSigned   string // YYYY-MM-DD (empty if unsigned)
	DateDeclared string // YYYY-MM-DD
	Position     PrimaryPosition
}
