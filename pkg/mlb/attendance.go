// Copyright (c) 2026 John Dewey
//
// SPDX-License-Identifier: MIT

package mlb

import (
	"context"
	"fmt"

	"github.com/retr0h/mlb-sdk/internal/gen"
)

const attendanceDateFmt = "2006-01-02"

// Attendance fetches attendance records. One of q.TeamID, q.LeagueID, or
// q.LeagueListID must be set; the MLB API rejects the call otherwise
// (toddrob99: required_params=[["teamId"], ["leagueId"], ["leagueListId"]]).
//
// Example:
//
//	a, _ := c.Attendance(ctx, mlb.AttendanceQuery{TeamID: 119, Season: 2024})
//	fmt.Println("total:", a.AggregateTotals.AttendanceTotal)
//	for _, r := range a.Records {
//	    fmt.Println(r.Team.Name, r.Year, r.AttendanceAverageHome)
//	}
func (c *Client) Attendance(ctx context.Context, q AttendanceQuery) (*Attendance, error) {
	if q.TeamID == 0 && q.LeagueID == 0 && q.LeagueListID == "" {
		return nil, fmt.Errorf(
			"mlb: attendance: %w: one of TeamID, LeagueID, or LeagueListID is required",
			ErrInvalidQuery,
		)
	}

	params := &gen.GetAttendanceParams{}
	if q.TeamID != 0 {
		params.TeamId = ptr(q.TeamID)
	}
	if q.LeagueID != 0 {
		params.LeagueId = ptr(q.LeagueID)
	}
	if q.LeagueListID != "" {
		params.LeagueListId = ptr(q.LeagueListID)
	}
	if q.Season != 0 {
		params.Season = ptr(q.Season)
	}
	if !q.On.IsZero() {
		params.Date = ptr(q.On.Format(attendanceDateFmt))
	}
	if q.GameType != "" {
		params.GameType = ptr(q.GameType)
	}
	if q.Fields != "" {
		params.Fields = ptr(q.Fields)
	}

	resp, err := c.raw.GetAttendanceWithResponse(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("mlb: attendance: %w", err)
	}
	if resp.StatusCode() == 404 {
		return nil, ErrNotFound
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("mlb: attendance: unexpected status %d", resp.StatusCode())
	}
	return attendanceFromGen(resp.JSON200), nil
}

func attendanceFromGen(r *gen.AttendanceResponse) *Attendance {
	out := &Attendance{}
	if r == nil {
		return out
	}
	if r.Records != nil {
		out.Records = make([]AttendanceRecord, 0, len(*r.Records))
		for _, rec := range *r.Records {
			out.Records = append(out.Records, attendanceRecordFromGen(rec))
		}
	}
	if r.AggregateTotals != nil {
		out.AggregateTotals = attendanceAggregateTotalsFromGen(*r.AggregateTotals)
	}
	return out
}

func attendanceRecordFromGen(r gen.AttendanceRecord) AttendanceRecord {
	out := AttendanceRecord{}
	if r.OpeningsTotal != nil {
		out.OpeningsTotal = *r.OpeningsTotal
	}
	if r.OpeningsTotalAway != nil {
		out.OpeningsTotalAway = *r.OpeningsTotalAway
	}
	if r.OpeningsTotalHome != nil {
		out.OpeningsTotalHome = *r.OpeningsTotalHome
	}
	if r.OpeningsTotalLost != nil {
		out.OpeningsTotalLost = *r.OpeningsTotalLost
	}
	if r.GamesTotal != nil {
		out.GamesTotal = *r.GamesTotal
	}
	if r.GamesAwayTotal != nil {
		out.GamesAwayTotal = *r.GamesAwayTotal
	}
	if r.GamesHomeTotal != nil {
		out.GamesHomeTotal = *r.GamesHomeTotal
	}
	if r.Year != nil {
		out.Year = *r.Year
	}
	if r.AttendanceAverageAway != nil {
		out.AttendanceAverageAway = *r.AttendanceAverageAway
	}
	if r.AttendanceAverageHome != nil {
		out.AttendanceAverageHome = *r.AttendanceAverageHome
	}
	if r.AttendanceAverageYtd != nil {
		out.AttendanceAverageYtd = *r.AttendanceAverageYtd
	}
	if r.AttendanceHigh != nil {
		out.AttendanceHigh = *r.AttendanceHigh
	}
	if r.AttendanceHighDate != nil {
		out.AttendanceHighDate = *r.AttendanceHighDate
	}
	if r.AttendanceHighGame != nil {
		out.AttendanceHighGame = attendanceGameRefFromGen(*r.AttendanceHighGame)
	}
	if r.AttendanceLow != nil {
		out.AttendanceLow = *r.AttendanceLow
	}
	if r.AttendanceLowDate != nil {
		out.AttendanceLowDate = *r.AttendanceLowDate
	}
	if r.AttendanceLowGame != nil {
		out.AttendanceLowGame = attendanceGameRefFromGen(*r.AttendanceLowGame)
	}
	if r.AttendanceOpeningAverage != nil {
		out.AttendanceOpeningAverage = *r.AttendanceOpeningAverage
	}
	if r.AttendanceTotal != nil {
		out.AttendanceTotal = *r.AttendanceTotal
	}
	if r.AttendanceTotalAway != nil {
		out.AttendanceTotalAway = *r.AttendanceTotalAway
	}
	if r.AttendanceTotalHome != nil {
		out.AttendanceTotalHome = *r.AttendanceTotalHome
	}
	if r.GameType != nil {
		out.GameType = gameTypeRefFromGen(*r.GameType)
	}
	if r.Team != nil {
		if r.Team.Id != nil {
			out.Team.ID = TeamID(*r.Team.Id)
		}
		if r.Team.Name != nil {
			out.Team.Name = *r.Team.Name
		}
	}
	return out
}

func attendanceGameRefFromGen(g gen.AttendanceGameRef) AttendanceGameRef {
	out := AttendanceGameRef{}
	if g.GamePk != nil {
		out.GamePk = *g.GamePk
	}
	if g.Link != nil {
		out.Link = *g.Link
	}
	if g.Content != nil {
		out.Content = refFromGen(g.Content)
	}
	if g.DayNight != nil {
		out.DayNight = *g.DayNight
	}
	return out
}

func gameTypeRefFromGen(g gen.GameTypeRef) GameTypeRef {
	out := GameTypeRef{}
	if g.Id != nil {
		out.ID = *g.Id
	}
	if g.Description != nil {
		out.Description = *g.Description
	}
	return out
}

func attendanceAggregateTotalsFromGen(a gen.AttendanceAggregateTotals) AttendanceAggregateTotals {
	out := AttendanceAggregateTotals{}
	if a.OpeningsTotalAway != nil {
		out.OpeningsTotalAway = *a.OpeningsTotalAway
	}
	if a.OpeningsTotalHome != nil {
		out.OpeningsTotalHome = *a.OpeningsTotalHome
	}
	if a.OpeningsTotalLost != nil {
		out.OpeningsTotalLost = *a.OpeningsTotalLost
	}
	if a.OpeningsTotalYtd != nil {
		out.OpeningsTotalYtd = *a.OpeningsTotalYtd
	}
	if a.AttendanceAverageAway != nil {
		out.AttendanceAverageAway = *a.AttendanceAverageAway
	}
	if a.AttendanceAverageHome != nil {
		out.AttendanceAverageHome = *a.AttendanceAverageHome
	}
	if a.AttendanceAverageYtd != nil {
		out.AttendanceAverageYtd = *a.AttendanceAverageYtd
	}
	if a.AttendanceHigh != nil {
		out.AttendanceHigh = *a.AttendanceHigh
	}
	if a.AttendanceHighDate != nil {
		out.AttendanceHighDate = *a.AttendanceHighDate
	}
	if a.AttendanceTotal != nil {
		out.AttendanceTotal = *a.AttendanceTotal
	}
	if a.AttendanceTotalAway != nil {
		out.AttendanceTotalAway = *a.AttendanceTotalAway
	}
	if a.AttendanceTotalHome != nil {
		out.AttendanceTotalHome = *a.AttendanceTotalHome
	}
	return out
}
