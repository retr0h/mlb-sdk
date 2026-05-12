// Copyright (c) 2026 John Dewey
//
// SPDX-License-Identifier: MIT

package mlb

import (
	"context"
	"fmt"

	"github.com/retr0h/mlb-sdk/internal/gen"
)

const rosterDateFmt = "2006-01-02"

// Roster fetches the roster for a team. Filter by rosterType ("active",
// "40Man", "fullSeason", …), season, and date.
//
// Example:
//
//	r, _ := c.Roster(ctx, 119, mlb.RosterQuery{Season: 2024, RosterType: "active"})
//	for _, e := range r.Roster {
//	    fmt.Println(e.Person.FullName, e.JerseyNumber, e.Position.Abbreviation)
//	}
func (c *Client) Roster(ctx context.Context, teamID int, q RosterQuery) (*Roster, error) {
	params := &gen.GetTeamRosterParams{}
	if q.RosterType != "" {
		params.RosterType = ptr(q.RosterType)
	}
	if q.Season != 0 {
		params.Season = ptr(q.Season)
	}
	if !q.On.IsZero() {
		params.Date = ptr(q.On.Format(rosterDateFmt))
	}
	if q.Hydrate != "" {
		params.Hydrate = ptr(q.Hydrate)
	}
	if q.Fields != "" {
		params.Fields = ptr(q.Fields)
	}

	resp, err := c.raw.GetTeamRosterWithResponse(ctx, teamID, params)
	if err != nil {
		return nil, fmt.Errorf("mlb: roster: %w", err)
	}
	if resp.StatusCode() == 404 {
		return nil, ErrNotFound
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("mlb: roster: unexpected status %d", resp.StatusCode())
	}
	return rosterFromGen(resp.JSON200), nil
}

func rosterFromGen(r *gen.RosterResponse) *Roster {
	out := &Roster{}
	if r == nil {
		return out
	}
	if r.Link != nil {
		out.Link = *r.Link
	}
	if r.Roster != nil {
		out.Roster = make([]RosterEntry, 0, len(*r.Roster))
		for _, e := range *r.Roster {
			out.Roster = append(out.Roster, rosterEntryFromGen(e))
		}
	}
	return out
}

func rosterEntryFromGen(e gen.RosterEntry) RosterEntry {
	out := RosterEntry{}
	if e.Person != nil {
		out.Person = personFromGen(*e.Person)
	}
	if e.JerseyNumber != nil {
		out.JerseyNumber = *e.JerseyNumber
	}
	if e.Position != nil {
		out.Position = primaryPositionFromGen(*e.Position)
	}
	if e.Status != nil {
		out.Status = rosterStatusFromGen(*e.Status)
	}
	return out
}

func rosterStatusFromGen(s gen.RosterStatus) RosterStatus {
	out := RosterStatus{}
	if s.Code != nil {
		out.Code = *s.Code
	}
	if s.Description != nil {
		out.Description = *s.Description
	}
	return out
}
