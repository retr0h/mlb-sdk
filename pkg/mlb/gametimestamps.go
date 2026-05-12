// Copyright (c) 2026 John Dewey
//
// SPDX-License-Identifier: MIT

package mlb

import (
	"context"
	"fmt"

	"github.com/retr0h/mlb-sdk/internal/gen"
)

// GameTimestamps fetches the live-feed timestamps for a game.
//
// Example:
//
//	ts, _ := c.GameTimestamps(ctx, 745455)
//	for _, t := range ts { fmt.Println(t) }
func (c *Client) GameTimestamps(ctx context.Context, gamePk int) ([]string, error) {
	resp, err := c.raw.GetGameTimestampsWithResponse(ctx, gamePk)
	if err != nil {
		return nil, fmt.Errorf("mlb: gameTimestamps: %w", err)
	}
	if resp.StatusCode() == 404 {
		return nil, ErrNotFound
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("mlb: gameTimestamps: unexpected status %d", resp.StatusCode())
	}
	return timestampsFromGen(resp.JSON200), nil
}

// GameColorTimestamps fetches the color-feed timestamps for a game.
func (c *Client) GameColorTimestamps(ctx context.Context, gamePk int) ([]string, error) {
	resp, err := c.raw.GetGameColorTimestampsWithResponse(ctx, gamePk)
	if err != nil {
		return nil, fmt.Errorf("mlb: gameColorTimestamps: %w", err)
	}
	if resp.StatusCode() == 404 {
		return nil, ErrNotFound
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf(
			"mlb: gameColorTimestamps: unexpected status %d",
			resp.StatusCode(),
		)
	}
	return timestampsFromGen(resp.JSON200), nil
}

func timestampsFromGen(r *gen.TimestampsResponse) []string {
	if r == nil || r.Timestamps == nil {
		return nil
	}
	return *r.Timestamps
}
