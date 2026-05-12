// Copyright (c) 2026 John Dewey
//
// SPDX-License-Identifier: MIT

package mlb

import (
	"context"
	"fmt"

	"github.com/retr0h/mlb-sdk/internal/gen"
)

// ContextMetrics fetches win probability and sac-fly probability for a game.
//
// Example:
//
//	cm, _ := c.ContextMetrics(ctx, 745455, mlb.ContextMetricsQuery{})
//	fmt.Printf("home %.1f%% away %.1f%%\n", cm.HomeWinProbability, cm.AwayWinProbability)
func (c *Client) ContextMetrics(
	ctx context.Context,
	gamePk int,
	q ContextMetricsQuery,
) (*ContextMetrics, error) {
	params := &gen.GetContextMetricsParams{}
	if q.Timecode != "" {
		params.Timecode = ptr(q.Timecode)
	}
	if q.Fields != "" {
		params.Fields = ptr(q.Fields)
	}

	resp, err := c.raw.GetContextMetricsWithResponse(ctx, gamePk, params)
	if err != nil {
		return nil, fmt.Errorf("mlb: contextMetrics: %w", err)
	}
	if resp.StatusCode() == 404 {
		return nil, ErrNotFound
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("mlb: contextMetrics: unexpected status %d", resp.StatusCode())
	}
	return contextMetricsFromGen(resp.JSON200), nil
}

func contextMetricsFromGen(r *gen.ContextMetricsResponse) *ContextMetrics {
	out := &ContextMetrics{}
	if r == nil {
		return out
	}
	if r.Game != nil {
		out.Game = refFromGen(r.Game)
	}
	if r.LeftFieldSacFlyProbability != nil {
		out.LeftFieldSacFlyProbability = map[string]any(*r.LeftFieldSacFlyProbability)
	}
	if r.CenterFieldSacFlyProbability != nil {
		out.CenterFieldSacFlyProbability = map[string]any(*r.CenterFieldSacFlyProbability)
	}
	if r.RightFieldSacFlyProbability != nil {
		out.RightFieldSacFlyProbability = map[string]any(*r.RightFieldSacFlyProbability)
	}
	if r.HomeWinProbability != nil {
		out.HomeWinProbability = *r.HomeWinProbability
	}
	if r.AwayWinProbability != nil {
		out.AwayWinProbability = *r.AwayWinProbability
	}
	return out
}
