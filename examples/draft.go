// Copyright (c) 2026 John Dewey
//
// SPDX-License-Identifier: MIT
//
// Draft prints top draft picks for a year. Run with:
//
//	go run ./examples/draft [YEAR]   # default: 2024
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
	year := 2024
	if len(os.Args) > 1 {
		n, err := strconv.Atoi(os.Args[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "year must be numeric: %v\n", err)
			os.Exit(2)
		}
		year = n
	}
//
	c := mlb.New()
	d, err := c.Draft(context.Background(), year, mlb.DraftQuery{Round: "1"})
	if err != nil {
		fmt.Fprintln(os.Stderr, "draft:", err)
		os.Exit(1)
	}
	fmt.Printf("Draft %d\n", d.DraftYear)
	for _, r := range d.Rounds {
		for _, p := range r.Picks {
			fmt.Printf("  #%-3d %-25s [%s]  %s (%s)\n",
				p.PickNumber, p.Person.FullName,
				p.Person.PrimaryPosition.Abbreviation,
				p.Team.Name, p.School.Name)
		}
	}
}
