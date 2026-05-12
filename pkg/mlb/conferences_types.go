// Copyright (c) 2026 John Dewey
//
// SPDX-License-Identifier: MIT

package mlb

// ConferencesQuery filters a conferences listing. Every field is optional.
type ConferencesQuery struct {
	ConferenceID int
	Season       int
	Fields       string
}

// Conferences is the typed view of /api/v1/conferences.
type Conferences struct {
	Conferences []Conference
}

// Conference returns the entry with the given id, or nil when not present.
func (c *Conferences) Conference(id int) *Conference {
	if c == nil {
		return nil
	}
	for i := range c.Conferences {
		if c.Conferences[i].ID == id {
			return &c.Conferences[i]
		}
	}
	return nil
}

// Conference is one conference tracked by the MLB Stats API.
type Conference struct {
	ID           int
	Name         string
	Link         string
	Abbreviation string
	NameShort    string
	HasWildcard  bool
	League       Ref
	Sport        Ref
}
