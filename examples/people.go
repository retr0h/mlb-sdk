// Copyright (c) 2026 John Dewey
//
// SPDX-License-Identifier: MIT

// People fetches multiple MLB players by id. Run with:
//
//	go run ./examples/people [IDS]   # default: 660271,545361
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/retr0h/mlb-sdk/pkg/mlb"
)

func main() {
	ids := "660271,545361"
	if len(os.Args) > 1 {
		ids = os.Args[1]
	}

	c := mlb.New()
	pp, err := c.People(context.Background(), mlb.PeopleQuery{PersonIDs: ids})
	if err != nil {
		fmt.Fprintln(os.Stderr, "people:", err)
		os.Exit(1)
	}
	for _, p := range pp {
		fmt.Printf("%-25s [%s] %s %s  active=%v\n",
			p.FullName, p.PrimaryPosition.Abbreviation,
			p.BatSide.Code, p.PitchHand.Code, p.Active)
	}
}
