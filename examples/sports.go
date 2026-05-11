// Copyright (c) 2026 John Dewey
//
// SPDX-License-Identifier: MIT

// Sports prints every sport the MLB API tracks. Run with:
//
//	go run ./examples/sports
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/retr0h/mlb-sdk/pkg/mlb"
)

func main() {
	c := mlb.New()
	s, err := c.Sports(context.Background(), mlb.SportsQuery{})
	if err != nil {
		fmt.Fprintln(os.Stderr, "sports:", err)
		os.Exit(1)
	}
	for _, sport := range s.Sports {
		fmt.Printf(
			"%3d  %-5s  %-30s  active=%v\n",
			sport.ID, sport.Abbreviation, sport.Name, sport.ActiveStatus,
		)
	}
}
