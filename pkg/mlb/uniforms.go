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
// GameUniformsQuery filters a game-uniforms lookup. GamePks is required.
type GameUniformsQuery struct {
	GamePks string // required, comma-separated
	Fields  string
}
//
// GameUniforms fetches uniform data for the given games.
func (c *Client) GameUniforms(ctx context.Context, q GameUniformsQuery) ([]map[string]any, error) {
	if q.GamePks == "" {
		return nil, fmt.Errorf("mlb: gameUniforms: %w: GamePks is required", ErrInvalidQuery)
	}
	params := &gen.GetGameUniformsParams{GamePks: q.GamePks}
	if q.Fields != "" {
		params.Fields = ptr(q.Fields)
	}
//
	resp, err := c.raw.GetGameUniformsWithResponse(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("mlb: gameUniforms: %w", err)
	}
	if resp.StatusCode() == 404 {
		return nil, ErrNotFound
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("mlb: gameUniforms: unexpected status %d", resp.StatusCode())
	}
	if resp.JSON200.Uniforms == nil {
		return nil, nil
	}
	out := make([]map[string]any, len(*resp.JSON200.Uniforms))
	copy(out, *resp.JSON200.Uniforms)
	return out, nil
}
//
// TeamUniformsQuery filters a team-uniforms lookup. TeamIDs is required.
type TeamUniformsQuery struct {
	TeamIDs string // required, comma-separated
	Season  int
	Fields  string
}
//
// TeamUniforms fetches the uniform catalog for the given teams.
func (c *Client) TeamUniforms(ctx context.Context, q TeamUniformsQuery) ([]map[string]any, error) {
	if q.TeamIDs == "" {
		return nil, fmt.Errorf("mlb: teamUniforms: %w: TeamIDs is required", ErrInvalidQuery)
	}
	params := &gen.GetTeamUniformsParams{TeamIds: q.TeamIDs}
	if q.Season != 0 {
		params.Season = ptr(q.Season)
	}
	if q.Fields != "" {
		params.Fields = ptr(q.Fields)
	}
//
	resp, err := c.raw.GetTeamUniformsWithResponse(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("mlb: teamUniforms: %w", err)
	}
	if resp.StatusCode() == 404 {
		return nil, ErrNotFound
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("mlb: teamUniforms: unexpected status %d", resp.StatusCode())
	}
	if resp.JSON200.Uniforms == nil {
		return nil, nil
	}
	out := make([]map[string]any, len(*resp.JSON200.Uniforms))
	copy(out, *resp.JSON200.Uniforms)
	return out, nil
}
