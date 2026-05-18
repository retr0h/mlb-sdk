// Copyright (c) 2026 John Dewey
//
// SPDX-License-Identifier: MIT
//
// StatsLeaders prints league stat leaders. Run with:
//
//	go run ./examples/statsleaders   # top 10 HR leaders 2024
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
	sl, err := c.StatsLeaders(context.Background(), mlb.StatsLeadersQuery{
		LeaderCategories: "homeRuns",
		Season:           2024,
		SportID:          1,
		Limit:            10,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "statsLeaders:", err)
		os.Exit(1)
	}
	for _, cat := range sl.LeagueLeaders {
		fmt.Printf("%s — %s %s\n", cat.LeaderCategory, cat.Season, cat.GameType)
		for _, l := range cat.Leaders {
			fmt.Printf("  #%-3d %-25s %s  (%s)\n",
				l.Rank, l.Player.FullName, l.Value, l.Team.Name)
		}
	}
}
