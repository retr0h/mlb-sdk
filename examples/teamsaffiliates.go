// Copyright (c) 2026 John Dewey
//
// SPDX-License-Identifier: MIT
//
// TeamsAffiliates prints affiliate teams. Run with:
//
//	go run ./examples/teamsaffiliates
package main
//
import (
	"context"
	"fmt"
	"os"
//
	"github.com/retr0h/mlb-sdk/pkg/mlb"
)
//
func main() {
	c := mlb.New()
	ts, err := c.TeamsAffiliates(context.Background(), mlb.TeamsAffiliatesQuery{
		TeamIDs: "119",
		SportID: 1,
		Season:  2024,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "teamsAffiliates:", err)
		os.Exit(1)
	}
	for _, t := range ts.Teams {
		fmt.Printf("%-4s %-30s id=%d\n", t.Abbreviation, t.Name, t.ID)
	}
}
