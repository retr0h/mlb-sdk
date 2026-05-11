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

import "github.com/retr0h/mlb-sdk/internal/gen"

// Boxscore is the per-game team-stats view served by the MLB Stats API,
// flattened to two sides and re-presented through methods that hide the
// API's awkward bits. The raw response is retained on the value so future
// helpers can pull from any field without changing the public type.
type Boxscore struct {
	Home *BoxscoreTeam
	Away *BoxscoreTeam

	raw *gen.BoxscoreResponse
}

// Team returns the BoxscoreTeam for a given team ID, or nil if neither side
// of this boxscore matches.
func (b *Boxscore) Team(id TeamID) *BoxscoreTeam {
	if b == nil {
		return nil
	}
	if b.Home != nil && b.Home.ID == id {
		return b.Home
	}
	if b.Away != nil && b.Away.ID == id {
		return b.Away
	}
	return nil
}

// BoxscoreTeam is one side of a boxscore. Every meaningful field from the
// API is exposed as a public Go field; DoublePlaysTurned() is an additive
// helper because the API hides team-level double-plays in a free-text info
// block rather than a structured field.
type BoxscoreTeam struct {
	ID       TeamID
	Name     string
	Pitching PitchingStats
	Batting  BattingStats

	// raw is retained ONLY so DoublePlaysTurned() can walk the Info block.
	// Field promotion is the default; raw is the escape hatch.
	raw *gen.BoxscoreSide
}

// PitchingStats is the team-level pitching line for a single game — i.e.
// what the team's pitchers did. Field names follow Go conventions; the
// underlying API uses camelCase (strikeOuts, baseOnBalls).
type PitchingStats struct {
	Strikeouts int
	Hits       int // hits allowed
	Runs       int // runs allowed
	HomeRuns   int // HRs allowed
	Walks      int // mapped from API's baseOnBalls
}

// BattingStats is the team-level offensive line for a single game.
type BattingStats struct {
	Runs                 int
	Hits                 int
	HomeRuns             int
	RBI                  int
	StolenBases          int
	GroundIntoDoublePlay int // batter-side GIDP (not double plays the team turned)
}
