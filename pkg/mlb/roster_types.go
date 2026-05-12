// Copyright (c) 2026 John Dewey
//
// SPDX-License-Identifier: MIT

package mlb

import "time"

// RosterQuery refines a team roster lookup. The teamId path parameter is
// taken as a method argument; everything in this struct is optional.
type RosterQuery struct {
	// RosterType selects the roster view: "active", "40Man", "fullSeason",
	// "allTime", "depthChart", etc. Zero value yields the API default.
	RosterType string

	// Season constrains the roster to a specific year.
	Season int

	// On views the roster as of a specific calendar date.
	On time.Time

	// Hydrate is a comma-separated hydrate string.
	Hydrate string

	// Fields restricts the response to a comma-separated field projection.
	Fields string
}

// Roster is the typed view of /api/v1/teams/{teamId}/roster.
type Roster struct {
	Link   string
	Roster []RosterEntry
}

// RosterEntry is one player slot on a team roster.
type RosterEntry struct {
	Person       Person
	JerseyNumber string
	Position     PrimaryPosition
	Status       RosterStatus
}

// RosterStatus is the roster-status code for a player (e.g. "A" for Active,
// "MIN" for Minor League Contract).
type RosterStatus struct {
	Code        string
	Description string
}
