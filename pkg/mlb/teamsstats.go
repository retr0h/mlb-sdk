// Copyright (c) 2026 John Dewey
//
// SPDX-License-Identifier: MIT

package mlb

import (
	"context"
	"fmt"

	"github.com/retr0h/mlb-sdk/internal/gen"
)

// TeamsStatsQuery filters a league-wide teams-stats lookup. Season, Group,
// and Stats are all required.
type TeamsStatsQuery struct {
	Season    int    // required
	Group     string // required: "hitting" | "pitching" | "fielding"
	Stats     string // required: "season" | "byDateRange" | …
	SportIDs  string
	GameType  string
	Order     string
	SortStat  string
	StartDate string
	EndDate   string
	Fields    string
}

// TeamsStats fetches league-wide aggregated team stats.
func (c *Client) TeamsStats(ctx context.Context, q TeamsStatsQuery) (*TeamStats, error) {
	if q.Season == 0 || q.Group == "" || q.Stats == "" {
		return nil, fmt.Errorf(
			"mlb: teamsStats: %w: Season, Group, and Stats are all required",
			ErrInvalidQuery,
		)
	}
	params := &gen.GetTeamsStatsParams{
		Season: q.Season,
		Group:  q.Group,
		Stats:  q.Stats,
	}
	if q.SportIDs != "" {
		params.SportIds = ptr(q.SportIDs)
	}
	if q.GameType != "" {
		params.GameType = ptr(q.GameType)
	}
	if q.Order != "" {
		params.Order = ptr(q.Order)
	}
	if q.SortStat != "" {
		params.SortStat = ptr(q.SortStat)
	}
	if q.StartDate != "" {
		params.StartDate = ptr(q.StartDate)
	}
	if q.EndDate != "" {
		params.EndDate = ptr(q.EndDate)
	}
	if q.Fields != "" {
		params.Fields = ptr(q.Fields)
	}

	resp, err := c.raw.GetTeamsStatsWithResponse(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("mlb: teamsStats: %w", err)
	}
	if resp.StatusCode() == 404 {
		return nil, ErrNotFound
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("mlb: teamsStats: unexpected status %d", resp.StatusCode())
	}
	return teamStatsFromGen(resp.JSON200), nil
}
