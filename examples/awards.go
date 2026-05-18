// Copyright (c) 2026 John Dewey
//
// SPDX-License-Identifier: MIT
//
// Awards prints recipients of an MLB award. Run with:
//
//	go run ./examples/awards [AWARD_ID]   # default: MLBHOF
package main
//
import (
	"context"
	"fmt"
	"os"
//
	"github.com/retr0h/mlb-sdk/pkg/mlb"
)
//
func main() {
	awardID := "MLBHOF"
	if len(os.Args) > 1 {
		awardID = os.Args[1]
	}
//
	c := mlb.New()
	a, err := c.AwardRecipients(context.Background(), awardID, mlb.AwardRecipientsQuery{})
	if err != nil {
		fmt.Fprintln(os.Stderr, "awardRecipients:", err)
		os.Exit(1)
	}
	for _, r := range a.Recipients {
		notes := ""
		if r.Notes != "" {
			notes = "  (" + r.Notes + ")"
		}
		fmt.Printf("%s  %-25s [%s]%s\n",
			r.Date, r.Player.NameFirstLast, r.Player.PrimaryPosition.Abbreviation, notes)
	}
}
