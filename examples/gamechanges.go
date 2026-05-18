// Copyright (c) 2026 John Dewey
//
// SPDX-License-Identifier: MIT
//
// GameChanges prints recently updated games. Run with:
//
//	go run ./examples/gamechanges
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
	since := time.Now().Add(-24 * time.Hour).Format(time.RFC3339)
	c := mlb.New()
	gc, err := c.GameChanges(context.Background(), mlb.GameChangesQuery{
		UpdatedSince: since,
		SportID:      1,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "gameChanges:", err)
		os.Exit(1)
	}
	fmt.Printf("changed games (last 24h): %d\n", gc.TotalGames)
	for _, d := range gc.Dates {
		for _, g := range d.Games {
			fmt.Printf("  %s  gamePk=%d  %s\n", d.Date, g.GamePk, g.Status)
		}
	}
}
