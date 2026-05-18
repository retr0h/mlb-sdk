// Copyright (c) 2026 John Dewey
//
// SPDX-License-Identifier: MIT
//
// PlayByPlay prints every play in a game and counts grounded-into-double-plays
// (the v1 endpoint). Run with:
//
//	go run ./examples/playbyplay <gamePk>
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
		fmt.Fprintln(os.Stderr, "usage: playbyplay <gamePk>")
		os.Exit(2)
	}
	gamePk, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "gamePk must be an integer")
		os.Exit(2)
	}
//
	c := mlb.New()
	plays, err := c.PlayByPlay(context.Background(), gamePk)
	if err != nil {
		fmt.Fprintln(os.Stderr, "playByPlay:", err)
		os.Exit(1)
	}
	dps := 0
	for _, p := range plays {
		if p.IsDoublePlay() {
			dps++
		}
	}
	fmt.Printf("%d plays, %d grounded-into-double-plays\n", len(plays), dps)
}
