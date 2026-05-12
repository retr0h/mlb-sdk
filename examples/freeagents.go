// Copyright (c) 2026 John Dewey
//
// SPDX-License-Identifier: MIT

// FreeAgents prints free-agent signings for a season. Run with:
//
//	go run ./examples/freeagents [SEASON]   # default: 2024
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
	fa, err := c.FreeAgents(context.Background(), mlb.FreeAgentsQuery{Season: year})
	if err != nil {
		fmt.Fprintln(os.Stderr, "freeAgents:", err)
		os.Exit(1)
	}
	for _, f := range fa.FreeAgents {
		to := "(unsigned)"
		if f.NewTeam.Name != "" {
			to = f.NewTeam.Name
		}
		fmt.Printf("%-25s [%s]  %s → %s  %s\n",
			f.Player.FullName, f.Position.Abbreviation,
			f.OriginalTeam.Name, to, f.DateSigned)
	}
}
