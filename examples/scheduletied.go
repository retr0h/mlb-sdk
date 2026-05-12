// Copyright (c) 2026 John Dewey
//
// SPDX-License-Identifier: MIT

// ScheduleTied prints tied/suspended games for a season. Run with:
//
//	go run ./examples/scheduletied [SEASON]   # default: 2024
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
	games, err := c.ScheduleTied(context.Background(), mlb.ScheduleTiedQuery{
		Season: year,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "scheduleTied:", err)
		os.Exit(1)
	}
	if len(games) == 0 {
		fmt.Println("no tied/suspended games")
		return
	}
	for _, g := range games {
		fmt.Printf("%-7d  %s @ %s  (%s)\n",
			g.GamePk, g.Away.Name, g.Home.Name, g.Status)
	}
}
