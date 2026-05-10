// Copyright (c) 2026 John Dewey
//
// SPDX-License-Identifier: MIT

// Standings prints the top of the NL East as of today. Run with:
//
//	go run ./examples/standings [LEAGUE]   # AL | NL (default NL)
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/retr0h/mlb-sdk/pkg/mlb"
)

func main() {
	league := mlb.NL
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "AL":
			league = mlb.AL
		case "NL":
			league = mlb.NL
		default:
			fmt.Fprintf(os.Stderr, "unknown league %q (AL|NL)\n", os.Args[1])
			os.Exit(2)
		}
	}

	c := mlb.New()
	st, err := c.Standings(context.Background(), mlb.StandingsQuery{
		League:         league,
		StandingsTypes: "regularSeason",
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "standings:", err)
		os.Exit(1)
	}
	for _, d := range st.Records {
		fmt.Printf("Division %d (league %d):\n", d.DivisionID, d.LeagueID)
		for _, tr := range d.TeamRecords {
			fmt.Printf("  %2s  %-25s  %3d-%-3d  %s  GB %s  streak %s\n",
				tr.DivisionRank, tr.Team.Name,
				tr.Wins, tr.Losses, tr.WinningPercentage,
				tr.GamesBack, tr.Streak.Code)
		}
		fmt.Println()
	}
}
