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
// Sports lists the sports tracked by the MLB Stats API. An empty query
// returns every sport (MLB plus affiliated minor / college / independent
// leagues). Filter to a single sport via q.SportID.
//
// Example:
//
//	s, _ := c.Sports(ctx, mlb.SportsQuery{})
//	if mlb := s.Sport(1); mlb != nil {
//	    fmt.Println(mlb.Name, mlb.Abbreviation)
//	}
func (c *Client) Sports(ctx context.Context, q SportsQuery) (*Sports, error) {
	params := &gen.GetSportsParams{}
	if q.SportID != 0 {
		params.SportId = ptr(q.SportID)
	}
	if q.Fields != "" {
		params.Fields = ptr(q.Fields)
	}
//
	resp, err := c.raw.GetSportsWithResponse(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("mlb: sports: %w", err)
	}
	if resp.StatusCode() == 404 {
		return nil, ErrNotFound
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("mlb: sports: unexpected status %d", resp.StatusCode())
	}
	return sportsFromGen(resp.JSON200), nil
}
//
func sportsFromGen(r *gen.SportsResponse) *Sports {
	out := &Sports{}
	if r == nil || r.Sports == nil {
		return out
	}
	out.Sports = make([]Sport, 0, len(*r.Sports))
	for _, s := range *r.Sports {
		out.Sports = append(out.Sports, sportFromGen(s))
	}
	return out
}
//
func sportFromGen(s gen.Sport) Sport {
	out := Sport{}
	if s.Id != nil {
		out.ID = *s.Id
	}
	if s.Code != nil {
		out.Code = *s.Code
	}
	if s.Link != nil {
		out.Link = *s.Link
	}
	if s.Name != nil {
		out.Name = *s.Name
	}
	if s.Abbreviation != nil {
		out.Abbreviation = *s.Abbreviation
	}
	if s.SortOrder != nil {
		out.SortOrder = *s.SortOrder
	}
	if s.ActiveStatus != nil {
		out.ActiveStatus = *s.ActiveStatus
	}
	return out
}
