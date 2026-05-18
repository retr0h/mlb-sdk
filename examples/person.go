// Copyright (c) 2026 John Dewey
//
// SPDX-License-Identifier: MIT
//
// Person prints a single MLB person's bio. Run with:
//
//	go run ./examples/person [PERSON_ID]   # default: 660271 (Ohtani)
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
	id := 660271
	if len(os.Args) > 1 {
		n, err := strconv.Atoi(os.Args[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "personId must be numeric: %v\n", err)
			os.Exit(2)
		}
		id = n
	}
//
	c := mlb.New()
	p, err := c.Person(context.Background(), id, mlb.PersonQuery{})
	if err != nil {
		fmt.Fprintln(os.Stderr, "person:", err)
		os.Exit(1)
	}
	fmt.Printf("%s (#%s) — %s %s\n",
		p.FullName, p.PrimaryNumber,
		p.PrimaryPosition.Name, p.PrimaryPosition.Abbreviation)
	fmt.Printf("  born: %s, %s  height: %s  weight: %d\n",
		p.BirthCity, p.BirthCountry, p.Height, p.Weight)
	fmt.Printf("  bats: %s  throws: %s\n",
		p.BatSide.Description, p.PitchHand.Description)
	fmt.Printf("  debut: %s  active: %v\n", p.MlbDebutDate, p.Active)
}
