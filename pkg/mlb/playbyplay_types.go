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

// EventType is the MLB Stats API's machine-readable play classifier. Reported
// at result.eventType in the live feed and play-by-play responses. Constants
// below cover the common values; values outside this set come through as
// EventType(<raw string>) so callers can still inspect them.
type EventType string

// Common play event types. The full upstream vocabulary is larger; add new
// constants when a deal or downstream feature needs to match on one.
const (
	EventSingle                 EventType = "single"
	EventDouble                 EventType = "double"
	EventTriple                 EventType = "triple"
	EventHomeRun                EventType = "home_run"
	EventWalk                   EventType = "walk"
	EventStrikeout              EventType = "strikeout"
	EventHitByPitch             EventType = "hit_by_pitch"
	EventFieldError             EventType = "field_error"
	EventFieldOut               EventType = "field_out"
	EventForceOut               EventType = "force_out"
	EventGroundedIntoDoublePlay EventType = "grounded_into_double_play"
)

// HalfInning marks whether a play occurred in the top or bottom of an inning.
type HalfInning string

// Top / bottom of an inning, as reported by the MLB API.
const (
	HalfTop    HalfInning = "top"
	HalfBottom HalfInning = "bottom"
)

// Play is a single batter-vs-pitcher event in a game, flattened from the
// nested API shape (result, count, about) into a single struct.
type Play struct {
	// Event is the human-readable name (e.g. "Grounded Into DP").
	Event string

	// EventType is the machine-readable classifier (e.g.
	// "grounded_into_double_play"). Match against the Event* constants
	// rather than string literals.
	EventType EventType

	// Description is the narrative line shown in MLB Gameday — usually
	// the most useful string for free-text searches.
	Description string

	// Inning is 1-indexed (1 through 9 in regulation; 10+ in extras).
	Inning int

	// HalfInning is "top" or "bottom" — see Half* constants.
	HalfInning HalfInning

	// Outs reflects the total number of outs after this play. Use
	// successive Plays to derive per-play outs recorded.
	Outs int
}

// IsDoublePlay reports whether this play was officially scored as a
// grounded-into-double-play. Other plays where two outs are recorded (e.g.
// caught stealing on a strikeout) are not flagged here — the MLB API
// distinguishes them by eventType.
func (p Play) IsDoublePlay() bool {
	return p.EventType == EventGroundedIntoDoublePlay
}
