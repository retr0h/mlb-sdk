// Copyright (c) 2026 John Dewey
//
// SPDX-License-Identifier: MIT

package mlb

// GamePaceQuery filters a game-pace lookup. Season is required (toddrob99:
// required_params=[["season"]]).
type GamePaceQuery struct {
	Season          int
	SportID         int
	TeamIDs         string
	LeagueIDs       string
	LeagueListID    string
	GameType        string
	StartDate       string
	EndDate         string
	VenueIDs        string
	OrgType         string
	IncludeChildren bool
	Fields          string
}

// GamePace is the typed view of /api/v1/gamePace. Time fields are strings
// in HH:MM:SS format as the API delivers them.
type GamePace struct {
	Season                         string
	Sport                          Ref
	HitsPer9Inn                    float64
	RunsPer9Inn                    float64
	PitchesPer9Inn                 float64
	PlateAppearancesPer9Inn        float64
	HitsPerGame                    float64
	RunsPerGame                    float64
	InningsPlayedPerGame           float64
	PitchesPerGame                 float64
	PitchersPerGame                float64
	PlateAppearancesPerGame        float64
	HitsPerRun                     float64
	PitchesPerPitcher              float64
	TotalGameTime                  string
	TotalInningsPlayed             float64
	TotalHits                      int
	TotalRuns                      int
	TotalPlateAppearances          int
	TotalPitchers                  int
	TotalPitches                   int
	TotalGames                     int
	Total7InnGames                 int
	Total9InnGames                 int
	Total9InnGamesCompletedEarly   int
	Total9InnGamesScheduled        int
	Total9InnGamesWithoutExtraInn  int
	TotalExtraInnGames             int
	TotalExtraInnTime              string
	TimePerGame                    string
	TimePerPitch                   string
	TimePerHit                     string
	TimePerRun                     string
	TimePerPlateAppearance         string
	TimePer9Inn                    string
	TimePer77PlateAppearances      string
	TimePer7InnGameWithoutExtraInn string
}
