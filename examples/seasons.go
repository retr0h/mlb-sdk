// Copyright (c) 2026 John Dewey
//
// SPDX-License-Identifier: MIT

// Seasons prints the regular-season window for a given MLB year. Run with:
//
//	go run ./examples/seasons [YEAR]   # default: 2024
package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/retr0h/mlb-sdk/pkg/mlb"
)

const dateFmt = "2006-01-02"

func main() {
	year := 2024
	if len(os.Args) > 1 {
		n, err := strconv.Atoi(os.Args[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "year must be numeric: %v\n", err)
			os.Exit(2)
		}
		year = n
	}

	c := mlb.New()
	s, err := c.Seasons(context.Background(), mlb.SeasonsQuery{
		SportID: 1,
		Season:  year,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "seasons:", err)
		os.Exit(1)
	}
	for _, season := range s.Seasons {
		fmt.Printf("season %s (wildcard=%v)\n", season.SeasonID, season.HasWildcard)
		fmt.Printf("  spring        %s → %s\n",
			season.SpringStartDate.Format(dateFmt),
			season.SpringEndDate.Format(dateFmt))
		fmt.Printf("  regular       %s → %s\n",
			season.RegularSeasonStartDate.Format(dateFmt),
			season.RegularSeasonEndDate.Format(dateFmt))
		if !season.AllStarDate.IsZero() {
			fmt.Printf("  all-star      %s\n", season.AllStarDate.Format(dateFmt))
		}
		fmt.Printf("  post-season   %s → %s\n",
			season.PostSeasonStartDate.Format(dateFmt),
			season.PostSeasonEndDate.Format(dateFmt))
	}
}
