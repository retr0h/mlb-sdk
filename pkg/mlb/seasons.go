// Copyright (c) 2026 John Dewey
//
// SPDX-License-Identifier: MIT

package mlb

import (
	"context"
	"fmt"
	"time"

	"github.com/retr0h/mlb-sdk/internal/gen"
)

const seasonDateFmt = "2006-01-02"

// Seasons fetches season metadata. At least one of q.SportID, q.DivisionID
// or q.LeagueID must be set (the toddrob99 catalog encodes this as
// `required_params: [["sportId"], ["divisionId"], ["leagueId"]]`). Set q.All
// to call the `/api/v1/seasons/all` path and receive every season the MLB
// API tracks.
//
// Example:
//
//	s, _ := c.Seasons(ctx, mlb.SeasonsQuery{SportID: 1, Season: 2024})
//	cur := s.Season("2024")
//	fmt.Println(cur.RegularSeasonStartDate, "→", cur.RegularSeasonEndDate)
func (c *Client) Seasons(ctx context.Context, q SeasonsQuery) (*Seasons, error) {
	if q.SportID == 0 && q.DivisionID == 0 && q.LeagueID == 0 {
		return nil, fmt.Errorf(
			"mlb: seasons: %w: one of SportID, DivisionID, or LeagueID is required",
			ErrInvalidQuery,
		)
	}

	if q.All {
		return c.seasonsAll(ctx, q)
	}
	return c.seasonsFiltered(ctx, q)
}

func (c *Client) seasonsFiltered(ctx context.Context, q SeasonsQuery) (*Seasons, error) {
	params := &gen.GetSeasonsParams{}
	if q.Season != 0 {
		params.Season = ptr(q.Season)
	}
	if q.SportID != 0 {
		params.SportId = ptr(q.SportID)
	}
	if q.DivisionID != 0 {
		params.DivisionId = ptr(q.DivisionID)
	}
	if q.LeagueID != 0 {
		params.LeagueId = ptr(q.LeagueID)
	}
	if q.Fields != "" {
		params.Fields = ptr(q.Fields)
	}

	resp, err := c.raw.GetSeasonsWithResponse(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("mlb: seasons: %w", err)
	}
	if resp.StatusCode() == 404 {
		return nil, ErrNotFound
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("mlb: seasons: unexpected status %d", resp.StatusCode())
	}
	return seasonsFromGen(resp.JSON200), nil
}

func (c *Client) seasonsAll(ctx context.Context, q SeasonsQuery) (*Seasons, error) {
	params := &gen.GetAllSeasonsParams{}
	if q.SportID != 0 {
		params.SportId = ptr(q.SportID)
	}
	if q.DivisionID != 0 {
		params.DivisionId = ptr(q.DivisionID)
	}
	if q.LeagueID != 0 {
		params.LeagueId = ptr(q.LeagueID)
	}
	if q.Fields != "" {
		params.Fields = ptr(q.Fields)
	}

	resp, err := c.raw.GetAllSeasonsWithResponse(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("mlb: seasons: %w", err)
	}
	if resp.StatusCode() == 404 {
		return nil, ErrNotFound
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("mlb: seasons: unexpected status %d", resp.StatusCode())
	}
	return seasonsFromGen(resp.JSON200), nil
}

func seasonsFromGen(r *gen.SeasonsResponse) *Seasons {
	out := &Seasons{}
	if r == nil || r.Seasons == nil {
		return out
	}
	out.Seasons = make([]Season, 0, len(*r.Seasons))
	for _, s := range *r.Seasons {
		out.Seasons = append(out.Seasons, seasonFromGen(s))
	}
	return out
}

func seasonFromGen(s gen.Season) Season {
	out := Season{}
	if s.SeasonId != nil {
		out.SeasonID = *s.SeasonId
	}
	if s.HasWildcard != nil {
		out.HasWildcard = *s.HasWildcard
	}
	out.PreSeasonStartDate = parseSeasonDate(s.PreSeasonStartDate)
	out.PreSeasonEndDate = parseSeasonDate(s.PreSeasonEndDate)
	out.SeasonStartDate = parseSeasonDate(s.SeasonStartDate)
	out.SpringStartDate = parseSeasonDate(s.SpringStartDate)
	out.SpringEndDate = parseSeasonDate(s.SpringEndDate)
	out.RegularSeasonStartDate = parseSeasonDate(s.RegularSeasonStartDate)
	out.LastDate1stHalf = parseSeasonDate(s.LastDate1stHalf)
	out.AllStarDate = parseSeasonDate(s.AllStarDate)
	out.FirstDate2ndHalf = parseSeasonDate(s.FirstDate2ndHalf)
	out.RegularSeasonEndDate = parseSeasonDate(s.RegularSeasonEndDate)
	out.PostSeasonStartDate = parseSeasonDate(s.PostSeasonStartDate)
	out.PostSeasonEndDate = parseSeasonDate(s.PostSeasonEndDate)
	out.SeasonEndDate = parseSeasonDate(s.SeasonEndDate)
	out.OffseasonStartDate = parseSeasonDate(s.OffseasonStartDate)
	out.OffSeasonEndDate = parseSeasonDate(s.OffSeasonEndDate)
	if s.SeasonLevelGamedayType != nil {
		out.SeasonLevelGamedayType = *s.SeasonLevelGamedayType
	}
	if s.GameLevelGamedayType != nil {
		out.GameLevelGamedayType = *s.GameLevelGamedayType
	}
	if s.QualifierPlateAppearances != nil {
		out.QualifierPlateAppearances = *s.QualifierPlateAppearances
	}
	if s.QualifierOutsPitched != nil {
		out.QualifierOutsPitched = *s.QualifierOutsPitched
	}
	return out
}

// parseSeasonDate parses a YYYY-MM-DD pointer string into a time.Time. Nil
// or unparseable input yields the zero time (callers can `.IsZero()`).
func parseSeasonDate(p *string) time.Time {
	if p == nil || *p == "" {
		return time.Time{}
	}
	t, err := time.Parse(seasonDateFmt, *p)
	if err != nil {
		return time.Time{}
	}
	return t
}
