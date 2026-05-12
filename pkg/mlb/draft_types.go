// Copyright (c) 2026 John Dewey
//
// SPDX-License-Identifier: MIT

package mlb

// DraftQuery refines a draft lookup. The year path parameter is taken as a
// method argument; everything in this struct is optional.
type DraftQuery struct {
	Round  string
	Fields string
}

// DraftData is the typed view of /api/v1/draft/{year}.
type DraftData struct {
	DraftYear int
	Rounds    []DraftRound
}

// DraftRound is one round of the draft.
type DraftRound struct {
	Round string
	Picks []DraftPick
}

// DraftPick is one pick within a draft round.
type DraftPick struct {
	BisPlayerID       int
	PickRound         string
	PickNumber        int
	DisplayPickNumber int
	RoundPickNumber   int
	Rank              int
	PickValue         string
	SigningBonus      string
	Home              DraftHome
	ScoutingReport    string
	School            DraftSchool
	Blurb             string
	HeadshotLink      string
	Person            PersonDetail
	Team              TeamInfo
	DraftType         DraftTypeRef
	IsDrafted         bool
	IsPass            bool
	Year              string
}

// DraftHome is the hometown location of a draft pick.
type DraftHome struct {
	City    string
	Country string
}

// DraftSchool is the school/college a draft pick attended.
type DraftSchool struct {
	Name        string
	SchoolClass string // "4YR JR", "HS", …
	City        string
	Country     string
	State       string
}

// DraftTypeRef is the short reference for a draft type.
type DraftTypeRef struct {
	Code        string // "JR", …
	Description string
}
