// Copyright (c) 2026 John Dewey
//
// SPDX-License-Identifier: MIT

package mlb

import (
	"context"
	"fmt"

	"github.com/retr0h/mlb-sdk/internal/gen"
)

// FreeAgents fetches free-agent signings and declarations.
//
// Example:
//
//	fa, _ := c.FreeAgents(ctx, mlb.FreeAgentsQuery{Season: 2024})
//	for _, f := range fa.FreeAgents {
//	    fmt.Println(f.Player.FullName, f.OriginalTeam.Name, "→", f.NewTeam.Name)
//	}
func (c *Client) FreeAgents(ctx context.Context, q FreeAgentsQuery) (*FreeAgents, error) {
	params := &gen.GetFreeAgentsParams{}
	if q.Season != 0 {
		params.Season = ptr(q.Season)
	}
	if q.Order != "" {
		params.Order = ptr(q.Order)
	}
	if q.Hydrate != "" {
		params.Hydrate = ptr(q.Hydrate)
	}
	if q.Fields != "" {
		params.Fields = ptr(q.Fields)
	}

	resp, err := c.raw.GetFreeAgentsWithResponse(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("mlb: freeAgents: %w", err)
	}
	if resp.StatusCode() == 404 {
		return nil, ErrNotFound
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("mlb: freeAgents: unexpected status %d", resp.StatusCode())
	}
	return freeAgentsFromGen(resp.JSON200), nil
}

func freeAgentsFromGen(r *gen.FreeAgentsResponse) *FreeAgents {
	out := &FreeAgents{}
	if r == nil || r.FreeAgents == nil {
		return out
	}
	out.FreeAgents = make([]FreeAgent, 0, len(*r.FreeAgents))
	for _, fa := range *r.FreeAgents {
		out.FreeAgents = append(out.FreeAgents, freeAgentFromGen(fa))
	}
	return out
}

func freeAgentFromGen(f gen.FreeAgent) FreeAgent {
	out := FreeAgent{}
	if f.Player != nil {
		out.Player = personFromGen(*f.Player)
	}
	if f.OriginalTeam != nil {
		if f.OriginalTeam.Id != nil {
			out.OriginalTeam.ID = TeamID(*f.OriginalTeam.Id)
		}
		if f.OriginalTeam.Name != nil {
			out.OriginalTeam.Name = *f.OriginalTeam.Name
		}
	}
	if f.NewTeam != nil {
		if f.NewTeam.Id != nil {
			out.NewTeam.ID = TeamID(*f.NewTeam.Id)
		}
		if f.NewTeam.Name != nil {
			out.NewTeam.Name = *f.NewTeam.Name
		}
	}
	if f.Notes != nil {
		out.Notes = *f.Notes
	}
	if f.DateSigned != nil {
		out.DateSigned = *f.DateSigned
	}
	if f.DateDeclared != nil {
		out.DateDeclared = *f.DateDeclared
	}
	if f.Position != nil {
		out.Position = primaryPositionFromGen(*f.Position)
	}
	return out
}
