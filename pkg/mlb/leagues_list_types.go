// Copyright (c) 2026 John Dewey
//
// SPDX-License-Identifier: MIT
//
package mlb
//
// LeaguesQuery filters a leagues listing. The MLB API requires either
// SportID or LeagueIDs to be set (toddrob99 encodes this as
// `required_params: [["sportId"], ["leagueIds"]]`); the SDK enforces the
// constraint with an ErrInvalidQuery runtime check.
type LeaguesQuery struct {
	// SportID restricts to a sport (1 = MLB, 11 = AAA, …). One of SportID
	// or LeagueIDs is required.
	SportID int
//
	// LeagueIDs is a comma-separated list of league ids (e.g. "103,104"
	// for AL+NL). One of SportID or LeagueIDs is required.
	LeagueIDs string
//
	// Seasons is a comma-separated list of season years.
	Seasons string
//
	// Fields restricts the response to a comma-separated field projection.
	Fields string
}
//
// Leagues is the typed view of /api/v1/league — a list of LeagueInfo.
type Leagues struct {
	Leagues []LeagueInfo
}
//
// League returns the entry with the given league id, or nil when not
// present in the response.
func (l *Leagues) League(id int) *LeagueInfo {
	if l == nil {
		return nil
	}
	for i := range l.Leagues {
		if l.Leagues[i].ID == id {
			return &l.Leagues[i]
		}
	}
	return nil
}
