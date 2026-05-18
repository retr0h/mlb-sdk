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
// GameColorQuery refines a color-feed lookup.
type GameColorQuery struct {
	Timecode string
	Fields   string
}
//
// GameColor fetches the color commentary feed. The response is opaque
// (map[string]any) because the full response shape is too large to model.
func (c *Client) GameColor(
	ctx context.Context,
	gamePk int,
	q GameColorQuery,
) (map[string]any, error) {
	params := &gen.GetGameColorParams{}
	if q.Timecode != "" {
		params.Timecode = ptr(q.Timecode)
	}
	if q.Fields != "" {
		params.Fields = ptr(q.Fields)
	}
//
	resp, err := c.raw.GetGameColorWithResponse(ctx, gamePk, params)
	if err != nil {
		return nil, fmt.Errorf("mlb: gameColor: %w", err)
	}
	if resp.StatusCode() == 404 {
		return nil, ErrNotFound
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("mlb: gameColor: unexpected status %d", resp.StatusCode())
	}
	return map[string]any(*resp.JSON200), nil
}
//
// GameColorDiffQuery refines a color-diff lookup. Both timecodes required.
type GameColorDiffQuery struct {
	StartTimecode string // required
	EndTimecode   string // required
}
//
// GameColorDiff fetches the color-feed diff patch between two timecodes.
func (c *Client) GameColorDiff(
	ctx context.Context,
	gamePk int,
	q GameColorDiffQuery,
) ([]map[string]any, error) {
	if q.StartTimecode == "" || q.EndTimecode == "" {
		return nil, fmt.Errorf(
			"mlb: gameColorDiff: %w: StartTimecode and EndTimecode are both required",
			ErrInvalidQuery,
		)
	}
	params := &gen.GetGameColorDiffParams{
		StartTimecode: q.StartTimecode,
		EndTimecode:   q.EndTimecode,
	}
	resp, err := c.raw.GetGameColorDiffWithResponse(ctx, gamePk, params)
	if err != nil {
		return nil, fmt.Errorf("mlb: gameColorDiff: %w", err)
	}
	if resp.StatusCode() == 404 {
		return nil, ErrNotFound
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("mlb: gameColorDiff: unexpected status %d", resp.StatusCode())
	}
	out := make([]map[string]any, len(*resp.JSON200))
	copy(out, *resp.JSON200)
	return out, nil
}
//
// GameDiffQuery refines a live-feed diff lookup. Both timecodes required.
type GameDiffQuery struct {
	StartTimecode string // required
	EndTimecode   string // required
}
//
// GameDiff fetches the live-feed diff patch between two timecodes.
func (c *Client) GameDiff(
	ctx context.Context,
	gamePk int,
	q GameDiffQuery,
) ([]map[string]any, error) {
	if q.StartTimecode == "" || q.EndTimecode == "" {
		return nil, fmt.Errorf(
			"mlb: gameDiff: %w: StartTimecode and EndTimecode are both required",
			ErrInvalidQuery,
		)
	}
	params := &gen.GetGameDiffParams{
		StartTimecode: q.StartTimecode,
		EndTimecode:   q.EndTimecode,
	}
	resp, err := c.raw.GetGameDiffWithResponse(ctx, gamePk, params)
	if err != nil {
		return nil, fmt.Errorf("mlb: gameDiff: %w", err)
	}
	if resp.StatusCode() == 404 {
		return nil, ErrNotFound
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("mlb: gameDiff: unexpected status %d", resp.StatusCode())
	}
	out := make([]map[string]any, len(*resp.JSON200))
	copy(out, *resp.JSON200)
	return out, nil
}
//
// GameContentQuery refines a game-content lookup.
type GameContentQuery struct {
	HighlightLimit int
}
//
// GameContent fetches game content (highlights, editorial). The response
// is opaque because the full shape is deeply nested.
func (c *Client) GameContent(
	ctx context.Context,
	gamePk int,
	q GameContentQuery,
) (map[string]any, error) {
	params := &gen.GetGameContentParams{}
	if q.HighlightLimit != 0 {
		params.HighlightLimit = ptr(q.HighlightLimit)
	}
//
	resp, err := c.raw.GetGameContentWithResponse(ctx, gamePk, params)
	if err != nil {
		return nil, fmt.Errorf("mlb: gameContent: %w", err)
	}
	if resp.StatusCode() == 404 {
		return nil, ErrNotFound
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("mlb: gameContent: unexpected status %d", resp.StatusCode())
	}
	return map[string]any(*resp.JSON200), nil
}
//
// GameWinProbabilityQuery refines a win-probability lookup.
type GameWinProbabilityQuery struct {
	Timecode string
	Fields   string
}
//
// GameWinProbability fetches win probability per at-bat. Returns opaque
// maps because each entry is a full play object with probability fields.
func (c *Client) GameWinProbability(
	ctx context.Context,
	gamePk int,
	q GameWinProbabilityQuery,
) ([]map[string]any, error) {
	params := &gen.GetGameWinProbabilityParams{}
	if q.Timecode != "" {
		params.Timecode = ptr(q.Timecode)
	}
	if q.Fields != "" {
		params.Fields = ptr(q.Fields)
	}
//
	resp, err := c.raw.GetGameWinProbabilityWithResponse(ctx, gamePk, params)
	if err != nil {
		return nil, fmt.Errorf("mlb: gameWinProbability: %w", err)
	}
	if resp.StatusCode() == 404 {
		return nil, ErrNotFound
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("mlb: gameWinProbability: unexpected status %d", resp.StatusCode())
	}
	out := make([]map[string]any, len(*resp.JSON200))
	copy(out, *resp.JSON200)
	return out, nil
}
//
// Meta fetches API metadata (gameTypes, statGroups, etc.). The metaType
// path parameter determines the resource — e.g. "gameTypes", "statGroups".
func (c *Client) Meta(ctx context.Context, metaType string) ([]map[string]any, error) {
	resp, err := c.raw.GetMetaWithResponse(ctx, metaType)
	if err != nil {
		return nil, fmt.Errorf("mlb: meta: %w", err)
	}
	if resp.StatusCode() == 404 {
		return nil, ErrNotFound
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("mlb: meta: unexpected status %d", resp.StatusCode())
	}
	out := make([]map[string]any, len(*resp.JSON200))
	copy(out, *resp.JSON200)
	return out, nil
}
