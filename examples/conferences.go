// Copyright (c) 2026 John Dewey
//
// SPDX-License-Identifier: MIT

// Conferences prints the conferences the MLB API tracks. Run with:
//
//	go run ./examples/conferences
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/retr0h/mlb-sdk/pkg/mlb"
)

func main() {
	c := mlb.New()
	co, err := c.Conferences(context.Background(), mlb.ConferencesQuery{})
	if err != nil {
		fmt.Fprintln(os.Stderr, "conferences:", err)
		os.Exit(1)
	}
	for _, conf := range co.Conferences {
		fmt.Printf("%3d  %-30s  %s\n", conf.ID, conf.Name, conf.Abbreviation)
	}
}
