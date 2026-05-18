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
// Leagues fetches league metadata. One of q.SportID or q.LeagueIDs must be
// set; the MLB API rejects the call otherwise (toddrob99 encodes this as
// `required_params: [["sportId"], ["leagueIds"]]`). Each LeagueInfo's
// SeasonDateInfo carries the league's calendar windows for that season.
//
// Example:
//
//	ls, _ := c.Leagues(ctx, mlb.LeaguesQuery{SportID: 1})
//	if nl := ls.League(104); nl != nil {
//	    fmt.Println(nl.Name, "→", nl.NumTeams, "teams")
//	}
func (c *Client) Leagues(ctx context.Context, q LeaguesQuery) (*Leagues, error) {
	if q.SportID == 0 && q.LeagueIDs == "" {
		return nil, fmt.Errorf(
			"mlb: leagues: %w: one of SportID or LeagueIDs is required",
			ErrInvalidQuery,
		)
	}
//
	params := &gen.GetLeaguesParams{}
	if q.SportID != 0 {
		params.SportId = ptr(q.SportID)
	}
	if q.LeagueIDs != "" {
		params.LeagueIds = ptr(q.LeagueIDs)
	}
	if q.Seasons != "" {
		params.Seasons = ptr(q.Seasons)
	}
	if q.Fields != "" {
		params.Fields = ptr(q.Fields)
	}
//
	resp, err := c.raw.GetLeaguesWithResponse(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("mlb: leagues: %w", err)
	}
	if resp.StatusCode() == 404 {
		return nil, ErrNotFound
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("mlb: leagues: unexpected status %d", resp.StatusCode())
	}
	return leaguesFromGen(resp.JSON200), nil
}
//
func leaguesFromGen(r *gen.LeaguesResponse) *Leagues {
	out := &Leagues{}
	if r == nil || r.Leagues == nil {
		return out
	}
	out.Leagues = make([]LeagueInfo, 0, len(*r.Leagues))
	for _, l := range *r.Leagues {
		out.Leagues = append(out.Leagues, leagueInfoFromGen(l))
	}
	return out
}
