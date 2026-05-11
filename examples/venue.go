// Copyright (c) 2026 John Dewey
//
// SPDX-License-Identifier: MIT

// Venue prints a hydrated venue lookup. Run with:
//
//	go run ./examples/venue [VENUE_ID]   # default: 22 (Dodger Stadium)
package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/retr0h/mlb-sdk/pkg/mlb"
)

func main() {
	venueID := 22
	if len(os.Args) > 1 {
		n, err := strconv.Atoi(os.Args[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "venue id must be numeric: %v\n", err)
			os.Exit(2)
		}
		venueID = n
	}

	c := mlb.New()
	v, err := c.Venue(context.Background(), venueID, mlb.VenueQuery{
		Hydrate: "location,fieldInfo,timezone",
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "venue:", err)
		os.Exit(1)
	}
	fmt.Printf("%d  %s (%s)\n", v.ID, v.Name, v.Season)
	if v.Location.City != "" {
		fmt.Printf("  %s, %s %s  (%s)\n",
			v.Location.City, v.Location.StateAbbrev,
			v.Location.PostalCode, v.Location.Country)
	}
	if v.FieldInfo.Capacity > 0 {
		fmt.Printf("  capacity %d, %s, roof: %s\n",
			v.FieldInfo.Capacity, v.FieldInfo.TurfType, v.FieldInfo.RoofType)
	}
	if v.TimeZone.ID != "" {
		fmt.Printf("  time zone %s (UTC%+d)\n", v.TimeZone.ID, v.TimeZone.Offset)
	}
}
