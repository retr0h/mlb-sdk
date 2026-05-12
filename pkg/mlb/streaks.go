// Copyright (c) 2026 John Dewey
//
// SPDX-License-Identifier: MIT

package mlb

import (
	"context"
	"fmt"

	"github.com/retr0h/mlb-sdk/internal/gen"
)

// StatsStreaksQuery filters a stat-streaks lookup. All five fields are
// required (toddrob99: required_params=[["streakType", "streakSpan",
// "season", "sportId", "limit"]]).
type StatsStreaksQuery struct {
	StreakType string // required
	StreakSpan string // required
	Season     int    // required
	SportID    int    // required
	Limit      int    // required
	GameType   string
	Hydrate    string
	Fields     string
}

// StatsStreaks fetches stat streak data. All five required fields must be
// set. Note: this endpoint may return 404 during the offseason.
func (c *Client) StatsStreaks(
	ctx context.Context,
	q StatsStreaksQuery,
) (map[string]any, error) {
	if q.StreakType == "" || q.StreakSpan == "" || q.Season == 0 ||
		q.SportID == 0 || q.Limit == 0 {
		return nil, fmt.Errorf(
			"mlb: statsStreaks: %w: StreakType, StreakSpan, Season, SportID, and Limit are all required",
			ErrInvalidQuery,
		)
	}
	params := &gen.GetStatsStreaksParams{
		StreakType: q.StreakType,
		StreakSpan: q.StreakSpan,
		Season:     q.Season,
		SportId:    q.SportID,
		Limit:      q.Limit,
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

	resp, err := c.raw.GetStatsStreaksWithResponse(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("mlb: statsStreaks: %w", err)
	}
	if resp.StatusCode() == 404 {
		return nil, ErrNotFound
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("mlb: statsStreaks: unexpected status %d", resp.StatusCode())
	}
	return map[string]any(*resp.JSON200), nil
}
