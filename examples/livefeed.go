// Copyright (c) 2026 John Dewey
//
// SPDX-License-Identifier: MIT
//
// LiveFeed prints every play in a game from the v1.1 /feed/live endpoint —
// the same source MLB Gameday consumes. Run with:
//
//	go run ./examples/livefeed <gamePk>
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
		fmt.Fprintln(os.Stderr, "usage: livefeed <gamePk>")
		os.Exit(2)
	}
	gamePk, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "gamePk must be an integer")
		os.Exit(2)
	}
//
	c := mlb.New()
	plays, err := c.LiveFeed(context.Background(), gamePk)
	if err != nil {
		fmt.Fprintln(os.Stderr, "liveFeed:", err)
		os.Exit(1)
	}
	for _, p := range plays {
		fmt.Printf("[%s %d] %s — %s\n", p.HalfInning, p.Inning, p.EventType, p.Description)
	}
}
