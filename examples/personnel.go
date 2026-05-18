// Copyright (c) 2026 John Dewey
//
// SPDX-License-Identifier: MIT
//
// Personnel prints a team's front-office personnel. Run with:
//
//	go run ./examples/personnel [TEAM_ID]   # default: 119
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
	teamID := 119
	if len(os.Args) > 1 {
		n, err := strconv.Atoi(os.Args[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "teamId must be numeric: %v\n", err)
			os.Exit(2)
		}
		teamID = n
	}
//
	c := mlb.New()
	s, err := c.Personnel(context.Background(), teamID, mlb.PersonnelQuery{})
	if err != nil {
		fmt.Fprintln(os.Stderr, "personnel:", err)
		os.Exit(1)
	}
	for _, e := range s.Roster {
		fmt.Printf("%-30s %s\n", e.Person.FullName, e.Title)
	}
}
