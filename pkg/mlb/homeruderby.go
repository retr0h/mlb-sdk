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
// HomeRunDerbyQuery refines a derby lookup.
type HomeRunDerbyQuery struct {
	Fields string
}
//
// HomeRunDerby fetches Home Run Derby data for a game.
func (c *Client) HomeRunDerby(
	ctx context.Context,
	gamePk int,
	q HomeRunDerbyQuery,
) (map[string]any, error) {
	params := &gen.GetHomeRunDerbyParams{}
	if q.Fields != "" {
		params.Fields = ptr(q.Fields)
	}
//
	resp, err := c.raw.GetHomeRunDerbyWithResponse(ctx, gamePk, params)
	if err != nil {
		return nil, fmt.Errorf("mlb: homeRunDerby: %w", err)
	}
	if resp.StatusCode() == 404 {
		return nil, ErrNotFound
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("mlb: homeRunDerby: unexpected status %d", resp.StatusCode())
	}
	return map[string]any(*resp.JSON200), nil
}
