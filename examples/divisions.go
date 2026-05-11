// Copyright (c) 2026 John Dewey
//
// SPDX-License-Identifier: MIT

// Divisions prints MLB-only divisions. Run with:
//
//	go run ./examples/divisions
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/retr0h/mlb-sdk/pkg/mlb"
)

func main() {
	c := mlb.New()
	d, err := c.Divisions(context.Background(), mlb.DivisionsQuery{SportID: 1})
	if err != nil {
		fmt.Fprintln(os.Stderr, "divisions:", err)
		os.Exit(1)
	}
	for _, div := range d.Divisions {
		fmt.Printf(
			"%-4d  %-25s  league=%d  sport=%d  playoff teams=%d  wildcard=%v\n",
			div.ID, div.Name, div.League.ID, div.Sport.ID,
			div.NumPlayoffTeams, div.HasWildcard,
		)
	}
}
