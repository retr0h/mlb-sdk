// Copyright (c) 2026 John Dewey
//
// SPDX-License-Identifier: MIT

package mlb

import (
	"context"
	"fmt"

	"github.com/retr0h/mlb-sdk/internal/gen"
)

// Team fetches a single team by id. The underlying endpoint responds with
// `{"teams": [<one entry>]}`; this method collapses the wrapper and
// returns the single TeamInfo. An empty array maps to ErrNotFound.
//
// Sub-objects (League, Division, Sport, SpringLeague, Venue) only carry
// id/name/link unless the request hydrates them via q.Hydrate (e.g.
// "league,division,sport,springLeague,venue").
//
// Example:
//
//	t, _ := c.Team(ctx, 119, mlb.TeamQuery{
//	    Hydrate: "league,division,sport,springLeague,venue",
//	})
//	fmt.Println(t.Name, "→", t.Division.Name, "/", t.League.Abbreviation)
func (c *Client) Team(ctx context.Context, teamID int, q TeamQuery) (*TeamInfo, error) {
	params := &gen.GetTeamParams{}
	if q.Season != 0 {
		params.Season = ptr(q.Season)
	}
	if q.SportID != 0 {
		params.SportId = ptr(q.SportID)
	}
	if q.Hydrate != "" {
		params.Hydrate = ptr(q.Hydrate)
	}
	if q.Fields != "" {
		params.Fields = ptr(q.Fields)
	}

	resp, err := c.raw.GetTeamWithResponse(ctx, teamID, params)
	if err != nil {
		return nil, fmt.Errorf("mlb: team: %w", err)
	}
	if resp.StatusCode() == 404 {
		return nil, ErrNotFound
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("mlb: team: unexpected status %d", resp.StatusCode())
	}
	t := teamFromResponse(resp.JSON200)
	if t == nil {
		return nil, ErrNotFound
	}
	return t, nil
}

// teamFromResponse collapses the teams array to a single TeamInfo. Empty
// array → nil so the caller maps to ErrNotFound.
func teamFromResponse(r *gen.TeamsResponse) *TeamInfo {
	if r == nil || r.Teams == nil || len(*r.Teams) == 0 {
		return nil
	}
	t := teamInfoFromGen((*r.Teams)[0])
	return &t
}

func teamInfoFromGen(t gen.TeamInfo) TeamInfo {
	out := TeamInfo{}
	if t.Id != nil {
		out.ID = *t.Id
	}
	if t.Name != nil {
		out.Name = *t.Name
	}
	if t.Link != nil {
		out.Link = *t.Link
	}
	if t.Season != nil {
		out.Season = *t.Season
	}
	if t.Venue != nil {
		out.Venue = venueFromGen(*t.Venue)
	}
	if t.SpringLeague != nil {
		out.SpringLeague = leagueInfoFromGen(*t.SpringLeague)
	}
	if t.SpringVenue != nil {
		out.SpringVenue = refFromGen(t.SpringVenue)
	}
	if t.TeamCode != nil {
		out.TeamCode = *t.TeamCode
	}
	if t.FileCode != nil {
		out.FileCode = *t.FileCode
	}
	if t.Abbreviation != nil {
		out.Abbreviation = *t.Abbreviation
	}
	if t.TeamName != nil {
		out.TeamName = *t.TeamName
	}
	if t.LocationName != nil {
		out.LocationName = *t.LocationName
	}
	if t.FirstYearOfPlay != nil {
		out.FirstYearOfPlay = *t.FirstYearOfPlay
	}
	if t.League != nil {
		out.League = leagueInfoFromGen(*t.League)
	}
	if t.Division != nil {
		out.Division = divisionFromGen(*t.Division)
	}
	if t.Sport != nil {
		out.Sport = sportFromGen(*t.Sport)
	}
	if t.ShortName != nil {
		out.ShortName = *t.ShortName
	}
	if t.FranchiseName != nil {
		out.FranchiseName = *t.FranchiseName
	}
	if t.ClubName != nil {
		out.ClubName = *t.ClubName
	}
	if t.Active != nil {
		out.Active = *t.Active
	}
	if t.AllStarStatus != nil {
		out.AllStarStatus = *t.AllStarStatus
	}
	return out
}

func leagueInfoFromGen(l gen.LeagueInfo) LeagueInfo {
	out := LeagueInfo{}
	if l.Id != nil {
		out.ID = *l.Id
	}
	if l.Name != nil {
		out.Name = *l.Name
	}
	if l.Link != nil {
		out.Link = *l.Link
	}
	if l.Abbreviation != nil {
		out.Abbreviation = *l.Abbreviation
	}
	if l.NameShort != nil {
		out.NameShort = *l.NameShort
	}
	if l.SeasonState != nil {
		out.SeasonState = *l.SeasonState
	}
	if l.HasWildCard != nil {
		out.HasWildCard = *l.HasWildCard
	}
	if l.HasSplitSeason != nil {
		out.HasSplitSeason = *l.HasSplitSeason
	}
	if l.NumGames != nil {
		out.NumGames = *l.NumGames
	}
	if l.HasPlayoffPoints != nil {
		out.HasPlayoffPoints = *l.HasPlayoffPoints
	}
	if l.NumTeams != nil {
		out.NumTeams = *l.NumTeams
	}
	if l.NumWildcardTeams != nil {
		out.NumWildcardTeams = *l.NumWildcardTeams
	}
	if l.SeasonDateInfo != nil {
		out.SeasonDateInfo = seasonFromGen(*l.SeasonDateInfo)
	}
	if l.Season != nil {
		out.Season = *l.Season
	}
	if l.OrgCode != nil {
		out.OrgCode = *l.OrgCode
	}
	if l.ConferencesInUse != nil {
		out.ConferencesInUse = *l.ConferencesInUse
	}
	if l.DivisionsInUse != nil {
		out.DivisionsInUse = *l.DivisionsInUse
	}
	if l.Sport != nil {
		out.Sport = refFromGen(l.Sport)
	}
	if l.SortOrder != nil {
		out.SortOrder = *l.SortOrder
	}
	if l.Active != nil {
		out.Active = *l.Active
	}
	return out
}
