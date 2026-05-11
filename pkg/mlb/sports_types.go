// Copyright (c) 2026 John Dewey
//
// SPDX-License-Identifier: MIT

package mlb

// SportsQuery filters a sports listing. Every field is optional; an empty
// query returns every sport the MLB API tracks (MLB + affiliated minor
// leagues + college + independent baseball).
type SportsQuery struct {
	// SportID restricts the response to a single sport (1 = MLB, 11 = AAA,
	// 12 = AA, …).
	SportID int

	// Fields is a comma-separated field projection passed to the MLB API.
	Fields string
}

// Sports is the typed view of /api/v1/sports.
type Sports struct {
	Sports []Sport
}

// Sport returns the entry with the given sport id, or nil when not present
// in the response.
func (s *Sports) Sport(id int) *Sport {
	if s == nil {
		return nil
	}
	for i := range s.Sports {
		if s.Sports[i].ID == id {
			return &s.Sports[i]
		}
	}
	return nil
}

// Sport is one sport tracked by the MLB Stats API.
type Sport struct {
	ID           int
	Code         string // e.g. "mlb", "aaa"
	Link         string
	Name         string
	Abbreviation string
	SortOrder    int
	ActiveStatus bool
}
