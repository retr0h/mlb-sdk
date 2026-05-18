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
// Linescore fetches the linescore for a single game. Returns a per-inning
// breakdown plus game-total runs/hits/errors and the current
// defense/offense situation.
//
// Example:
//
//	ls, _ := c.Linescore(ctx, 745455, mlb.LinescoreQuery{})
//	fmt.Println("score:", ls.Teams.Away.Runs, "-", ls.Teams.Home.Runs)
//	for _, inn := range ls.Innings {
//	    fmt.Printf("  %s: away=%d home=%d\n", inn.OrdinalNum, inn.Away.Runs, inn.Home.Runs)
//	}
func (c *Client) Linescore(
	ctx context.Context,
	gamePk int,
	q LinescoreQuery,
) (*Linescore, error) {
	params := &gen.GetLinescoreParams{}
	if q.Timecode != "" {
		params.Timecode = ptr(q.Timecode)
	}
	if q.Fields != "" {
		params.Fields = ptr(q.Fields)
	}
//
	resp, err := c.raw.GetLinescoreWithResponse(ctx, gamePk, params)
	if err != nil {
		return nil, fmt.Errorf("mlb: linescore: %w", err)
	}
	if resp.StatusCode() == 404 {
		return nil, ErrNotFound
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("mlb: linescore: unexpected status %d", resp.StatusCode())
	}
	return linescoreFromGen(resp.JSON200), nil
}
//
func linescoreFromGen(r *gen.LinescoreResponse) *Linescore {
	out := &Linescore{}
	if r == nil {
		return out
	}
	if r.CurrentInning != nil {
		out.CurrentInning = *r.CurrentInning
	}
	if r.CurrentInningOrdinal != nil {
		out.CurrentInningOrdinal = *r.CurrentInningOrdinal
	}
	if r.InningState != nil {
		out.InningState = *r.InningState
	}
	if r.InningHalf != nil {
		out.InningHalf = *r.InningHalf
	}
	if r.IsTopInning != nil {
		out.IsTopInning = *r.IsTopInning
	}
	if r.ScheduledInnings != nil {
		out.ScheduledInnings = *r.ScheduledInnings
	}
	if r.Innings != nil {
		out.Innings = make([]LinescoreInning, 0, len(*r.Innings))
		for _, inn := range *r.Innings {
			out.Innings = append(out.Innings, linescoreInningFromGen(inn))
		}
	}
	if r.Teams != nil {
		out.Teams = linescoreTeamsFromGen(*r.Teams)
	}
	if r.Defense != nil {
		out.Defense = linescoreDefenseFromGen(*r.Defense)
	}
	if r.Offense != nil {
		out.Offense = linescoreOffenseFromGen(*r.Offense)
	}
	if r.Balls != nil {
		out.Balls = *r.Balls
	}
	if r.Strikes != nil {
		out.Strikes = *r.Strikes
	}
	if r.Outs != nil {
		out.Outs = *r.Outs
	}
	return out
}
//
func linescoreInningFromGen(i gen.LinescoreInning) LinescoreInning {
	out := LinescoreInning{}
	if i.Num != nil {
		out.Num = *i.Num
	}
	if i.OrdinalNum != nil {
		out.OrdinalNum = *i.OrdinalNum
	}
	if i.Home != nil {
		out.Home = linescoreInningHalfFromGen(*i.Home)
	}
	if i.Away != nil {
		out.Away = linescoreInningHalfFromGen(*i.Away)
	}
	return out
}
//
func linescoreInningHalfFromGen(h gen.LinescoreInningHalf) LinescoreInningHalf {
	out := LinescoreInningHalf{}
	if h.Runs != nil {
		out.Runs = *h.Runs
	}
	if h.Hits != nil {
		out.Hits = *h.Hits
	}
	if h.Errors != nil {
		out.Errors = *h.Errors
	}
	if h.LeftOnBase != nil {
		out.LeftOnBase = *h.LeftOnBase
	}
	return out
}
//
func linescoreTeamsFromGen(t gen.LinescoreTeams) LinescoreTeams {
	out := LinescoreTeams{}
	if t.Home != nil {
		out.Home = linescoreTeamTotalsFromGen(*t.Home)
	}
	if t.Away != nil {
		out.Away = linescoreTeamTotalsFromGen(*t.Away)
	}
	return out
}
//
func linescoreTeamTotalsFromGen(t gen.LinescoreTeamTotals) LinescoreTeamTotals {
	out := LinescoreTeamTotals{}
	if t.Runs != nil {
		out.Runs = *t.Runs
	}
	if t.Hits != nil {
		out.Hits = *t.Hits
	}
	if t.Errors != nil {
		out.Errors = *t.Errors
	}
	if t.LeftOnBase != nil {
		out.LeftOnBase = *t.LeftOnBase
	}
	if t.IsWinner != nil {
		out.IsWinner = *t.IsWinner
	}
	return out
}
//
func linescoreDefenseFromGen(d gen.LinescoreDefense) LinescoreDefense {
	out := LinescoreDefense{}
	if d.Pitcher != nil {
		out.Pitcher = personFromGen(*d.Pitcher)
	}
	if d.Catcher != nil {
		out.Catcher = personFromGen(*d.Catcher)
	}
	if d.First != nil {
		out.First = personFromGen(*d.First)
	}
	if d.Second != nil {
		out.Second = personFromGen(*d.Second)
	}
	if d.Third != nil {
		out.Third = personFromGen(*d.Third)
	}
	if d.Shortstop != nil {
		out.Shortstop = personFromGen(*d.Shortstop)
	}
	if d.Left != nil {
		out.Left = personFromGen(*d.Left)
	}
	if d.Center != nil {
		out.Center = personFromGen(*d.Center)
	}
	if d.Right != nil {
		out.Right = personFromGen(*d.Right)
	}
	return out
}
//
func linescoreOffenseFromGen(o gen.LinescoreOffense) LinescoreOffense {
	out := LinescoreOffense{}
	if o.Batter != nil {
		out.Batter = personFromGen(*o.Batter)
	}
	if o.OnDeck != nil {
		out.OnDeck = personFromGen(*o.OnDeck)
	}
	if o.InHole != nil {
		out.InHole = personFromGen(*o.InHole)
	}
	if o.First != nil {
		out.First = personFromGen(*o.First)
	}
	if o.Second != nil {
		out.Second = personFromGen(*o.Second)
	}
	if o.Third != nil {
		out.Third = personFromGen(*o.Third)
	}
	return out
}
