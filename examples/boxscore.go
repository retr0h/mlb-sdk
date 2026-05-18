// Copyright (c) 2026 John Dewey
//
// SPDX-License-Identifier: MIT
//
// Boxscore prints the team-stats summary for a single game, including the
// number of double plays each side turned. Run with:
//
//	go run ./examples/boxscore <gamePk>
package main
//
import (
	"context"
	"fmt"
	"os"
	"strconv"
//
	"github.com/retr0h/mlb-sdk/pkg/mlb"
)
//
func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: boxscore <gamePk>")
		os.Exit(2)
	}
	gamePk, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "gamePk must be an integer")
		os.Exit(2)
	}
//
	c := mlb.New()
	box, err := c.Boxscore(context.Background(), gamePk)
	if err != nil {
		fmt.Fprintln(os.Stderr, "boxscore:", err)
		os.Exit(1)
	}
	for _, side := range []*mlb.BoxscoreTeam{box.Home, box.Away} {
		if side == nil {
			continue
		}
		fmt.Printf("%-30s  DPs turned: %d\n", side.Name, side.DoublePlaysTurned())
	}
}
