// Copyright (c) 2026 John Dewey
//
// SPDX-License-Identifier: MIT

package mlb

import "time"

// StandingsQuery filters a standings lookup. League is required; the rest are
// optional. Set Date to view standings as of a specific calendar date.
type StandingsQuery struct {
	// League is the MLB league to fetch standings for (required).
	League LeagueID

	// Season is the year (e.g. 2026). Zero = current season per MLB.
	Season int

	// StandingsTypes is the API's `standingsTypes` query — e.g.
	// "regularSeason", "wildCard", "divisionLeaders". Zero value = default.
	StandingsTypes string

	// On views standings as of a specific date. Zero = current.
	On time.Time

	// Hydrate is a comma-separated free-form hydrate string. The MLB API's
	// hydrate vocabulary grows constantly; we don't constrain it.
	Hydrate string
}

// Standings is the typed view of /api/v1/standings — one outer slice of
// per-division standings, each containing a team-record list.
type Standings struct {
	Records []DivisionStandings
}

// Division returns the division-standings entry for a given division id, or
// nil when not present in this response.
func (s *Standings) Division(id int) *DivisionStandings {
	if s == nil {
		return nil
	}
	for i := range s.Records {
		if s.Records[i].Division.ID == id {
			return &s.Records[i]
		}
	}
	return nil
}

// DivisionStandings is one division block in a standings response.
type DivisionStandings struct {
	StandingsType string
	League        Ref
	Division      Ref
	Sport         Ref
	LastUpdated   time.Time
	TeamRecords   []TeamRecord
}

// Team returns the TeamRecord for a given TeamID within this division, or
// nil when the team isn't part of this division block.
func (d *DivisionStandings) Team(id TeamID) *TeamRecord {
	if d == nil {
		return nil
	}
	for i := range d.TeamRecords {
		if d.TeamRecords[i].Team.ID == id {
			return &d.TeamRecords[i]
		}
	}
	return nil
}

// Ref is the API's lightweight reference object — `{id, link}` — used for
// league, division, and sport pointers. The MLB API doesn't include a name
// in these references; resolve via the Sport/League/Division endpoints if
// you need one.
type Ref struct {
	ID   int
	Link string
}

// TeamRecord is one team's standings slot. Several rank fields are typed as
// strings because the MLB API uses "-" for "not applicable" — an int can't
// represent that without losing information.
type TeamRecord struct {
	Team              TeamRef
	Streak            Streak
	Wins              int
	Losses            int
	GamesPlayed       int
	RunsScored        int
	RunsAllowed       int
	RunDifferential   int
	WinningPercentage string // ".593"-style decimal
	GamesBack         string // "-" or "1.5"
	WildCardGamesBack string
	DivisionRank      string
	LeagueRank        string
	SportRank         string
	EliminationNumber string
	MagicNumber       string
	Clinched          bool
	DivisionLeader    bool
	DivisionChamp     bool
	HasWildcard       bool
	Season            string
	LastUpdated       time.Time
}

// TeamRef is a lightweight team pointer used in standings (and other
// endpoints where the API returns just `{id, name}` without the full Team
// detail).
type TeamRef struct {
	ID   TeamID
	Name string
}

// Streak is a team's current win/loss streak as the API reports it.
type Streak struct {
	Code   string // "W3", "L1"
	Type   string // "wins", "losses"
	Number int
}
