// Copyright (c) 2026 John Dewey
//
// SPDX-License-Identifier: MIT

package mlb

// DivisionsQuery filters a divisions listing. Every field is optional; an
// empty query returns every division across MLB and affiliated leagues.
type DivisionsQuery struct {
	// DivisionID restricts the response to a single division by id (e.g.
	// 200 = AL West, 204 = NL East).
	DivisionID int

	// LeagueID restricts the response to a single league (e.g. 103 = AL,
	// 104 = NL).
	LeagueID int

	// SportID restricts the response to a sport (1 = MLB, 11 = AAA, etc.).
	SportID int

	// Season constrains the divisions to a specific season's snapshot.
	Season int
}

// Divisions is the typed view of /api/v1/divisions.
type Divisions struct {
	Divisions []Division
}

// Division returns the entry with the given division id, or nil when not
// present in the response.
func (d *Divisions) Division(id int) *Division {
	if d == nil {
		return nil
	}
	for i := range d.Divisions {
		if d.Divisions[i].ID == id {
			return &d.Divisions[i]
		}
	}
	return nil
}

// Division is one entry in the divisions list. NumPlayoffTeams is only
// populated for divisions where the API reports playoff structure (it is
// omitted for several historical / informal divisions).
type Division struct {
	ID              int
	Name            string
	Season          string
	NameShort       string
	Link            string
	Abbreviation    string
	League          Ref
	Sport           Ref
	HasWildcard     bool
	SortOrder       int
	NumPlayoffTeams int
	Active          bool
}
