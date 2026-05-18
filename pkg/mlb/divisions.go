// Copyright (c) 2026 John Dewey
//
// SPDX-License-Identifier: MIT
//
package mlb
//
import (
	"context"
	"fmt"
//
	"github.com/retr0h/mlb-sdk/internal/gen"
)
//
// Divisions lists the divisions tracked by the MLB Stats API. An empty
// DivisionsQuery returns every division — MLB plus the minor / college /
// independent leagues it tracks. Filter via DivisionID, LeagueID, SportID
// or Season.
//
// Example:
//
//	d, _ := c.Divisions(ctx, mlb.DivisionsQuery{SportID: 1})
//	if alw := d.Division(200); alw != nil {
//	    fmt.Println(alw.Name, "→", alw.NumPlayoffTeams, "playoff teams")
//	}
func (c *Client) Divisions(ctx context.Context, q DivisionsQuery) (*Divisions, error) {
	params := &gen.GetDivisionsParams{}
	if q.DivisionID != 0 {
		params.DivisionId = ptr(q.DivisionID)
	}
	if q.LeagueID != 0 {
		params.LeagueId = ptr(q.LeagueID)
	}
	if q.SportID != 0 {
		params.SportId = ptr(q.SportID)
	}
	if q.Season != 0 {
		params.Season = ptr(q.Season)
	}
//
	resp, err := c.raw.GetDivisionsWithResponse(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("mlb: divisions: %w", err)
	}
	if resp.StatusCode() == 404 {
		return nil, ErrNotFound
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("mlb: divisions: unexpected status %d", resp.StatusCode())
	}
	return divisionsFromGen(resp.JSON200), nil
}
//
func divisionsFromGen(r *gen.DivisionsResponse) *Divisions {
	out := &Divisions{}
	if r == nil || r.Divisions == nil {
		return out
	}
	out.Divisions = make([]Division, 0, len(*r.Divisions))
	for _, d := range *r.Divisions {
		out.Divisions = append(out.Divisions, divisionFromGen(d))
	}
	return out
}
//
func divisionFromGen(d gen.Division) Division {
	out := Division{}
	if d.Id != nil {
		out.ID = *d.Id
	}
	if d.Name != nil {
		out.Name = *d.Name
	}
	if d.Season != nil {
		out.Season = *d.Season
	}
	if d.NameShort != nil {
		out.NameShort = *d.NameShort
	}
	if d.Link != nil {
		out.Link = *d.Link
	}
	if d.Abbreviation != nil {
		out.Abbreviation = *d.Abbreviation
	}
	if d.League != nil {
		out.League = refFromGen(d.League)
	}
	if d.Sport != nil {
		out.Sport = refFromGen(d.Sport)
	}
	if d.HasWildcard != nil {
		out.HasWildcard = *d.HasWildcard
	}
	if d.SortOrder != nil {
		out.SortOrder = *d.SortOrder
	}
	if d.NumPlayoffTeams != nil {
		out.NumPlayoffTeams = *d.NumPlayoffTeams
	}
	if d.Active != nil {
		out.Active = *d.Active
	}
	return out
}
