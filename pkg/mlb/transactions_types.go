// Copyright (c) 2026 John Dewey
//
// SPDX-License-Identifier: MIT

package mlb

import "time"

// TransactionsQuery filters a transactions lookup. toddrob99 encodes
// `required_params: [["teamId"], ["playerId"], ["date"], ["startDate",
// "endDate"]]` — one of the following must hold:
//
//   - TeamID alone,
//   - PlayerID alone,
//   - On alone,
//   - StartDate AND EndDate together.
//
// The SDK enforces this with an ErrInvalidQuery runtime check.
type TransactionsQuery struct {
	// TeamID restricts transactions to a single team.
	TeamID int

	// PlayerID restricts transactions to a single player.
	PlayerID int

	// On is the single calendar day (the MLB API's `date` query).
	On time.Time

	// StartDate / EndDate bracket a closed date range. Both must be set
	// together — one alone is rejected.
	StartDate time.Time
	EndDate   time.Time

	// SportID restricts to a sport (1 = MLB, 11 = AAA, …).
	SportID int

	// Fields restricts the response to a comma-separated field projection.
	Fields string
}

// Transactions is the typed view of /api/v1/transactions.
type Transactions struct {
	Transactions []Transaction
}

// Transaction is one roster / assignment transaction. FromTeam is omitted on
// initial assignments; Description is the free-text the MLB feed populates.
// Date fields are kept as strings because the MLB API delivers them as
// YYYY-MM-DD without time-of-day context.
type Transaction struct {
	ID             int
	Person         Person
	FromTeam       TeamRef
	ToTeam         TeamRef
	Date           string // YYYY-MM-DD
	EffectiveDate  string // YYYY-MM-DD
	ResolutionDate string // YYYY-MM-DD
	TypeCode       string // "ASG", "SE", "TR", …
	TypeDesc       string
	Description    string
}
