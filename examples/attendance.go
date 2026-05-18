// Copyright (c) 2026 John Dewey
//
// SPDX-License-Identifier: MIT
//
// Attendance prints a team's attendance record for a season. Run with:
//
//	go run ./examples/attendance [TEAM_ID] [SEASON]   # default: 119 2024
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
	teamID, year := 119, 2024
	if len(os.Args) > 1 {
		n, err := strconv.Atoi(os.Args[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "teamId must be numeric: %v\n", err)
			os.Exit(2)
		}
		teamID = n
	}
	if len(os.Args) > 2 {
		n, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Fprintf(os.Stderr, "season must be numeric: %v\n", err)
			os.Exit(2)
		}
		year = n
	}
//
	c := mlb.New()
	a, err := c.Attendance(context.Background(), mlb.AttendanceQuery{
		TeamID: teamID,
		Season: year,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "attendance:", err)
		os.Exit(1)
	}
	for _, r := range a.Records {
		fmt.Printf("%s (%s) — %s\n", r.Team.Name, r.Year, r.GameType.Description)
		fmt.Printf("  total:    %d\n", r.AttendanceTotal)
		fmt.Printf("  avg home: %d\n", r.AttendanceAverageHome)
		fmt.Printf("  avg away: %d\n", r.AttendanceAverageAway)
		fmt.Printf("  high:     %d on %s (gamePk %d)\n",
			r.AttendanceHigh,
			r.AttendanceHighDate.Format("2006-01-02"),
			r.AttendanceHighGame.GamePk)
	}
	fmt.Printf("\naggregate total: %d\n", a.AggregateTotals.AttendanceTotal)
}
