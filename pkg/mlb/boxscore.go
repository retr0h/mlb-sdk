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
	"fmt"
	"strconv"
	"strings"

	"github.com/retr0h/mlb-sdk/internal/gen"
)

// Boxscore fetches the boxscore for a game by its MLB-assigned gamePk.
func (c *Client) Boxscore(ctx context.Context, gamePk int) (*Boxscore, error) {
	resp, err := c.raw.GetBoxscoreWithResponse(ctx, gamePk)
	if err != nil {
		return nil, fmt.Errorf("mlb: boxscore: %w", err)
	}
	if resp.StatusCode() == 404 {
		return nil, ErrNotFound
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("mlb: boxscore: unexpected status %d", resp.StatusCode())
	}
	return boxscoreFromGen(resp.JSON200), nil
}

func boxscoreFromGen(r *gen.BoxscoreResponse) *Boxscore {
	out := &Boxscore{raw: r}
	if r == nil || r.Teams == nil {
		return out
	}
	out.Home = teamFromGen(r.Teams.Home)
	out.Away = teamFromGen(r.Teams.Away)
	return out
}

func teamFromGen(s *gen.BoxscoreSide) *BoxscoreTeam {
	if s == nil {
		return nil
	}
	out := &BoxscoreTeam{raw: s}
	if s.Team != nil {
		if s.Team.Id != nil {
			out.ID = TeamID(*s.Team.Id)
		}
		if s.Team.Name != nil {
			out.Name = *s.Team.Name
		}
	}
	return out
}

// DoublePlaysTurned returns the number of double plays this team turned in
// the game. The MLB Stats API does NOT expose this in the structured
// teamStats.fielding block; it appears only in the boxscore's free-text
// info section under FIELDING / DP, with values like "(Smith-Jones)." for
// one DP or "2 (...)." for two. Returns 0 when no DP entry is present.
func (t *BoxscoreTeam) DoublePlaysTurned() int {
	if t == nil || t.raw == nil || t.raw.Info == nil {
		return 0
	}
	for _, section := range *t.raw.Info {
		if section.Title == nil || *section.Title != "FIELDING" {
			continue
		}
		if section.FieldList == nil {
			continue
		}
		for _, item := range *section.FieldList {
			if item.Label == nil || *item.Label != "DP" {
				continue
			}
			if item.Value == nil {
				return 0
			}
			return parseDPCount(*item.Value)
		}
	}
	return 0
}

// parseDPCount reads the leading integer from a DP value. Single DPs are
// recorded without a leading number ("(Smith-Jones)."); the count then
// defaults to 1. Empty input returns 0 — caller treats this as "no DP entry".
func parseDPCount(value string) int {
	value = strings.TrimSpace(value)
	if value == "" {
		return 0
	}
	end := 0
	for end < len(value) && value[end] >= '0' && value[end] <= '9' {
		end++
	}
	if end > 0 {
		if n, err := strconv.Atoi(value[:end]); err == nil {
			return n
		}
	}
	return 1
}
