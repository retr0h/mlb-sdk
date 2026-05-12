// Copyright (c) 2026 John Dewey
//
// SPDX-License-Identifier: MIT

// Umpires prints the MLB umpire roster. Run with:
//
//	go run ./examples/umpires
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/retr0h/mlb-sdk/pkg/mlb"
)

func main() {
	c := mlb.New()
	s, err := c.Umpires(context.Background(), mlb.UmpiresQuery{SportID: 1})
	if err != nil {
		fmt.Fprintln(os.Stderr, "umpires:", err)
		os.Exit(1)
	}
	for _, e := range s.Roster {
		fmt.Printf("#%-3s %-30s %s\n", e.JerseyNumber, e.Person.FullName, e.Job)
	}
}
