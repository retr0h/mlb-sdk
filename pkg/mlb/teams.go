// Copyright (c) 2026 John Dewey
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to
// deal in the Software without restriction, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
// sell copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER
// DEALINGS IN THE SOFTWARE.

package mlb

// TeamID is the MLB Stats API numeric team identifier. Constants below give
// each franchise a typed name; the underlying value is the integer the API
// expects in `teamId` query parameters.
type TeamID int

// MLB teams (regular season; expansion / relocation TBD as needed).
const (
	LAD TeamID = 119 // Los Angeles Dodgers
	SD  TeamID = 135 // San Diego Padres
	SF  TeamID = 137 // San Francisco Giants
	ARI TeamID = 109 // Arizona Diamondbacks
	COL TeamID = 115 // Colorado Rockies
	LAA TeamID = 108 // Los Angeles Angels
	OAK TeamID = 133 // Oakland Athletics
	SEA TeamID = 136 // Seattle Mariners
	TEX TeamID = 140 // Texas Rangers
	HOU TeamID = 117 // Houston Astros
	NYY TeamID = 147 // New York Yankees
	NYM TeamID = 121 // New York Mets
	BOS TeamID = 111 // Boston Red Sox
	CHC TeamID = 112 // Chicago Cubs
	CHW TeamID = 145 // Chicago White Sox
	ATL TeamID = 144 // Atlanta Braves
	MIA TeamID = 146 // Miami Marlins
	PHI TeamID = 143 // Philadelphia Phillies
	WSH TeamID = 120 // Washington Nationals
	BAL TeamID = 110 // Baltimore Orioles
	TB  TeamID = 139 // Tampa Bay Rays
	TOR TeamID = 141 // Toronto Blue Jays
	CLE TeamID = 114 // Cleveland Guardians
	DET TeamID = 116 // Detroit Tigers
	KC  TeamID = 118 // Kansas City Royals
	MIN TeamID = 142 // Minnesota Twins
	CIN TeamID = 113 // Cincinnati Reds
	MIL TeamID = 158 // Milwaukee Brewers
	PIT TeamID = 134 // Pittsburgh Pirates
	STL TeamID = 138 // St. Louis Cardinals
)

// Int returns the numeric form expected by the MLB Stats API.
func (t TeamID) Int() int { return int(t) }
