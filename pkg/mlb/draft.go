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
// Draft fetches draft data for a given year.
//
// Example:
//
//	d, _ := c.Draft(ctx, 2024, mlb.DraftQuery{})
//	for _, r := range d.Rounds {
//	    for _, p := range r.Picks {
//	        fmt.Printf("#%d %s — %s\n", p.PickNumber, p.Person.FullName, p.Team.Name)
//	    }
//	}
func (c *Client) Draft(ctx context.Context, year int, q DraftQuery) (*DraftData, error) {
	params := &gen.GetDraftParams{}
	if q.Round != "" {
		params.Round = ptr(q.Round)
	}
	if q.Fields != "" {
		params.Fields = ptr(q.Fields)
	}
//
	resp, err := c.raw.GetDraftWithResponse(ctx, year, params)
	if err != nil {
		return nil, fmt.Errorf("mlb: draft: %w", err)
	}
	if resp.StatusCode() == 404 {
		return nil, ErrNotFound
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("mlb: draft: unexpected status %d", resp.StatusCode())
	}
	return draftFromGen(resp.JSON200), nil
}
//
func draftFromGen(r *gen.DraftResponse) *DraftData {
	out := &DraftData{}
	if r == nil || r.Drafts == nil {
		return out
	}
	d := r.Drafts
	if d.DraftYear != nil {
		out.DraftYear = *d.DraftYear
	}
	if d.Rounds != nil {
		out.Rounds = make([]DraftRound, 0, len(*d.Rounds))
		for _, rd := range *d.Rounds {
			out.Rounds = append(out.Rounds, draftRoundFromGen(rd))
		}
	}
	return out
}
//
func draftRoundFromGen(r gen.DraftRound) DraftRound {
	out := DraftRound{}
	if r.Round != nil {
		out.Round = *r.Round
	}
	if r.Picks != nil {
		out.Picks = make([]DraftPick, 0, len(*r.Picks))
		for _, p := range *r.Picks {
			out.Picks = append(out.Picks, draftPickFromGen(p))
		}
	}
	return out
}
//
func draftPickFromGen(p gen.DraftPick) DraftPick {
	out := DraftPick{}
	if p.BisPlayerId != nil {
		out.BisPlayerID = *p.BisPlayerId
	}
	if p.PickRound != nil {
		out.PickRound = *p.PickRound
	}
	if p.PickNumber != nil {
		out.PickNumber = *p.PickNumber
	}
	if p.DisplayPickNumber != nil {
		out.DisplayPickNumber = *p.DisplayPickNumber
	}
	if p.RoundPickNumber != nil {
		out.RoundPickNumber = *p.RoundPickNumber
	}
	if p.Rank != nil {
		out.Rank = *p.Rank
	}
	if p.PickValue != nil {
		out.PickValue = *p.PickValue
	}
	if p.SigningBonus != nil {
		out.SigningBonus = *p.SigningBonus
	}
	if p.Home != nil {
		out.Home = draftHomeFromGen(*p.Home)
	}
	if p.ScoutingReport != nil {
		out.ScoutingReport = *p.ScoutingReport
	}
	if p.School != nil {
		out.School = draftSchoolFromGen(*p.School)
	}
	if p.Blurb != nil {
		out.Blurb = *p.Blurb
	}
	if p.HeadshotLink != nil {
		out.HeadshotLink = *p.HeadshotLink
	}
	if p.Person != nil {
		out.Person = personDetailFromGen(*p.Person)
	}
	if p.Team != nil {
		out.Team = teamInfoFromGen(*p.Team)
	}
	if p.DraftType != nil {
		out.DraftType = draftTypeRefFromGen(*p.DraftType)
	}
	if p.IsDrafted != nil {
		out.IsDrafted = *p.IsDrafted
	}
	if p.IsPass != nil {
		out.IsPass = *p.IsPass
	}
	if p.Year != nil {
		out.Year = *p.Year
	}
	return out
}
//
func draftHomeFromGen(h gen.DraftHome) DraftHome {
	out := DraftHome{}
	if h.City != nil {
		out.City = *h.City
	}
	if h.Country != nil {
		out.Country = *h.Country
	}
	return out
}
//
func draftSchoolFromGen(s gen.DraftSchool) DraftSchool {
	out := DraftSchool{}
	if s.Name != nil {
		out.Name = *s.Name
	}
	if s.SchoolClass != nil {
		out.SchoolClass = *s.SchoolClass
	}
	if s.City != nil {
		out.City = *s.City
	}
	if s.Country != nil {
		out.Country = *s.Country
	}
	if s.State != nil {
		out.State = *s.State
	}
	return out
}
//
func draftTypeRefFromGen(d gen.DraftTypeRef) DraftTypeRef {
	out := DraftTypeRef{}
	if d.Code != nil {
		out.Code = *d.Code
	}
	if d.Description != nil {
		out.Description = *d.Description
	}
	return out
}
