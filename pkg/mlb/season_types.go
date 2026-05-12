// Copyright (c) 2026 John Dewey
//
// SPDX-License-Identifier: MIT

package mlb

// SeasonQuery refines a single-season lookup. The seasonId path parameter is
// taken as a method argument; SportID is required by the MLB API and must be
// non-zero.
type SeasonQuery struct {
	// SportID is the sport the season belongs to (1 = MLB, 11 = AAA, …).
	// Required.
	SportID int

	// Fields restricts the response to a comma-separated field projection.
	Fields string
}
