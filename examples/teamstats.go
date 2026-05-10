// Copyright (c) 2026 John Dewey
//
// SPDX-License-Identifier: MIT

// TeamStats prints the season-fielding row for a team, including double
// plays and errors. Run with:
//
//	go run ./examples/teamstats <season>   # default season: current year
package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/retr0h/mlb-sdk/pkg/mlb"
)

func main() {
	season := time.Now().Year()
	if len(os.Args) > 1 {
		s, err := strconv.Atoi(os.Args[1])
		if err != nil {
			fmt.Fprintln(os.Stderr, "season must be an integer")
			os.Exit(2)
		}
		season = s
	}

	c := mlb.New()
	ts, err := c.TeamStats(context.Background(), mlb.TeamStatsQuery{
		Team:   mlb.LAD,
		Season: season,
		Type:   mlb.TeamStatTypeSeason,
		Group:  mlb.TeamStatGroupFielding,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "teamStats:", err)
		os.Exit(1)
	}
	split := ts.Group(mlb.TeamStatGroupFielding).Season(strconv.Itoa(season))
	if split == nil {
		fmt.Println("no fielding data for that season")
		return
	}
	fmt.Printf("Dodgers %d fielding — DPs: %d, errors: %d\n",
		season, split.DoublePlays(), split.Int("errors"))
}
