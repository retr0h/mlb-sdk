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
// Season fetches one season's metadata by id. q.SportID is required; the
// MLB API rejects the call otherwise (toddrob99 encodes this as
// `required_params: [["sportId"]]`).
//
// The underlying endpoint responds with `{"seasons": [<one entry>]}`. This
// method collapses the wrapper and returns the single Season; an empty
// array maps to ErrNotFound.
//
// Example:
//
//	s, _ := c.Season(ctx, "2024", mlb.SeasonQuery{SportID: 1})
//	fmt.Println(s.RegularSeasonStartDate, "→", s.RegularSeasonEndDate)
func (c *Client) Season(ctx context.Context, seasonID string, q SeasonQuery) (*Season, error) {
	if q.SportID == 0 {
		return nil, fmt.Errorf("mlb: season: %w: SportID is required", ErrInvalidQuery)
	}
	params := &gen.GetSeasonParams{SportId: q.SportID}
	if q.Fields != "" {
		params.Fields = ptr(q.Fields)
	}
//
	resp, err := c.raw.GetSeasonWithResponse(ctx, seasonID, params)
	if err != nil {
		return nil, fmt.Errorf("mlb: season: %w", err)
	}
	if resp.StatusCode() == 404 {
		return nil, ErrNotFound
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("mlb: season: unexpected status %d", resp.StatusCode())
	}
	s := seasonFromResponse(resp.JSON200)
	if s == nil {
		return nil, ErrNotFound
	}
	return s, nil
}
//
// seasonFromResponse collapses the seasons array to a single Season. Empty
// array → nil so the caller maps to ErrNotFound.
func seasonFromResponse(r *gen.SeasonsResponse) *Season {
	if r == nil || r.Seasons == nil || len(*r.Seasons) == 0 {
		return nil
	}
	s := seasonFromGen((*r.Seasons)[0])
	return &s
}
