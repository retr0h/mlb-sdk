// Copyright (c) 2026 John Dewey
//
// SPDX-License-Identifier: MIT
//
// TeamsHistory prints historical team records. Run with:
//
//	go run ./examples/teamshistory
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
	ts, err := c.TeamsHistory(context.Background(), mlb.TeamsHistoryQuery{
		TeamIDs: "119",
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "teamsHistory:", err)
		os.Exit(1)
	}
	for _, t := range ts.Teams {
		fmt.Printf("%-4s %-30s season=%d  active=%v\n",
			t.Abbreviation, t.Name, t.Season, t.Active)
	}
}
