// Copyright (c) 2026 John Dewey
//
// SPDX-License-Identifier: MIT

// Team prints a single MLB team's metadata. Run with:
//
//	go run ./examples/team [TEAM_ID]   # default: 119 (Dodgers)
package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/retr0h/mlb-sdk/pkg/mlb"
)

func main() {
	id := 119
	if len(os.Args) > 1 {
		n, err := strconv.Atoi(os.Args[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "teamId must be numeric: %v\n", err)
			os.Exit(2)
		}
		id = n
	}

	c := mlb.New()
	t, err := c.Team(context.Background(), id, mlb.TeamQuery{
		Hydrate: "league,division,sport,springLeague,venue",
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "team:", err)
		os.Exit(1)
	}
	fmt.Printf("%s (%s) — id=%d, firstYearOfPlay=%s\n",
		t.Name, t.Abbreviation, t.ID, t.FirstYearOfPlay)
	fmt.Printf("  venue:    %s\n", t.Venue.Name)
	fmt.Printf("  league:   %s (%s)\n", t.League.Name, t.League.Abbreviation)
	fmt.Printf("  division: %s\n", t.Division.Name)
	fmt.Printf("  sport:    %s\n", t.Sport.Name)
	fmt.Printf("  spring:   %s\n", t.SpringLeague.Name)
}
