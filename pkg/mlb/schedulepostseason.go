// Copyright (c) 2026 John Dewey
//
// SPDX-License-Identifier: MIT

package mlb

import (
	"context"
	"fmt"

	"github.com/retr0h/mlb-sdk/internal/gen"
)

// SchedulePostseasonQuery filters a postseason schedule lookup.
type SchedulePostseasonQuery struct {
	Season       int
	GameTypes    string
	SeriesNumber int
	TeamID       int
	SportID      int
	Hydrate      string
	Fields       string
}

// SchedulePostseason fetches the postseason schedule.
func (c *Client) SchedulePostseason(
	ctx context.Context,
	q SchedulePostseasonQuery,
) ([]Game, error) {
	params := &gen.GetSchedulePostseasonParams{}
	if q.Season != 0 {
		params.Season = ptr(q.Season)
	}
	if q.GameTypes != "" {
		params.GameTypes = ptr(q.GameTypes)
	}
	if q.SeriesNumber != 0 {
		params.SeriesNumber = ptr(q.SeriesNumber)
	}
	if q.TeamID != 0 {
		params.TeamId = ptr(q.TeamID)
	}
	if q.SportID != 0 {
		params.SportId = ptr(q.SportID)
	}
	if q.Hydrate != "" {
		params.Hydrate = ptr(q.Hydrate)
	}
	if q.Fields != "" {
		params.Fields = ptr(q.Fields)
	}

	resp, err := c.raw.GetSchedulePostseasonWithResponse(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("mlb: schedulePostseason: %w", err)
	}
	if resp.StatusCode() == 404 {
		return nil, ErrNotFound
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("mlb: schedulePostseason: unexpected status %d", resp.StatusCode())
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

// SchedulePostseasonTuneInQuery filters a postseason tune-in lookup.
type SchedulePostseasonTuneInQuery struct {
	Season  int
	TeamID  int
	SportID int
	Hydrate string
	Fields  string
}

// SchedulePostseasonTuneIn fetches postseason tune-in info.
func (c *Client) SchedulePostseasonTuneIn(
	ctx context.Context,
	q SchedulePostseasonTuneInQuery,
) ([]Game, error) {
	params := &gen.GetSchedulePostseasonTuneInParams{}
	if q.Season != 0 {
		params.Season = ptr(q.Season)
	}
	if q.TeamID != 0 {
		params.TeamId = ptr(q.TeamID)
	}
	if q.SportID != 0 {
		params.SportId = ptr(q.SportID)
	}
	if q.Hydrate != "" {
		params.Hydrate = ptr(q.Hydrate)
	}
	if q.Fields != "" {
		params.Fields = ptr(q.Fields)
	}

	resp, err := c.raw.GetSchedulePostseasonTuneInWithResponse(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("mlb: schedulePostseasonTuneIn: %w", err)
	}
	if resp.StatusCode() == 404 {
		return nil, ErrNotFound
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf(
			"mlb: schedulePostseasonTuneIn: unexpected status %d",
			resp.StatusCode(),
		)
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
