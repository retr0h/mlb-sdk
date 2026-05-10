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

// dateFmt is the YYYY-MM-DD format the MLB schedule API expects.
const dateFmt = "2006-01-02"

// sportIDMLB selects MLB regular-season games. The schedule endpoint serves
// every league MLB tracks (minors, college, etc.); we always pin to MLB.
const sportIDMLB = 1

// Schedule returns games matching q. SportId is always pinned to 1 (MLB);
// callers cannot request other leagues through this method.
func (c *Client) Schedule(ctx context.Context, q ScheduleQuery) ([]Game, error) {
	params := &gen.GetScheduleParams{
		SportId: ptr(sportIDMLB),
	}
	if q.Team != 0 {
		params.TeamId = ptr(q.Team.Int())
	}
	if !q.On.IsZero() {
		d := q.On.Format(dateFmt)
		params.Date = &d
	}
	if !q.From.IsZero() && !q.To.IsZero() {
		from := q.From.Format(dateFmt)
		to := q.To.Format(dateFmt)
		params.StartDate = &from
		params.EndDate = &to
	}

	resp, err := c.raw.GetScheduleWithResponse(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("mlb: schedule: %w", err)
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("mlb: schedule: unexpected status %d", resp.StatusCode())
	}
	if resp.JSON200.Dates == nil {
		return []Game{}, nil
	}

	out := make([]Game, 0)
	for _, d := range *resp.JSON200.Dates {
		if d.Games == nil {
			continue
		}
		for _, g := range *d.Games {
			out = append(out, gameFromGen(g))
		}
	}
	return out, nil
}

func gameFromGen(g gen.ScheduleGame) Game {
	out := Game{}
	if g.GamePk != nil {
		out.GamePk = *g.GamePk
	}
	if g.GameDate != nil {
		out.Date = *g.GameDate
	}
	if g.Status != nil && g.Status.AbstractGameState != nil {
		out.Status = GameStatus(*g.Status.AbstractGameState)
	}
	if g.Teams != nil {
		out.Home = teamScoreFromGen(g.Teams.Home)
		out.Away = teamScoreFromGen(g.Teams.Away)
	}
	return out
}

func teamScoreFromGen(s *gen.SideScoreboard) TeamScore {
	if s == nil {
		return TeamScore{}
	}
	out := TeamScore{}
	if s.Team != nil {
		if s.Team.Id != nil {
			out.ID = TeamID(*s.Team.Id)
		}
		if s.Team.Name != nil {
			out.Name = *s.Team.Name
		}
	}
	if s.Score != nil {
		out.Score = *s.Score
	}
	return out
}

// ptr is a tiny helper for taking the address of a literal — needed because
// every gen request param is *T.
func ptr[T any](v T) *T { return &v }

// ErrNotFound is returned when the API replies 404 for a resource lookup.
var ErrNotFound = errors.New("mlb: not found")
