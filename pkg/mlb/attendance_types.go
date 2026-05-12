// Copyright (c) 2026 John Dewey
//
// SPDX-License-Identifier: MIT

package mlb

import "time"

// AttendanceQuery filters an attendance lookup. The MLB API requires one of
// TeamID, LeagueID, or LeagueListID to be set (toddrob99 encodes this as
// `required_params: [["teamId"], ["leagueId"], ["leagueListId"]]`). The
// SDK enforces this with an ErrInvalidQuery runtime check.
type AttendanceQuery struct {
	// TeamID restricts attendance to a single team. One of TeamID,
	// LeagueID, or LeagueListID is required.
	TeamID int

	// LeagueID restricts attendance to a single league. One of TeamID,
	// LeagueID, or LeagueListID is required.
	LeagueID int

	// LeagueListID picks a named league list (e.g. "milb_all"). One of
	// TeamID, LeagueID, or LeagueListID is required.
	LeagueListID string

	// Season constrains to a single season year.
	Season int

	// On views attendance through a specific date (the MLB API's `date`
	// query — year-to-date totals as of this calendar date).
	On time.Time

	// GameType filters by game-type code: "R" (regular), "S" (spring), …
	GameType string

	// Fields restricts the response to a comma-separated field projection.
	Fields string
}

// Attendance is the typed view of /api/v1/attendance. AggregateTotals is the
// per-call rollup the API always returns alongside the per-record details.
type Attendance struct {
	Records         []AttendanceRecord
	AggregateTotals AttendanceAggregateTotals
}

// AttendanceRecord is one attendance row — typically per (team, season,
// gameType) combination, or aggregated per league.
type AttendanceRecord struct {
	OpeningsTotal            int
	OpeningsTotalAway        int
	OpeningsTotalHome        int
	OpeningsTotalLost        int
	GamesTotal               int
	GamesAwayTotal           int
	GamesHomeTotal           int
	Year                     string
	AttendanceAverageAway    int
	AttendanceAverageHome    int
	AttendanceAverageYtd     int
	AttendanceHigh           int
	AttendanceHighDate       time.Time
	AttendanceHighGame       AttendanceGameRef
	AttendanceLow            int
	AttendanceLowDate        time.Time
	AttendanceLowGame        AttendanceGameRef
	AttendanceOpeningAverage int
	AttendanceTotal          int
	AttendanceTotalAway      int
	AttendanceTotalHome      int
	GameType                 GameTypeRef
	Team                     TeamRef
}

// AttendanceGameRef points at one of the games anchoring a high/low
// attendance record. Content is a nested link to the game-content endpoint.
type AttendanceGameRef struct {
	GamePk   int
	Link     string
	Content  Ref
	DayNight string // "day" | "night"
}

// GameTypeRef is the short reference `{id, description}` for an MLB
// game-type. Used by attendance records to label the gameType bucket.
type GameTypeRef struct {
	ID          string // "R", "S", "E", …
	Description string
}

// AttendanceAggregateTotals is the per-call rollup the API returns alongside
// each attendance response — sum / average across every record returned.
type AttendanceAggregateTotals struct {
	OpeningsTotalAway     int
	OpeningsTotalHome     int
	OpeningsTotalLost     int
	OpeningsTotalYtd      int
	AttendanceAverageAway int
	AttendanceAverageHome int
	AttendanceAverageYtd  int
	AttendanceHigh        int
	AttendanceHighDate    time.Time
	AttendanceTotal       int
	AttendanceTotalAway   int
	AttendanceTotalHome   int
}
