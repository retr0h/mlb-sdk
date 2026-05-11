// Copyright (c) 2026 John Dewey
//
// SPDX-License-Identifier: MIT

// Season prints the calendar windows for a single MLB season. Run with:
//
//	go run ./examples/season [YEAR]   # default: 2024
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/retr0h/mlb-sdk/pkg/mlb"
)

const dateFmt = "2006-01-02"

func main() {
	year := "2024"
	if len(os.Args) > 1 {
		year = os.Args[1]
	}

	c := mlb.New()
	s, err := c.Season(context.Background(), year, mlb.SeasonQuery{SportID: 1})
	if err != nil {
		fmt.Fprintln(os.Stderr, "season:", err)
		os.Exit(1)
	}
	fmt.Printf("season %s (wildcard=%v)\n", s.SeasonID, s.HasWildcard)
	fmt.Printf("  regular     %s → %s\n",
		s.RegularSeasonStartDate.Format(dateFmt),
		s.RegularSeasonEndDate.Format(dateFmt))
	if !s.AllStarDate.IsZero() {
		fmt.Printf("  all-star    %s\n", s.AllStarDate.Format(dateFmt))
	}
	fmt.Printf("  post-season %s → %s\n",
		s.PostSeasonStartDate.Format(dateFmt),
		s.PostSeasonEndDate.Format(dateFmt))
}
