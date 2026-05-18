// Copyright (c) 2026 John Dewey
//
// SPDX-License-Identifier: MIT
//
package mlb
//
// StatsLeadersQuery filters a stat-leaders lookup. LeaderCategories is
// required (toddrob99: required_params=[["leaderCategories"]]).
type StatsLeadersQuery struct {
	LeaderCategories string // required, e.g. "homeRuns"
	Season           int
	SportID          int
	LeagueID         int
	StatGroup        string
	PlayerPool       string
	LeaderGameTypes  string
	StatType         string
	Hydrate          string
	Limit            int
	Fields           string
}
//
// StatsLeaders is the typed view of /api/v1/stats/leaders.
type StatsLeaders struct {
	LeagueLeaders []LeaderCategory
}
//
// LeaderCategory is one stat-category block (e.g. homeRuns).
type LeaderCategory struct {
	LeaderCategory string
	Season         string
	GameType       string // display name
	StatGroup      string
	TotalSplits    int
	Leaders        []LeaderEntry
}
//
// LeaderEntry is one leader row.
type LeaderEntry struct {
	Rank   int
	Value  string
	Team   TeamRef
	League Ref
	Player Person
	Sport  Ref
}
