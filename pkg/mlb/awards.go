// Copyright (c) 2026 John Dewey
//
// SPDX-License-Identifier: MIT

package mlb

import (
	"context"
	"fmt"

	"github.com/retr0h/mlb-sdk/internal/gen"
)

// AwardRecipients fetches every recipient of a single award (e.g. MLBHOF,
// ALMVP). awardId is required; the underlying MLB API returns 404 when the
// id is unknown.
//
// Example:
//
//	a, _ := c.AwardRecipients(ctx, "MLBHOF", mlb.AwardRecipientsQuery{})
//	for _, r := range a.Recipients {
//	    fmt.Println(r.Date, r.Player.NameFirstLast, "—", r.Notes)
//	}
func (c *Client) AwardRecipients(
	ctx context.Context,
	awardID string,
	q AwardRecipientsQuery,
) (*Awards, error) {
	params := &gen.GetAwardRecipientsParams{}
	if q.SportID != 0 {
		params.SportId = ptr(q.SportID)
	}
	if q.LeagueID != 0 {
		params.LeagueId = ptr(q.LeagueID)
	}
	if q.Season != 0 {
		params.Season = ptr(q.Season)
	}
	if q.Hydrate != "" {
		params.Hydrate = ptr(q.Hydrate)
	}
	if q.Fields != "" {
		params.Fields = ptr(q.Fields)
	}

	resp, err := c.raw.GetAwardRecipientsWithResponse(ctx, awardID, params)
	if err != nil {
		return nil, fmt.Errorf("mlb: awardRecipients: %w", err)
	}
	if resp.StatusCode() == 404 {
		return nil, ErrNotFound
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("mlb: awardRecipients: unexpected status %d", resp.StatusCode())
	}
	return awardsFromGen(resp.JSON200), nil
}

func awardsFromGen(r *gen.AwardsResponse) *Awards {
	out := &Awards{}
	if r == nil || r.Awards == nil {
		return out
	}
	out.Recipients = make([]AwardRecipient, 0, len(*r.Awards))
	for _, a := range *r.Awards {
		out.Recipients = append(out.Recipients, awardRecipientFromGen(a))
	}
	return out
}

func awardRecipientFromGen(a gen.AwardRecipient) AwardRecipient {
	out := AwardRecipient{}
	if a.Id != nil {
		out.ID = *a.Id
	}
	if a.Name != nil {
		out.Name = *a.Name
	}
	if a.Date != nil {
		out.Date = *a.Date
	}
	if a.Season != nil {
		out.Season = *a.Season
	}
	if a.Team != nil {
		out.Team = refFromGen(a.Team)
	}
	if a.Player != nil {
		out.Player = personFromGen(*a.Player)
	}
	if a.Notes != nil {
		out.Notes = *a.Notes
	}
	return out
}

func personFromGen(p gen.Person) Person {
	out := Person{}
	if p.Id != nil {
		out.ID = *p.Id
	}
	if p.FullName != nil {
		out.FullName = *p.FullName
	}
	if p.NameFirstLast != nil {
		out.NameFirstLast = *p.NameFirstLast
	}
	if p.Link != nil {
		out.Link = *p.Link
	}
	if p.PrimaryPosition != nil {
		out.PrimaryPosition = primaryPositionFromGen(*p.PrimaryPosition)
	}
	return out
}

func primaryPositionFromGen(p gen.PrimaryPosition) PrimaryPosition {
	out := PrimaryPosition{}
	if p.Code != nil {
		out.Code = *p.Code
	}
	if p.Name != nil {
		out.Name = *p.Name
	}
	if p.Type != nil {
		out.Type = *p.Type
	}
	if p.Abbreviation != nil {
		out.Abbreviation = *p.Abbreviation
	}
	return out
}
