// Copyright (c) 2026 John Dewey
//
// SPDX-License-Identifier: MIT

package mlb

// AwardRecipientsQuery refines an award-recipients lookup. The awardId path
// parameter is taken as a method argument; everything below is optional.
type AwardRecipientsQuery struct {
	// SportID restricts to a sport (1 = MLB).
	SportID int

	// LeagueID restricts to a league (103 = AL, 104 = NL).
	LeagueID int

	// Season restricts to a single season year.
	Season int

	// Hydrate is a comma-separated hydrate string. The MLB API's hydrate
	// vocabulary grows; we don't constrain it.
	Hydrate string

	// Fields restricts the response to a comma-separated field projection.
	Fields string
}

// Awards is the typed view of /api/v1/awards/{awardId}/recipients — a list
// of recipient rows.
type Awards struct {
	Recipients []AwardRecipient
}

// AwardRecipient is one award-recipient row — award metadata plus the team
// and person it was awarded to.
type AwardRecipient struct {
	ID     string // award id, e.g. "MLBHOF"
	Name   string // award display name
	Date   string // YYYY-MM-DD as the API delivers it (date-only string)
	Season string
	Team   Ref
	Player Person
	Notes  string
}

// Person is the lightweight person reference returned by awards /
// transactions / roster endpoints. PrimaryPosition is only present when the
// request hydrates the position sub-object.
type Person struct {
	ID              int
	FullName        string
	NameFirstLast   string
	Link            string
	PrimaryPosition PrimaryPosition
}

// PrimaryPosition is the player-position metadata the MLB API embeds in
// person references.
type PrimaryPosition struct {
	Code         string // "1", "2", …
	Name         string // "Pitcher", "Catcher", …
	Type         string // "Pitcher", "Infielder", …
	Abbreviation string // "P", "C", "1B", …
}
