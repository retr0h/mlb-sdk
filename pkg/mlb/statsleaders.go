// Copyright (c) 2026 John Dewey
//
// SPDX-License-Identifier: MIT

package mlb

import (
	"context"
	"fmt"

	"github.com/retr0h/mlb-sdk/internal/gen"
)

// StatsLeaders fetches league stat leaders. q.LeaderCategories is required
// (toddrob99: required_params=[["leaderCategories"]]).
//
// Example:
//
//	sl, _ := c.StatsLeaders(ctx, mlb.StatsLeadersQuery{
//	    LeaderCategories: "homeRuns", Season: 2024, SportID: 1, Limit: 5,
//	})
//	for _, cat := range sl.LeagueLeaders {
//	    for _, l := range cat.Leaders {
//	        fmt.Printf("#%d %s %s\n", l.Rank, l.Player.FullName, l.Value)
//	    }
//	}
func (c *Client) StatsLeaders(
	ctx context.Context,
	q StatsLeadersQuery,
) (*StatsLeaders, error) {
	if q.LeaderCategories == "" {
		return nil, fmt.Errorf(
			"mlb: statsLeaders: %w: LeaderCategories is required",
			ErrInvalidQuery,
		)
	}

	params := &gen.GetStatsLeadersParams{LeaderCategories: q.LeaderCategories}
	if q.Season != 0 {
		params.Season = ptr(q.Season)
	}
	if q.SportID != 0 {
		params.SportId = ptr(q.SportID)
	}
	if q.LeagueID != 0 {
		params.LeagueId = ptr(q.LeagueID)
	}
	if q.StatGroup != "" {
		params.StatGroup = ptr(q.StatGroup)
	}
	if q.PlayerPool != "" {
		params.PlayerPool = ptr(q.PlayerPool)
	}
	if q.LeaderGameTypes != "" {
		params.LeaderGameTypes = ptr(q.LeaderGameTypes)
	}
	if q.StatType != "" {
		params.StatType = ptr(q.StatType)
	}
	if q.Hydrate != "" {
		params.Hydrate = ptr(q.Hydrate)
	}
	if q.Limit != 0 {
		params.Limit = ptr(q.Limit)
	}
	if q.Fields != "" {
		params.Fields = ptr(q.Fields)
	}

	resp, err := c.raw.GetStatsLeadersWithResponse(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("mlb: statsLeaders: %w", err)
	}
	if resp.StatusCode() == 404 {
		return nil, ErrNotFound
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("mlb: statsLeaders: unexpected status %d", resp.StatusCode())
	}
	return statsLeadersFromGen(resp.JSON200), nil
}

func statsLeadersFromGen(r *gen.StatsLeadersResponse) *StatsLeaders {
	out := &StatsLeaders{}
	if r == nil || r.LeagueLeaders == nil {
		return out
	}
	out.LeagueLeaders = make([]LeaderCategory, 0, len(*r.LeagueLeaders))
	for _, c := range *r.LeagueLeaders {
		out.LeagueLeaders = append(out.LeagueLeaders, leaderCategoryFromGen(c))
	}
	return out
}

func leaderCategoryFromGen(c gen.LeaderCategory) LeaderCategory {
	out := LeaderCategory{}
	if c.LeaderCategory != nil {
		out.LeaderCategory = *c.LeaderCategory
	}
	if c.Season != nil {
		out.Season = *c.Season
	}
	if c.GameType != nil && c.GameType.DisplayName != nil {
		out.GameType = *c.GameType.DisplayName
	}
	if c.StatGroup != nil {
		out.StatGroup = *c.StatGroup
	}
	if c.TotalSplits != nil {
		out.TotalSplits = *c.TotalSplits
	}
	if c.Leaders != nil {
		out.Leaders = make([]LeaderEntry, 0, len(*c.Leaders))
		for _, l := range *c.Leaders {
			out.Leaders = append(out.Leaders, leaderEntryFromGen(l))
		}
	}
	return out
}

func leaderEntryFromGen(l gen.LeaderEntry) LeaderEntry {
	out := LeaderEntry{}
	if l.Rank != nil {
		out.Rank = *l.Rank
	}
	if l.Value != nil {
		out.Value = *l.Value
	}
	if l.Team != nil {
		if l.Team.Id != nil {
			out.Team.ID = TeamID(*l.Team.Id)
		}
		if l.Team.Name != nil {
			out.Team.Name = *l.Team.Name
		}
	}
	if l.League != nil {
		out.League = refFromGen(l.League)
	}
	if l.Person != nil {
		out.Player = personFromGen(*l.Person)
	}
	if l.Sport != nil {
		out.Sport = refFromGen(l.Sport)
	}
	return out
}
