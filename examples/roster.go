// Copyright (c) 2026 John Dewey
//
// SPDX-License-Identifier: MIT

// Roster prints a team's active roster. Run with:
//
//	go run ./examples/roster [TEAM_ID] [SEASON]   # default: 119 2024
package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/retr0h/mlb-sdk/pkg/mlb"
)

func main() {
	teamID, year := 119, 2024
	if len(os.Args) > 1 {
		n, _ := strconv.Atoi(os.Args[1])
		teamID = n
	}
	if len(os.Args) > 2 {
		n, _ := strconv.Atoi(os.Args[2])
		year = n
	}

	c := mlb.New()
	r, err := c.Roster(context.Background(), teamID, mlb.RosterQuery{
		RosterType: "active",
		Season:     year,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "roster:", err)
		os.Exit(1)
	}
	for _, e := range r.Roster {
		fmt.Printf("#%-3s %-25s [%s]  %s\n",
			e.JerseyNumber, e.Person.FullName,
			e.Position.Abbreviation, e.Status.Description)
	}
}
