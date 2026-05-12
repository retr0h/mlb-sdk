// Copyright (c) 2026 John Dewey
//
// SPDX-License-Identifier: MIT

package mlb

// PersonQuery refines a single-person lookup. The personId path parameter is
// taken as a method argument; everything in this struct is optional.
type PersonQuery struct {
	// Hydrate is a comma-separated hydrate string.
	Hydrate string

	// Fields restricts the response to a comma-separated field projection.
	Fields string
}

// PeopleQuery filters a multi-person lookup. PersonIDs is required.
type PeopleQuery struct {
	// PersonIDs is a comma-separated list of person ids (required).
	PersonIDs string

	// Hydrate is a comma-separated hydrate string.
	Hydrate string

	// Fields restricts the response to a comma-separated field projection.
	Fields string
}

// PersonDetail is the rich person record returned by /api/v1/people/{personId}.
// Many fields are absent for non-player persons (coaches, umpires, etc.).
type PersonDetail struct {
	ID               int
	FullName         string
	FirstName        string
	LastName         string
	Link             string
	PrimaryNumber    string
	BirthDate        string // YYYY-MM-DD
	CurrentAge       int
	BirthCity        string
	BirthCountry     string
	Height           string // "6' 4\""
	Weight           int
	Active           bool
	PrimaryPosition  PrimaryPosition
	UseName          string
	UseLastName      string
	BoxscoreName     string
	NickName         string
	Gender           string
	IsPlayer         bool
	IsVerified       bool
	Pronunciation    string
	MlbDebutDate     string // YYYY-MM-DD
	BatSide          HandSide
	PitchHand        HandSide
	NameFirstLast    string
	NameSlug         string
	FirstLastName    string
	LastFirstName    string
	LastInitName     string
	InitLastName     string
	FullFMLName      string
	FullLFMName      string
	StrikeZoneTop    float64
	StrikeZoneBottom float64
}

// HandSide is the `{code, description}` pair for bat side or pitch hand.
type HandSide struct {
	Code        string // "L", "R", "S"
	Description string // "Left", "Right", "Switch"
}
