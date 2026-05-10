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

// BoxscoreTeam is one side of a boxscore. Methods on this type read from the
// generated response shape without exposing it.
type BoxscoreTeam struct {
	ID   TeamID
	Name string

	raw *gen.BoxscoreSide
}
