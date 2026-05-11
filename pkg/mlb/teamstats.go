// Copyright (c) 2026 John Dewey

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to
// deal in the Software without restriction, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
// sell copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER
// DEALINGS IN THE SOFTWARE.

package mlb

import (
	"context"
	"errors"
	"fmt"

	"github.com/retr0h/mlb-sdk/internal/gen"
)

// ErrInvalidQuery is returned when a query struct is missing a field the SDK
// requires up-front (e.g. TeamStatsQuery.Team). Callers can `errors.Is` to
// distinguish this from server-side errors.
var ErrInvalidQuery = errors.New("mlb: invalid query")

// TeamStats fetches aggregated team stats. q.Team is required; supplying
// q.Group and q.Type narrows the response to a single group×window cell.
//
// Example:
//
//	ts, _ := c.TeamStats(ctx, mlb.TeamStatsQuery{
//	    Team:   mlb.LAD,
//	    Season: 2026,
//	    Type:   mlb.TeamStatTypeSeason,
//	    Group:  mlb.TeamStatGroupFielding,
//	})
//	dp := ts.Group(mlb.TeamStatGroupFielding).Season("2026").DoublePlays()
func (c *Client) TeamStats(ctx context.Context, q TeamStatsQuery) (*TeamStats, error) {
	if q.Team == 0 {
		return nil, fmt.Errorf("mlb: teamStats: %w: Team is required", ErrInvalidQuery)
	}
	params := &gen.GetTeamStatsParams{}
	if q.Season != 0 {
		params.Season = ptr(q.Season)
	}
	if q.Type != "" {
		s := string(q.Type)
		params.Stats = &s
	}
	if q.Group != "" {
		g := string(q.Group)
		params.Group = &g
	}

	resp, err := c.raw.GetTeamStatsWithResponse(ctx, q.Team.Int(), params)
	if err != nil {
		return nil, fmt.Errorf("mlb: teamStats: %w", err)
	}
	if resp.StatusCode() == 404 {
		return nil, ErrNotFound
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("mlb: teamStats: unexpected status %d", resp.StatusCode())
	}
	return teamStatsFromGen(resp.JSON200), nil
}

func teamStatsFromGen(r *gen.TeamStatsResponse) *TeamStats {
	out := &TeamStats{raw: r}
	if r == nil || r.Stats == nil {
		return out
	}
	out.Groups = make([]TeamStatGroupResult, 0, len(*r.Stats))
	for _, g := range *r.Stats {
		out.Groups = append(out.Groups, statGroupFromGen(g))
	}
	return out
}

func statGroupFromGen(g gen.StatGroup) TeamStatGroupResult {
	out := TeamStatGroupResult{}
	if g.Type != nil && g.Type.DisplayName != nil {
		out.Type = *g.Type.DisplayName
	}
	if g.Group != nil && g.Group.DisplayName != nil {
		out.Group = *g.Group.DisplayName
	}
	if g.Splits != nil {
		out.Splits = make([]TeamStatsSplit, 0, len(*g.Splits))
		for _, s := range *g.Splits {
			out.Splits = append(out.Splits, splitFromGen(s))
		}
	}
	return out
}

func splitFromGen(s gen.StatSplit) TeamStatsSplit {
	out := TeamStatsSplit{Stat: map[string]any{}}
	if s.Season != nil {
		out.Season = *s.Season
	}
	if s.Stat != nil {
		out.Stat = *s.Stat
	}
	return out
}
