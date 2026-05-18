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
// TeamLeadersQuery filters a team-leaders lookup. LeaderCategories and
// Season are required.
type TeamLeadersQuery struct {
	LeaderCategories string // required
	Season           int    // required
	LeaderGameTypes  string
	Hydrate          string
	Limit            int
	Fields           string
}
//
// TeamLeaders fetches stat leaders for a team.
func (c *Client) TeamLeaders(
	ctx context.Context,
	teamID int,
	q TeamLeadersQuery,
) (*StatsLeaders, error) {
	if q.LeaderCategories == "" || q.Season == 0 {
		return nil, fmt.Errorf(
			"mlb: teamLeaders: %w: LeaderCategories and Season are both required",
			ErrInvalidQuery,
		)
	}
	params := &gen.GetTeamLeadersParams{
		LeaderCategories: q.LeaderCategories,
		Season:           q.Season,
	}
	if q.LeaderGameTypes != "" {
		params.LeaderGameTypes = ptr(q.LeaderGameTypes)
	}
	if q.Hydrate != "" {
		params.Hydrate = ptr(q.Hydrate)
	}
	if q.Limit != 0 {
		params.Limit = ptr(q.Limit)
	}
	if q.Fields != "" {
		params.Fields = ptr(q.Fields)
	}
//
	resp, err := c.raw.GetTeamLeadersWithResponse(ctx, teamID, params)
	if err != nil {
		return nil, fmt.Errorf("mlb: teamLeaders: %w", err)
	}
	if resp.StatusCode() == 404 {
		return nil, ErrNotFound
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("mlb: teamLeaders: unexpected status %d", resp.StatusCode())
	}
	return teamLeadersFromGen(resp.JSON200), nil
}
//
func teamLeadersFromGen(r *gen.TeamLeadersResponse) *StatsLeaders {
	out := &StatsLeaders{}
	if r == nil || r.TeamLeaders == nil {
		return out
	}
	out.LeagueLeaders = make([]LeaderCategory, 0, len(*r.TeamLeaders))
	for _, c := range *r.TeamLeaders {
		out.LeagueLeaders = append(out.LeagueLeaders, leaderCategoryFromGen(c))
	}
	return out
}
