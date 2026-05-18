// Copyright (c) 2026 John Dewey
//
// SPDX-License-Identifier: MIT
//
package mlb
//
import "time"
//
// SeasonsQuery filters a seasons lookup. At least one of SportID, DivisionID
// or LeagueID is required — the toddrob99 catalog encodes this constraint as
// `required_params: [["sportId"], ["divisionId"], ["leagueId"]]`. Set All to
// true to query the full history (the `/api/v1/seasons/all` path).
type SeasonsQuery struct {
	// Season constrains the response to a specific year (e.g. 2024). Ignored
	// when All is true.
	Season int
//
	// SportID restricts to a sport (1 = MLB, 11 = AAA, …).
	SportID int
//
	// DivisionID restricts to a single division.
	DivisionID int
//
	// LeagueID restricts to a single league (103 = AL, 104 = NL).
	LeagueID int
//
	// Fields is a comma-separated field projection passed to the MLB API.
	Fields string
//
	// All switches the underlying path to `/api/v1/seasons/all`, returning
	// every season the API tracks. When set, Season is ignored.
	All bool
}
//
// Seasons is the typed view of /api/v1/seasons (and /api/v1/seasons/all).
type Seasons struct {
	Seasons []Season
}
//
// Season returns the entry with the given seasonId ("2026", "1876", …),
// or nil when not present in the response.
func (s *Seasons) Season(seasonID string) *Season {
	if s == nil {
		return nil
	}
	for i := range s.Seasons {
		if s.Seasons[i].SeasonID == seasonID {
			return &s.Seasons[i]
		}
	}
	return nil
}
//
// Season is one season's metadata. Date fields are parsed from the API's
// YYYY-MM-DD strings; zero time.Time means the API omitted the field for
// that historical row (most pre-divisional seasons omit allStarDate,
// firstDate2ndHalf, etc.). SeasonID is kept as a string because the MLB API
// uses it as an opaque key.
type Season struct {
	SeasonID                  string
	HasWildcard               bool
	PreSeasonStartDate        time.Time
	PreSeasonEndDate          time.Time
	SeasonStartDate           time.Time
	SpringStartDate           time.Time
	SpringEndDate             time.Time
	RegularSeasonStartDate    time.Time
	LastDate1stHalf           time.Time
	AllStarDate               time.Time
	FirstDate2ndHalf          time.Time
	RegularSeasonEndDate      time.Time
	PostSeasonStartDate       time.Time
	PostSeasonEndDate         time.Time
	SeasonEndDate             time.Time
	OffseasonStartDate        time.Time
	OffSeasonEndDate          time.Time
	SeasonLevelGamedayType    string // e.g. "P" (post-1968), "S" (Statcast-era)
	GameLevelGamedayType      string
	QualifierPlateAppearances float64
	QualifierOutsPitched      float64
}
