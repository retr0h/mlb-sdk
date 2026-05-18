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

import "time"

// ScheduleQuery filters the schedule listing. All fields are optional; the
// zero value returns the full league schedule for the configured day.
type ScheduleQuery struct {
	// Team filters to a single franchise. Zero means no team filter.
	Team TeamID

	// On filters to a single calendar date. Zero means no date filter.
	On time.Time

	// From / To bound a date range (inclusive). Both must be set to apply.
	From time.Time
	To   time.Time
}

// GameStatus is the MLB API's coarse-grained game state.
type GameStatus string

// Coarse-grained game states reported by the MLB API.
const (
	StatusFinal   GameStatus = "Final"
	StatusLive    GameStatus = "Live"
	StatusPreview GameStatus = "Preview"
)

// TeamScore is one side of a scheduled matchup.
type TeamScore struct {
	ID    TeamID
	Name  string
	Score int
}

// Game is a single scheduled MLB game in idiomatic form. The MLB API's raw
// gameDate string has been parsed into a UTC time.Time.
type Game struct {
	GamePk int
	Date   time.Time

	// Status is the coarse-grained game state (Final / Live / Preview).
	Status GameStatus

	// DetailedStatus is the API's free-form detailed state string —
	// values like "Final", "In Progress", "Scheduled", "Postponed",
	// "Suspended: Rain". Use this when the coarse Status isn't enough.
	DetailedStatus string

	Home TeamScore
	Away TeamScore
}
