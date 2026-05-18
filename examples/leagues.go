// Copyright (c) 2026 John Dewey
//
// SPDX-License-Identifier: MIT
//
// Leagues prints leagues for a sport. Run with:
//
//	go run ./examples/leagues   # MLB AL/NL
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
	ls, err := c.Leagues(context.Background(), mlb.LeaguesQuery{SportID: 1})
	if err != nil {
		fmt.Fprintln(os.Stderr, "leagues:", err)
		os.Exit(1)
	}
	for _, l := range ls.Leagues {
		fmt.Printf("%3d  %-20s (%s) teams=%d wildcards=%d\n",
			l.ID, l.Name, l.Abbreviation, l.NumTeams, l.NumWildcardTeams)
	}
}
