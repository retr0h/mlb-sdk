// Copyright (c) 2026 John Dewey
//
// SPDX-License-Identifier: MIT

package mlb

import (
	"context"
	"fmt"

	"github.com/retr0h/mlb-sdk/internal/gen"
)

// UmpireGamesQuery filters an umpire-games lookup. Season is required.
type UmpireGamesQuery struct {
	Season int // required
	Fields string
}

// UmpireGames fetches games assigned to an umpire for a season.
func (c *Client) UmpireGames(
	ctx context.Context,
	umpireID int,
	q UmpireGamesQuery,
) ([]Game, error) {
	if q.Season == 0 {
		return nil, fmt.Errorf("mlb: umpireGames: %w: Season is required", ErrInvalidQuery)
	}
	params := &gen.GetUmpireGamesParams{Season: q.Season}
	if q.Fields != "" {
		params.Fields = ptr(q.Fields)
	}

	resp, err := c.raw.GetUmpireGamesWithResponse(ctx, umpireID, params)
	if err != nil {
		return nil, fmt.Errorf("mlb: umpireGames: %w", err)
	}
	if resp.StatusCode() == 404 {
		return nil, ErrNotFound
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("mlb: umpireGames: unexpected status %d", resp.StatusCode())
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
