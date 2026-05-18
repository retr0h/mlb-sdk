// Copyright (c) 2026 John Dewey
//
// SPDX-License-Identifier: MIT
//
package mlb
//
// TeamQuery refines a single-team lookup. The teamId path parameter is
// taken as a method argument; everything in this struct is optional.
type TeamQuery struct {
	// Season constrains the team metadata to a specific season (some fields
	// — venue, league context — vary year-over-year).
	Season int
//
	// SportID restricts to a sport (1 = MLB, 11 = AAA, …). Rarely needed
	// for /teams/{id} because the id is globally unique, but accepted.
	SportID int
//
	// Hydrate is a comma-separated hydrate string. Common values for this
	// endpoint: "league", "division", "sport", "springLeague", "venue".
	// Without hydration, sub-objects carry only id/name/link.
	Hydrate string
//
	// Fields restricts the response to a comma-separated field projection.
	Fields string
}
//
// TeamInfo is the rich team record returned by /api/v1/teams/{teamId} and
// /api/v1/teams. Sub-objects (League, Division, Sport, SpringLeague, Venue)
// only carry id/name/link until the request hydrates the matching key —
// pass them via TeamQuery.Hydrate or TeamsQuery.Hydrate.
//
// Distinct from TeamRef (the lightweight {id, name} pointer used in
// standings) and TeamID (the typed numeric identifier constant).
type TeamInfo struct {
	ID              int
	Name            string
	Link            string
	Season          int
	Venue           Venue
	SpringLeague    LeagueInfo
	SpringVenue     Ref
	TeamCode        string
	FileCode        string
	Abbreviation    string
	TeamName        string
	LocationName    string
	FirstYearOfPlay string
	League          LeagueInfo
	Division        Division
	Sport           Sport
	ShortName       string
	FranchiseName   string
	ClubName        string
	Active          bool
	AllStarStatus   string
}
//
// LeagueInfo is the rich league record used both by /api/v1/league and as a
// nested object in /api/v1/teams responses. Most fields are zero until the
// API hydrates the league sub-object — embedded references typically only
// carry ID/Name/Link.
type LeagueInfo struct {
	ID               int
	Name             string
	Link             string
	Abbreviation     string
	NameShort        string
	SeasonState      string // "inseason" | "offseason" | …
	HasWildCard      bool
	HasSplitSeason   bool
	NumGames         int
	HasPlayoffPoints bool
	NumTeams         int
	NumWildcardTeams int
	SeasonDateInfo   Season
	Season           string
	OrgCode          string
	ConferencesInUse bool
	DivisionsInUse   bool
	Sport            Ref
	SortOrder        int
	Active           bool
}
