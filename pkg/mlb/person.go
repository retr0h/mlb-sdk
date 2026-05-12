// Copyright (c) 2026 John Dewey
//
// SPDX-License-Identifier: MIT

package mlb

import (
	"context"
	"fmt"

	"github.com/retr0h/mlb-sdk/internal/gen"
)

// Person fetches a single person by id. The underlying endpoint responds
// with `{"people": [<one entry>]}`; this method collapses the wrapper and
// maps an empty array to ErrNotFound.
//
// Example:
//
//	p, _ := c.Person(ctx, 660271, mlb.PersonQuery{})
//	fmt.Println(p.FullName, p.PrimaryPosition.Abbreviation)
func (c *Client) Person(
	ctx context.Context,
	personID int,
	q PersonQuery,
) (*PersonDetail, error) {
	params := &gen.GetPersonParams{}
	if q.Hydrate != "" {
		params.Hydrate = ptr(q.Hydrate)
	}
	if q.Fields != "" {
		params.Fields = ptr(q.Fields)
	}

	resp, err := c.raw.GetPersonWithResponse(ctx, personID, params)
	if err != nil {
		return nil, fmt.Errorf("mlb: person: %w", err)
	}
	if resp.StatusCode() == 404 {
		return nil, ErrNotFound
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("mlb: person: unexpected status %d", resp.StatusCode())
	}
	p := personDetailFromResponse(resp.JSON200)
	if p == nil {
		return nil, ErrNotFound
	}
	return p, nil
}

// People fetches multiple people by comma-separated ids. q.PersonIDs is
// required.
//
// Example:
//
//	pp, _ := c.People(ctx, mlb.PeopleQuery{PersonIDs: "660271,545361"})
//	for _, p := range pp {
//	    fmt.Println(p.FullName)
//	}
func (c *Client) People(
	ctx context.Context,
	q PeopleQuery,
) ([]PersonDetail, error) {
	if q.PersonIDs == "" {
		return nil, fmt.Errorf(
			"mlb: people: %w: PersonIDs is required",
			ErrInvalidQuery,
		)
	}
	params := &gen.GetPeopleParams{PersonIds: q.PersonIDs}
	if q.Hydrate != "" {
		params.Hydrate = ptr(q.Hydrate)
	}
	if q.Fields != "" {
		params.Fields = ptr(q.Fields)
	}

	resp, err := c.raw.GetPeopleWithResponse(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("mlb: people: %w", err)
	}
	if resp.StatusCode() == 404 {
		return nil, ErrNotFound
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("mlb: people: unexpected status %d", resp.StatusCode())
	}
	return peopleFromGen(resp.JSON200), nil
}

func personDetailFromResponse(r *gen.PeopleResponse) *PersonDetail {
	if r == nil || r.People == nil || len(*r.People) == 0 {
		return nil
	}
	p := personDetailFromGen((*r.People)[0])
	return &p
}

func peopleFromGen(r *gen.PeopleResponse) []PersonDetail {
	if r == nil || r.People == nil {
		return nil
	}
	out := make([]PersonDetail, 0, len(*r.People))
	for _, p := range *r.People {
		out = append(out, personDetailFromGen(p))
	}
	return out
}

func personDetailFromGen(p gen.PersonDetail) PersonDetail {
	out := PersonDetail{}
	if p.Id != nil {
		out.ID = *p.Id
	}
	if p.FullName != nil {
		out.FullName = *p.FullName
	}
	if p.FirstName != nil {
		out.FirstName = *p.FirstName
	}
	if p.LastName != nil {
		out.LastName = *p.LastName
	}
	if p.Link != nil {
		out.Link = *p.Link
	}
	if p.PrimaryNumber != nil {
		out.PrimaryNumber = *p.PrimaryNumber
	}
	if p.BirthDate != nil {
		out.BirthDate = *p.BirthDate
	}
	if p.CurrentAge != nil {
		out.CurrentAge = *p.CurrentAge
	}
	if p.BirthCity != nil {
		out.BirthCity = *p.BirthCity
	}
	if p.BirthCountry != nil {
		out.BirthCountry = *p.BirthCountry
	}
	if p.Height != nil {
		out.Height = *p.Height
	}
	if p.Weight != nil {
		out.Weight = *p.Weight
	}
	if p.Active != nil {
		out.Active = *p.Active
	}
	if p.PrimaryPosition != nil {
		out.PrimaryPosition = primaryPositionFromGen(*p.PrimaryPosition)
	}
	if p.UseName != nil {
		out.UseName = *p.UseName
	}
	if p.UseLastName != nil {
		out.UseLastName = *p.UseLastName
	}
	if p.BoxscoreName != nil {
		out.BoxscoreName = *p.BoxscoreName
	}
	if p.NickName != nil {
		out.NickName = *p.NickName
	}
	if p.Gender != nil {
		out.Gender = *p.Gender
	}
	if p.IsPlayer != nil {
		out.IsPlayer = *p.IsPlayer
	}
	if p.IsVerified != nil {
		out.IsVerified = *p.IsVerified
	}
	if p.Pronunciation != nil {
		out.Pronunciation = *p.Pronunciation
	}
	if p.MlbDebutDate != nil {
		out.MlbDebutDate = *p.MlbDebutDate
	}
	if p.BatSide != nil {
		out.BatSide = handSideFromGen(*p.BatSide)
	}
	if p.PitchHand != nil {
		out.PitchHand = handSideFromGen(*p.PitchHand)
	}
	if p.NameFirstLast != nil {
		out.NameFirstLast = *p.NameFirstLast
	}
	if p.NameSlug != nil {
		out.NameSlug = *p.NameSlug
	}
	if p.FirstLastName != nil {
		out.FirstLastName = *p.FirstLastName
	}
	if p.LastFirstName != nil {
		out.LastFirstName = *p.LastFirstName
	}
	if p.LastInitName != nil {
		out.LastInitName = *p.LastInitName
	}
	if p.InitLastName != nil {
		out.InitLastName = *p.InitLastName
	}
	if p.FullFMLName != nil {
		out.FullFMLName = *p.FullFMLName
	}
	if p.FullLFMName != nil {
		out.FullLFMName = *p.FullLFMName
	}
	if p.StrikeZoneTop != nil {
		out.StrikeZoneTop = *p.StrikeZoneTop
	}
	if p.StrikeZoneBottom != nil {
		out.StrikeZoneBottom = *p.StrikeZoneBottom
	}
	return out
}

func handSideFromGen(h gen.HandSide) HandSide {
	out := HandSide{}
	if h.Code != nil {
		out.Code = *h.Code
	}
	if h.Description != nil {
		out.Description = *h.Description
	}
	return out
}
