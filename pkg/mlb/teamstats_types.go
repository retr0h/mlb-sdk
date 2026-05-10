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
	"strings"

	"github.com/retr0h/mlb-sdk/internal/gen"
)

// TeamStatGroup is the MLB Stats API's `group` query parameter — the side of
// the game the requested stats are about.
type TeamStatGroup string

// Standard team-stat groups. The MLB API also accepts other values (e.g.
// "catching"); add constants here when a downstream feature needs to match
// on one.
const (
	TeamStatGroupHitting  TeamStatGroup = "hitting"
	TeamStatGroupPitching TeamStatGroup = "pitching"
	TeamStatGroupFielding TeamStatGroup = "fielding"
)

// TeamStatType is the MLB Stats API's `stats` query parameter — the time
// window for the requested stats.
type TeamStatType string

// Common team-stat windows.
const (
	TeamStatTypeSeason      TeamStatType = "season"
	TeamStatTypeByDateRange TeamStatType = "byDateRange"
	TeamStatTypeYearByYear  TeamStatType = "yearByYear"
)

// TeamStatsQuery filters a team-stats lookup. Team is required; the rest are
// optional but typically all set for meaningful results.
type TeamStatsQuery struct {
	// Team is the team whose stats to fetch (required).
	Team TeamID

	// Season is the year (e.g. 2026). Zero leaves the param off, in which
	// case the API defaults to the current season.
	Season int

	// Type is the time-window classifier (e.g. season, byDateRange).
	Type TeamStatType

	// Group is the side-of-game classifier (hitting, pitching, fielding).
	Group TeamStatGroup
}

// TeamStats is the typed view of /api/v1/teams/{teamId}/stats. The response
// is naturally a list of stat groups (one per group×type combination
// requested), each containing one or more season splits.
type TeamStats struct {
	Groups []TeamStatGroupResult

	raw *gen.TeamStatsResponse
}

// Group returns the first group in the response whose group-name matches g
// (case-insensitively), or nil if no such group is present. Use this when
// you queried for a single group and want a direct path to its splits.
func (t *TeamStats) Group(g TeamStatGroup) *TeamStatGroupResult {
	if t == nil {
		return nil
	}
	for i := range t.Groups {
		if strings.EqualFold(t.Groups[i].Group, string(g)) {
			return &t.Groups[i]
		}
	}
	return nil
}

// TeamStatGroupResult is a single (type × group) cell of the response —
// e.g. season-fielding, yearByYear-hitting — with one or more splits inside.
type TeamStatGroupResult struct {
	// Type is the displayName of the time window (e.g. "season").
	Type string
	// Group is the displayName of the side (e.g. "fielding").
	Group string
	// Splits is the list of per-window rows. For Type="season" this is
	// usually a single row; for "yearByYear" it's many.
	Splits []TeamStatsSplit
}

// Season returns the split matching the given season string (e.g. "2026"),
// or nil if no such split exists.
func (g *TeamStatGroupResult) Season(season string) *TeamStatsSplit {
	if g == nil {
		return nil
	}
	for i := range g.Splits {
		if g.Splits[i].Season == season {
			return &g.Splits[i]
		}
	}
	return nil
}

// TeamStatsSplit is a single time-window row of stats. The MLB API returns
// the stat block as a free-form object whose keys vary by group; we keep it
// as a map and provide typed accessors.
type TeamStatsSplit struct {
	Season string

	stat map[string]any
}

// Int reads key from the underlying stat map and coerces it to int. Returns
// 0 when the key is missing, the receiver is nil, or the value is not a
// number-like type.
func (s *TeamStatsSplit) Int(key string) int {
	if s == nil {
		return 0
	}
	v, ok := s.stat[key]
	if !ok {
		return 0
	}
	switch n := v.(type) {
	case float64:
		return int(n)
	case int:
		return n
	case int64:
		return int(n)
	default:
		return 0
	}
}

// Float reads key from the underlying stat map as float64. Returns 0 when
// the key is missing or the value is not a JSON number.
func (s *TeamStatsSplit) Float(key string) float64 {
	if s == nil {
		return 0
	}
	v, ok := s.stat[key]
	if !ok {
		return 0
	}
	f, _ := v.(float64)
	return f
}

// String reads key from the underlying stat map as string. Returns ""
// when the key is missing or the value is not a JSON string.
func (s *TeamStatsSplit) String(key string) string {
	if s == nil {
		return ""
	}
	v, ok := s.stat[key]
	if !ok {
		return ""
	}
	str, _ := v.(string)
	return str
}

// DoublePlays returns the doublePlays field from this split. Equivalent to
// s.Int("doublePlays") but typed for discoverability.
func (s *TeamStatsSplit) DoublePlays() int { return s.Int("doublePlays") }
