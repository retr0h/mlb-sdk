// Copyright (c) 2026 John Dewey
//
// SPDX-License-Identifier: MIT
//
package mlb
//
// LinescoreQuery refines a linescore lookup. The gamePk path parameter is
// taken as a method argument; everything in this struct is optional.
type LinescoreQuery struct {
	// Timecode is an YYYYMMDD_HHmmss string for point-in-time linescores.
	Timecode string
//
	// Fields restricts the response to a comma-separated field projection.
	Fields string
}
//
// Linescore is the typed view of /api/v1/game/{gamePk}/linescore.
type Linescore struct {
	CurrentInning        int
	CurrentInningOrdinal string
	InningState          string
	InningHalf           string
	IsTopInning          bool
	ScheduledInnings     int
	Innings              []LinescoreInning
	Teams                LinescoreTeams
	Defense              LinescoreDefense
	Offense              LinescoreOffense
	Balls                int
	Strikes              int
	Outs                 int
}
//
// LinescoreInning is one inning's line in the linescore.
type LinescoreInning struct {
	Num        int
	OrdinalNum string // "1st", "2nd", …
	Home       LinescoreInningHalf
	Away       LinescoreInningHalf
}
//
// LinescoreInningHalf is one side's (home/away) stats for one inning.
type LinescoreInningHalf struct {
	Runs       int
	Hits       int
	Errors     int
	LeftOnBase int
}
//
// LinescoreTeams holds the game-total runs/hits/errors for both sides.
type LinescoreTeams struct {
	Home LinescoreTeamTotals
	Away LinescoreTeamTotals
}
//
// LinescoreTeamTotals is one side's total runs/hits/errors for the game.
type LinescoreTeamTotals struct {
	Runs       int
	Hits       int
	Errors     int
	LeftOnBase int
	IsWinner   bool
}
//
// LinescoreDefense holds the current defensive lineup (Person refs).
type LinescoreDefense struct {
	Pitcher   Person
	Catcher   Person
	First     Person
	Second    Person
	Third     Person
	Shortstop Person
	Left      Person
	Center    Person
	Right     Person
}
//
// LinescoreOffense holds the current offensive situation (Person refs).
type LinescoreOffense struct {
	Batter Person
	OnDeck Person
	InHole Person
	First  Person
	Second Person
	Third  Person
}
