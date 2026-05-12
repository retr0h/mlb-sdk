// Copyright (c) 2026 John Dewey
//
// SPDX-License-Identifier: MIT

package mlb

// HighLowQuery filters a high/low lookup. The orgType path parameter is
// taken as a method argument. toddrob99 marks sortStat+season as required
// together; the SDK enforces this with an ErrInvalidQuery runtime check.
type HighLowQuery struct {
	SortStat  string
	Season    int
	GameType  string
	TeamID    int
	LeagueID  int
	SportIDs  string
	StatGroup string
	Limit     int
	Fields    string
}

// HighLow is the typed view of /api/v1/highLow/{orgType}.
type HighLow struct {
	Results []HighLowGroup
}

// HighLowGroup is one stat-group block (e.g. hitting).
type HighLowGroup struct {
	Group       string // display name
	TotalSplits int
	Splits      []HighLowSplit
}

// HighLowSplit is one stat row. Stat is a free-form map because the
// fields vary by sortStat.
type HighLowSplit struct {
	Season string
	Stat   map[string]any
	Team   TeamRef
	Player Person
}
