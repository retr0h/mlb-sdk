// Copyright (c) 2026 John Dewey
//
// SPDX-License-Identifier: MIT

package mlb

import (
	"context"
	"fmt"

	"github.com/retr0h/mlb-sdk/internal/gen"
)

// Teams lists teams the MLB API tracks. An empty query returns every team
// across every sport; filter via TeamsQuery.
//
// Example:
//
//	all, _ := c.Teams(ctx, mlb.TeamsQuery{SportID: 1, Season: 2024})
//	for _, t := range all.Teams {
//	    fmt.Println(t.Abbreviation, t.Name)
//	}
func (c *Client) Teams(ctx context.Context, q TeamsQuery) (*Teams, error) {
	params := &gen.GetTeamsParams{}
	if q.Season != 0 {
		params.Season = ptr(q.Season)
	}
	if q.ActiveStatus != "" {
		params.ActiveStatus = ptr(q.ActiveStatus)
	}
	if q.LeagueIDs != "" {
		params.LeagueIds = ptr(q.LeagueIDs)
	}
	if q.SportID != 0 {
		params.SportId = ptr(q.SportID)
	}
	if q.SportIDs != "" {
		params.SportIds = ptr(q.SportIDs)
	}
	if q.GameType != "" {
		params.GameType = ptr(q.GameType)
	}
	if q.Hydrate != "" {
		params.Hydrate = ptr(q.Hydrate)
	}
	if q.Fields != "" {
		params.Fields = ptr(q.Fields)
	}

	resp, err := c.raw.GetTeamsWithResponse(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("mlb: teams: %w", err)
	}
	if resp.StatusCode() == 404 {
		return nil, ErrNotFound
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("mlb: teams: unexpected status %d", resp.StatusCode())
	}
	return teamsFromGen(resp.JSON200), nil
}

func teamsFromGen(r *gen.TeamsResponse) *Teams {
	out := &Teams{}
	if r == nil || r.Teams == nil {
		return out
	}
	out.Teams = make([]TeamInfo, 0, len(*r.Teams))
	for _, t := range *r.Teams {
		out.Teams = append(out.Teams, teamInfoFromGen(t))
	}
	return out
}
