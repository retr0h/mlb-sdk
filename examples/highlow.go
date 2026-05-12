// Copyright (c) 2026 John Dewey
//
// SPDX-License-Identifier: MIT

// HighLow prints season high/low records. Run with:
//
//	go run ./examples/highlow   # top HR leaders 2024
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/retr0h/mlb-sdk/pkg/mlb"
)

func main() {
	c := mlb.New()
	hl, err := c.HighLow(context.Background(), "player", mlb.HighLowQuery{
		SortStat: "homeRuns",
		Season:   2024,
		SportIDs: "1",
		Limit:    10,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "highLow:", err)
		os.Exit(1)
	}
	for _, g := range hl.Results {
		fmt.Printf("group: %s (%d total)\n", g.Group, g.TotalSplits)
		for _, s := range g.Splits {
			fmt.Printf("  %-25s %s  %v\n", s.Player.FullName, s.Team.Name, s.Stat)
		}
	}
}
