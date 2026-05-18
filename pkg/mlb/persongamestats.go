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
// PersonGameStatsQuery refines a person-game-stats lookup.
type PersonGameStatsQuery struct {
	Fields string
}
//
// PersonGameStats fetches a player's stats for a specific game.
func (c *Client) PersonGameStats(
	ctx context.Context,
	personID int,
	gamePk int,
	q PersonGameStatsQuery,
) (*TeamStats, error) {
	params := &gen.GetPersonGameStatsParams{}
	if q.Fields != "" {
		params.Fields = ptr(q.Fields)
	}
//
	resp, err := c.raw.GetPersonGameStatsWithResponse(ctx, personID, gamePk, params)
	if err != nil {
		return nil, fmt.Errorf("mlb: personGameStats: %w", err)
	}
	if resp.StatusCode() == 404 {
		return nil, ErrNotFound
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("mlb: personGameStats: unexpected status %d", resp.StatusCode())
	}
	return teamStatsFromGen(resp.JSON200), nil
}
