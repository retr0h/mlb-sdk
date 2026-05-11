// Copyright (c) 2026 John Dewey
//
// SPDX-License-Identifier: MIT

// Teams prints every MLB team for a season. Run with:
//
//	go run ./examples/teams [SEASON]   # default: 2024
package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/retr0h/mlb-sdk/pkg/mlb"
)

func main() {
	year := 2024
	if len(os.Args) > 1 {
		n, err := strconv.Atoi(os.Args[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "season must be numeric: %v\n", err)
			os.Exit(2)
		}
		year = n
	}

	c := mlb.New()
	ts, err := c.Teams(context.Background(), mlb.TeamsQuery{
		SportID: 1,
		Season:  year,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "teams:", err)
		os.Exit(1)
	}
	for _, t := range ts.Teams {
		fmt.Printf("%-4s %-25s id=%d  venue=%s\n",
			t.Abbreviation, t.Name, t.ID, t.Venue.Name)
	}
}
