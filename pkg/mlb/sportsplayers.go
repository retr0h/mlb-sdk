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
// SportsPlayersQuery filters a sports-players lookup. Season is required.
type SportsPlayersQuery struct {
	Season   int // required
	GameType string
	Fields   string
}
//
// SportsPlayers fetches all players for a sport + season.
func (c *Client) SportsPlayers(
	ctx context.Context,
	sportID int,
	q SportsPlayersQuery,
) ([]PersonDetail, error) {
	if q.Season == 0 {
		return nil, fmt.Errorf("mlb: sportsPlayers: %w: Season is required", ErrInvalidQuery)
	}
	params := &gen.GetSportsPlayersParams{Season: q.Season}
	if q.GameType != "" {
		params.GameType = ptr(q.GameType)
	}
	if q.Fields != "" {
		params.Fields = ptr(q.Fields)
	}
//
	resp, err := c.raw.GetSportsPlayersWithResponse(ctx, sportID, params)
	if err != nil {
		return nil, fmt.Errorf("mlb: sportsPlayers: %w", err)
	}
	if resp.StatusCode() == 404 {
		return nil, ErrNotFound
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("mlb: sportsPlayers: unexpected status %d", resp.StatusCode())
	}
	return peopleFromGen(resp.JSON200), nil
}
