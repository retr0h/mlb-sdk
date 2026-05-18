// Copyright (c) 2026 John Dewey
//
// SPDX-License-Identifier: MIT
//
// Schedule prints today's MLB games for a given team. Run with:
//
//	go run ./examples/schedule [TEAM]   # default: LAD
package main
//
import (
	"context"
	"fmt"
	"os"
	"time"
//
	"github.com/retr0h/mlb-sdk/pkg/mlb"
)
//
func main() {
	team := mlb.LAD
	if len(os.Args) > 1 {
		// For brevity this example only knows the Dodgers/Yankees/Braves;
		// a real CLI would map the full string→TeamID set.
		switch os.Args[1] {
		case "NYY":
			team = mlb.NYY
		case "ATL":
			team = mlb.ATL
		case "LAD":
			team = mlb.LAD
		default:
			fmt.Fprintf(os.Stderr, "unknown team %q (LAD|NYY|ATL)\n", os.Args[1])
			os.Exit(2)
		}
	}
//
	c := mlb.New()
	games, err := c.Schedule(context.Background(), mlb.ScheduleQuery{
		Team: team,
		On:   time.Now(),
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "schedule:", err)
		os.Exit(1)
	}
	if len(games) == 0 {
		fmt.Println("no games today")
		return
	}
	for _, g := range games {
		fmt.Printf("%-7d  %s @ %s  (%s)\n",
			g.GamePk, g.Away.Name, g.Home.Name, g.Status)
	}
}
