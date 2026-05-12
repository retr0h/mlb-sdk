// Copyright (c) 2026 John Dewey
//
// SPDX-License-Identifier: MIT

// GamePace prints pace-of-play stats for an MLB season. Run with:
//
//	go run ./examples/gamepace [SEASON]   # default: 2024
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
	gp, err := c.GamePace(context.Background(), mlb.GamePaceQuery{
		Season:  year,
		SportID: 1,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "gamePace:", err)
		os.Exit(1)
	}
	fmt.Printf("Pace of Play — %s\n", gp.Season)
	fmt.Printf("  games: %d  total time: %s\n", gp.TotalGames, gp.TotalGameTime)
	fmt.Printf("  time per game: %s  per 9 inn: %s\n", gp.TimePerGame, gp.TimePer9Inn)
	fmt.Printf("  pitches/game: %.1f  hits/game: %.1f  runs/game: %.1f\n",
		gp.PitchesPerGame, gp.HitsPerGame, gp.RunsPerGame)
}
