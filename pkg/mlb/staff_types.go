// Copyright (c) 2026 John Dewey
//
// SPDX-License-Identifier: MIT
//
package mlb
//
import "time"
//
// CoachesQuery refines a team-coaches lookup.
type CoachesQuery struct {
	Season int
	On     time.Time
	Fields string
}
//
// PersonnelQuery refines a team-personnel lookup.
type PersonnelQuery struct {
	On     time.Time
	Fields string
}
//
// UmpiresQuery refines an umpires lookup.
type UmpiresQuery struct {
	SportID int
	On      time.Time
	Fields  string
}
//
// Staff is the typed view returned by coaches, personnel, and umpires
// endpoints — a list of StaffEntry records.
type Staff struct {
	Link       string
	TeamID     int
	RosterType string
	Roster     []StaffEntry
}
//
// StaffEntry is one person in a staff/umpire roster.
type StaffEntry struct {
	Person       Person
	JerseyNumber string
	Job          string // "Manager", "Umpire", …
	JobID        string // "MNGR", "UMPR", …
	Title        string
}
