// Copyright (c) 2026 John Dewey
//
// SPDX-License-Identifier: MIT

package mlb

import (
	"context"
	"fmt"

	"github.com/retr0h/mlb-sdk/internal/gen"
)

// TeamsAffiliatesQuery filters an affiliates lookup. TeamIDs is required.
type TeamsAffiliatesQuery struct {
	TeamIDs string // required, comma-separated
	SportID int
	Season  int
	Hydrate string
	Fields  string
}

// TeamsHistoryQuery filters a teams-history lookup. TeamIDs is required.
type TeamsHistoryQuery struct {
	TeamIDs     string // required, comma-separated
	StartSeason int
	EndSeason   int
	Fields      string
}

// TeamsAffiliates fetches affiliate teams for the given team ids.
func (c *Client) TeamsAffiliates(
	ctx context.Context,
	q TeamsAffiliatesQuery,
) (*Teams, error) {
	if q.TeamIDs == "" {
		return nil, fmt.Errorf("mlb: teamsAffiliates: %w: TeamIDs is required", ErrInvalidQuery)
	}
	params := &gen.GetTeamsAffiliatesParams{TeamIds: q.TeamIDs}
	if q.SportID != 0 {
		params.SportId = ptr(q.SportID)
	}
	if q.Season != 0 {
		params.Season = ptr(q.Season)
	}
	if q.Hydrate != "" {
		params.Hydrate = ptr(q.Hydrate)
	}
	if q.Fields != "" {
		params.Fields = ptr(q.Fields)
	}

	resp, err := c.raw.GetTeamsAffiliatesWithResponse(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("mlb: teamsAffiliates: %w", err)
	}
	if resp.StatusCode() == 404 {
		return nil, ErrNotFound
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("mlb: teamsAffiliates: unexpected status %d", resp.StatusCode())
	}
	return teamsFromGen(resp.JSON200), nil
}

// TeamsHistory fetches historical team records for the given ids.
func (c *Client) TeamsHistory(
	ctx context.Context,
	q TeamsHistoryQuery,
) (*Teams, error) {
	if q.TeamIDs == "" {
		return nil, fmt.Errorf("mlb: teamsHistory: %w: TeamIDs is required", ErrInvalidQuery)
	}
	params := &gen.GetTeamsHistoryParams{TeamIds: q.TeamIDs}
	if q.StartSeason != 0 {
		params.StartSeason = ptr(q.StartSeason)
	}
	if q.EndSeason != 0 {
		params.EndSeason = ptr(q.EndSeason)
	}
	if q.Fields != "" {
		params.Fields = ptr(q.Fields)
	}

	resp, err := c.raw.GetTeamsHistoryWithResponse(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("mlb: teamsHistory: %w", err)
	}
	if resp.StatusCode() == 404 {
		return nil, ErrNotFound
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("mlb: teamsHistory: unexpected status %d", resp.StatusCode())
	}
	return teamsFromGen(resp.JSON200), nil
}

// ScheduleTiedQuery filters a tied-games lookup. Season is required.
type ScheduleTiedQuery struct {
	Season    int // required
	GameTypes string
	Hydrate   string
	Fields    string
}

// ScheduleTied fetches tied/suspended games for a season.
func (c *Client) ScheduleTied(
	ctx context.Context,
	q ScheduleTiedQuery,
) ([]Game, error) {
	if q.Season == 0 {
		return nil, fmt.Errorf("mlb: scheduleTied: %w: Season is required", ErrInvalidQuery)
	}
	params := &gen.GetScheduleTiedParams{Season: q.Season}
	if q.GameTypes != "" {
		params.GameTypes = ptr(q.GameTypes)
	}
	if q.Hydrate != "" {
		params.Hydrate = ptr(q.Hydrate)
	}
	if q.Fields != "" {
		params.Fields = ptr(q.Fields)
	}

	resp, err := c.raw.GetScheduleTiedWithResponse(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("mlb: scheduleTied: %w", err)
	}
	if resp.StatusCode() == 404 {
		return nil, ErrNotFound
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("mlb: scheduleTied: unexpected status %d", resp.StatusCode())
	}
	var games []Game
	if resp.JSON200.Dates != nil {
		for _, d := range *resp.JSON200.Dates {
			if d.Games != nil {
				for _, g := range *d.Games {
					games = append(games, gameFromGen(g))
				}
			}
		}
	}
	return games, nil
}
