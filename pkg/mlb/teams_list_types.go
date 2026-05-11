// Copyright (c) 2026 John Dewey
//
// SPDX-License-Identifier: MIT

package mlb

// TeamsQuery filters a teams listing. Every field is optional; an empty
// query returns every team the MLB API tracks across every sport.
type TeamsQuery struct {
	// Season constrains the snapshot to a specific year (team rosters,
	// names, and league affiliations vary year-over-year).
	Season int

	// ActiveStatus filters by activity: "ACTIVE", "INACTIVE", or "BOTH".
	// Zero value yields the API default.
	ActiveStatus string

	// LeagueIDs is a comma-separated list of league ids (e.g. "103,104"
	// for AL+NL). Zero value yields every league for the matching sport.
	LeagueIDs string

	// SportID restricts to a sport (1 = MLB, 11 = AAA, …).
	SportID int

	// SportIDs is a comma-separated list of sport ids. Use instead of
	// SportID when you want multiple sports in one call.
	SportIDs string

	// GameType filters by game-type code: "R" (regular), "S" (spring),
	// "E" (exhibition), "A" (all-star), "D" / "F" / "L" / "W" (post).
	GameType string

	// Hydrate is a comma-separated hydrate string — see TeamQuery.Hydrate
	// for common values.
	Hydrate string

	// Fields restricts the response to a comma-separated field projection.
	Fields string
}

// Teams is the typed view of /api/v1/teams — a list of TeamInfo records.
type Teams struct {
	Teams []TeamInfo
}

// Team returns the entry with the given team id, or nil when not present
// in the response.
func (t *Teams) Team(id int) *TeamInfo {
	if t == nil {
		return nil
	}
	for i := range t.Teams {
		if t.Teams[i].ID == id {
			return &t.Teams[i]
		}
	}
	return nil
}
