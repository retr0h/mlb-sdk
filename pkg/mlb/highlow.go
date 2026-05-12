// Copyright (c) 2026 John Dewey
//
// SPDX-License-Identifier: MIT

package mlb

import (
	"context"
	"fmt"

	"github.com/retr0h/mlb-sdk/internal/gen"
)

// HighLow fetches season high/low records. orgType is one of: player, team,
// division, league, sport, types. q.SortStat and q.Season must both be set
// (toddrob99: required_params=[["sortStat", "season"]]).
//
// Example:
//
//	hl, _ := c.HighLow(ctx, "player", mlb.HighLowQuery{
//	    SortStat: "homeRuns", Season: 2024, SportIDs: "1", Limit: 5,
//	})
//	for _, g := range hl.Results {
//	    for _, s := range g.Splits {
//	        fmt.Println(s.Player.FullName, s.Stat)
//	    }
//	}
func (c *Client) HighLow(ctx context.Context, orgType string, q HighLowQuery) (*HighLow, error) {
	if q.SortStat == "" || q.Season == 0 {
		return nil, fmt.Errorf(
			"mlb: highLow: %w: SortStat and Season are both required",
			ErrInvalidQuery,
		)
	}

	params := &gen.GetHighLowParams{}
	params.SortStat = ptr(q.SortStat)
	params.Season = ptr(q.Season)
	if q.GameType != "" {
		params.GameType = ptr(q.GameType)
	}
	if q.TeamID != 0 {
		params.TeamId = ptr(q.TeamID)
	}
	if q.LeagueID != 0 {
		params.LeagueId = ptr(q.LeagueID)
	}
	if q.SportIDs != "" {
		params.SportIds = ptr(q.SportIDs)
	}
	if q.StatGroup != "" {
		params.StatGroup = ptr(q.StatGroup)
	}
	if q.Limit != 0 {
		params.Limit = ptr(q.Limit)
	}
	if q.Fields != "" {
		params.Fields = ptr(q.Fields)
	}

	resp, err := c.raw.GetHighLowWithResponse(ctx, orgType, params)
	if err != nil {
		return nil, fmt.Errorf("mlb: highLow: %w", err)
	}
	if resp.StatusCode() == 404 {
		return nil, ErrNotFound
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("mlb: highLow: unexpected status %d", resp.StatusCode())
	}
	return highLowFromGen(resp.JSON200), nil
}

func highLowFromGen(r *gen.HighLowResponse) *HighLow {
	out := &HighLow{}
	if r == nil || r.HighLowResults == nil {
		return out
	}
	out.Results = make([]HighLowGroup, 0, len(*r.HighLowResults))
	for _, g := range *r.HighLowResults {
		out.Results = append(out.Results, highLowGroupFromGen(g))
	}
	return out
}

func highLowGroupFromGen(g gen.HighLowGroup) HighLowGroup {
	out := HighLowGroup{}
	if g.Group != nil && g.Group.DisplayName != nil {
		out.Group = *g.Group.DisplayName
	}
	if g.TotalSplits != nil {
		out.TotalSplits = *g.TotalSplits
	}
	if g.Splits != nil {
		out.Splits = make([]HighLowSplit, 0, len(*g.Splits))
		for _, s := range *g.Splits {
			out.Splits = append(out.Splits, highLowSplitFromGen(s))
		}
	}
	return out
}

func highLowSplitFromGen(s gen.HighLowSplit) HighLowSplit {
	out := HighLowSplit{}
	if s.Season != nil {
		out.Season = *s.Season
	}
	if s.Stat != nil {
		out.Stat = *s.Stat
	}
	if s.Team != nil {
		if s.Team.Id != nil {
			out.Team.ID = TeamID(*s.Team.Id)
		}
		if s.Team.Name != nil {
			out.Team.Name = *s.Team.Name
		}
	}
	if s.Player != nil {
		out.Player = personFromGen(*s.Player)
	}
	return out
}
