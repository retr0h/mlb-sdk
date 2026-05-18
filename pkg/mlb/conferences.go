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
// Conferences lists the conferences tracked by the MLB Stats API.
//
// Example:
//
//	c, _ := client.Conferences(ctx, mlb.ConferencesQuery{})
//	for _, conf := range c.Conferences {
//	    fmt.Println(conf.Name, conf.Abbreviation)
//	}
func (c *Client) Conferences(ctx context.Context, q ConferencesQuery) (*Conferences, error) {
	params := &gen.GetConferencesParams{}
	if q.ConferenceID != 0 {
		params.ConferenceId = ptr(q.ConferenceID)
	}
	if q.Season != 0 {
		params.Season = ptr(q.Season)
	}
	if q.Fields != "" {
		params.Fields = ptr(q.Fields)
	}
//
	resp, err := c.raw.GetConferencesWithResponse(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("mlb: conferences: %w", err)
	}
	if resp.StatusCode() == 404 {
		return nil, ErrNotFound
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("mlb: conferences: unexpected status %d", resp.StatusCode())
	}
	return conferencesFromGen(resp.JSON200), nil
}
//
func conferencesFromGen(r *gen.ConferencesResponse) *Conferences {
	out := &Conferences{}
	if r == nil || r.Conferences == nil {
		return out
	}
	out.Conferences = make([]Conference, 0, len(*r.Conferences))
	for _, c := range *r.Conferences {
		out.Conferences = append(out.Conferences, conferenceFromGen(c))
	}
	return out
}
//
func conferenceFromGen(c gen.Conference) Conference {
	out := Conference{}
	if c.Id != nil {
		out.ID = *c.Id
	}
	if c.Name != nil {
		out.Name = *c.Name
	}
	if c.Link != nil {
		out.Link = *c.Link
	}
	if c.Abbreviation != nil {
		out.Abbreviation = *c.Abbreviation
	}
	if c.NameShort != nil {
		out.NameShort = *c.NameShort
	}
	if c.HasWildcard != nil {
		out.HasWildcard = *c.HasWildcard
	}
	if c.League != nil {
		out.League = refFromGen(c.League)
	}
	if c.Sport != nil {
		out.Sport = refFromGen(c.Sport)
	}
	return out
}
