// Copyright (c) 2026 John Dewey
//
// SPDX-License-Identifier: MIT
//
// Linescore prints the inning-by-inning linescore for a game. Run with:
//
//	go run ./examples/linescore [GAME_PK]   # default: 745455
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
	gamePk := 745455
	if len(os.Args) > 1 {
		n, err := strconv.Atoi(os.Args[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "gamePk must be numeric: %v\n", err)
			os.Exit(2)
		}
		gamePk = n
	}
//
	c := mlb.New()
	ls, err := c.Linescore(context.Background(), gamePk, mlb.LinescoreQuery{})
	if err != nil {
		fmt.Fprintln(os.Stderr, "linescore:", err)
		os.Exit(1)
	}
	fmt.Printf("inning %d %s  (%d-%d)\n",
		ls.CurrentInning, ls.InningState,
		ls.Teams.Away.Runs, ls.Teams.Home.Runs)
	for _, inn := range ls.Innings {
		fmt.Printf("  %3s:  away=%d  home=%d\n",
			inn.OrdinalNum, inn.Away.Runs, inn.Home.Runs)
	}
	fmt.Printf("pitcher: %s  batter: %s\n",
		ls.Defense.Pitcher.FullName, ls.Offense.Batter.FullName)
}
