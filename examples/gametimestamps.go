// Copyright (c) 2026 John Dewey
//
// SPDX-License-Identifier: MIT
//
// GameTimestamps prints the live-feed timestamps for a game. Run with:
//
//	go run ./examples/gametimestamps [GAME_PK]   # default: 745455
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
	ts, err := c.GameTimestamps(context.Background(), gamePk)
	if err != nil {
		fmt.Fprintln(os.Stderr, "gameTimestamps:", err)
		os.Exit(1)
	}
	for _, t := range ts {
		fmt.Println(t)
	}
}
