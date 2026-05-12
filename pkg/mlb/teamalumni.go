// Copyright (c) 2026 John Dewey
//
// SPDX-License-Identifier: MIT

package mlb

import (
	"context"
	"fmt"

	"github.com/retr0h/mlb-sdk/internal/gen"
)

// TeamAlumniQuery filters an alumni lookup. Season and Group are required.
type TeamAlumniQuery struct {
	Season  int    // required
	Group   string // required: "hitting" | "pitching" | "fielding"
	Hydrate string
	Fields  string
}

// TeamAlumni fetches alumni for a team. Season and Group are required.
func (c *Client) TeamAlumni(
	ctx context.Context,
	teamID int,
	q TeamAlumniQuery,
) ([]PersonDetail, error) {
	if q.Season == 0 || q.Group == "" {
		return nil, fmt.Errorf(
			"mlb: teamAlumni: %w: Season and Group are both required",
			ErrInvalidQuery,
		)
	}
	params := &gen.GetTeamAlumniParams{Season: q.Season, Group: q.Group}
	if q.Hydrate != "" {
		params.Hydrate = ptr(q.Hydrate)
	}
	if q.Fields != "" {
		params.Fields = ptr(q.Fields)
	}

	resp, err := c.raw.GetTeamAlumniWithResponse(ctx, teamID, params)
	if err != nil {
		return nil, fmt.Errorf("mlb: teamAlumni: %w", err)
	}
	if resp.StatusCode() == 404 {
		return nil, ErrNotFound
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("mlb: teamAlumni: unexpected status %d", resp.StatusCode())
	}
	return peopleFromGen(resp.JSON200), nil
}
