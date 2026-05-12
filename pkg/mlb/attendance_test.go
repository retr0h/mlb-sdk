// Copyright (c) 2026 John Dewey
//
// SPDX-License-Identifier: MIT

package mlb

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

const attendanceHappyBody = `{
  "copyright": "Copyright 2026 MLB Advanced Media, L.P.",
  "records": [{
    "openingsTotal": 161,
    "openingsTotalAway": 80, "openingsTotalHome": 81, "openingsTotalLost": 0,
    "gamesTotal": 163, "gamesAwayTotal": 82, "gamesHomeTotal": 81,
    "year": "2024",
    "attendanceAverageAway": 36253,
    "attendanceAverageHome": 48657,
    "attendanceAverageYtd":  42494,
    "attendanceHigh": 54070,
    "attendanceHighDate": "2024-07-24T00:00:00Z",
    "attendanceHighGame": {
      "gamePk": 746119,
      "link": "/api/v1.1/game/746119/feed/live",
      "content": {"id": 0, "link": "/api/v1/game/746119/content"},
      "dayNight": "day"
    },
    "attendanceLow": 15928,
    "attendanceLowDate": "2024-03-21T00:00:00Z",
    "attendanceLowGame": {
      "gamePk": 746175,
      "link": "/api/v1.1/game/746175/feed/live",
      "content": {"id": 0, "link": "/api/v1/game/746175/content"},
      "dayNight": "day"
    },
    "attendanceOpeningAverage": 48657,
    "attendanceTotal":      6841502,
    "attendanceTotalAway":  2900251,
    "attendanceTotalHome":  3941251,
    "gameType": {"id": "R", "description": "Regular Season"},
    "team": {"id": 119, "name": "Los Angeles Dodgers"}
  }],
  "aggregateTotals": {
    "openingsTotalAway": 80, "openingsTotalHome": 81,
    "openingsTotalLost": 0, "openingsTotalYtd": 0,
    "attendanceAverageAway": 36253,
    "attendanceAverageHome": 48657,
    "attendanceAverageYtd": 42494,
    "attendanceHigh": 54070,
    "attendanceHighDate": "2024-07-24T00:00:00Z",
    "attendanceTotal": 6841502,
    "attendanceTotalAway": 2900251,
    "attendanceTotalHome": 3941251
  }
}`

func TestAttendanceFromGen(t *testing.T) {
	// Cover the defensive nil-response branch — the wrapper guards against
	// nil in production, but the helper is exposed to the package and
	// should tolerate it.
	got := attendanceFromGen(nil)
	if got == nil {
		t.Fatal("attendanceFromGen(nil) = nil, want empty *Attendance")
	}
	if len(got.Records) != 0 || got.AggregateTotals != (AttendanceAggregateTotals{}) {
		t.Errorf("attendanceFromGen(nil) = %+v, want zero", got)
	}
}

func TestClient_Attendance(t *testing.T) {
	on := time.Date(2024, 7, 24, 0, 0, 0, 0, time.UTC)
	cases := []struct {
		name         string
		query        AttendanceQuery
		respStatus   int
		respBody     string
		wantPath     string
		wantQuery    url.Values
		wantErr      string
		wantIs       error
		wantRecs     int
		wantHydrated bool
	}{
		{
			name: "happy path: every filter set",
			query: AttendanceQuery{
				TeamID:       119,
				LeagueID:     104,
				LeagueListID: "milb_all",
				Season:       2024,
				On:           on,
				GameType:     "R",
				Fields:       "records,team,attendanceTotal",
			},
			respStatus: 200,
			respBody:   attendanceHappyBody,
			wantPath:   "/api/v1/attendance",
			wantQuery: url.Values{
				"teamId":       {"119"},
				"leagueId":     {"104"},
				"leagueListId": {"milb_all"},
				"season":       {"2024"},
				"date":         {"2024-07-24"},
				"gameType":     {"R"},
				"fields":       {"records,team,attendanceTotal"},
			},
			wantRecs:     1,
			wantHydrated: true,
		},
		{
			name:       "filter with LeagueID only",
			query:      AttendanceQuery{LeagueID: 104},
			respStatus: 200,
			respBody:   `{"records":[{"year":"2024"}]}`,
			wantPath:   "/api/v1/attendance",
			wantQuery:  url.Values{"leagueId": {"104"}},
			wantRecs:   1,
		},
		{
			name:       "filter with LeagueListID only",
			query:      AttendanceQuery{LeagueListID: "milb_all"},
			respStatus: 200,
			respBody:   `{"records":[]}`,
			wantPath:   "/api/v1/attendance",
			wantQuery:  url.Values{"leagueListId": {"milb_all"}},
			wantRecs:   0,
		},
		{
			name:    "missing required combo rejected before HTTP",
			query:   AttendanceQuery{Season: 2024},
			wantErr: "one of TeamID, LeagueID, or LeagueListID",
			wantIs:  ErrInvalidQuery,
		},
		{
			name:       "200 with no records yields empty slice",
			query:      AttendanceQuery{TeamID: 119},
			respStatus: 200,
			respBody:   `{}`,
			wantRecs:   0,
		},
		{
			name:       "404 returns ErrNotFound",
			query:      AttendanceQuery{TeamID: 119},
			respStatus: 404,
			respBody:   `{}`,
			wantIs:     ErrNotFound,
			wantErr:    "not found",
		},
		{
			name:       "5xx is wrapped",
			query:      AttendanceQuery{TeamID: 119},
			respStatus: 500,
			respBody:   `oops`,
			wantErr:    "unexpected status 500",
		},
		{
			name:       "malformed JSON is wrapped",
			query:      AttendanceQuery{TeamID: 119},
			respStatus: 200,
			respBody:   `not json`,
			wantErr:    "attendance",
		},
		{
			name:       "network failure is wrapped",
			query:      AttendanceQuery{TeamID: 119},
			respStatus: 0,
			wantErr:    "attendance",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var seenPath string
			var seenQuery url.Values
			srv := httptest.NewServer(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					seenPath = r.URL.Path
					seenQuery = r.URL.Query()
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(c.respStatus)
					_, _ = w.Write([]byte(c.respBody))
				}),
			)
			urlStr := srv.URL
			if c.respStatus == 0 {
				srv.Close()
			} else {
				defer srv.Close()
			}

			client := New(WithBaseURL(urlStr))
			a, err := client.Attendance(context.Background(), c.query)

			if c.wantErr != "" {
				if err == nil {
					t.Fatalf("expected error containing %q, got nil", c.wantErr)
				}
				if !strings.Contains(err.Error(), c.wantErr) {
					t.Errorf("err = %v, want substring %q", err, c.wantErr)
				}
				if c.wantIs != nil && !errors.Is(err, c.wantIs) {
					t.Errorf("errors.Is(err, %v) = false, want true", c.wantIs)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if a == nil {
				t.Fatal("expected non-nil Attendance")
			}
			if c.wantPath != "" && seenPath != c.wantPath {
				t.Errorf("path = %q, want %q", seenPath, c.wantPath)
			}
			if c.wantQuery != nil {
				for k, want := range c.wantQuery {
					if got := seenQuery.Get(k); got != want[0] {
						t.Errorf("query[%q] = %q, want %q", k, got, want[0])
					}
				}
			}
			if got := len(a.Records); got != c.wantRecs {
				t.Errorf("len(Records) = %d, want %d", got, c.wantRecs)
			}
			if !c.wantHydrated {
				return
			}
			r := a.Records[0]
			if r.OpeningsTotal != 161 || r.OpeningsTotalAway != 80 ||
				r.OpeningsTotalHome != 81 || r.OpeningsTotalLost != 0 ||
				r.GamesTotal != 163 || r.GamesAwayTotal != 82 ||
				r.GamesHomeTotal != 81 || r.Year != "2024" {
				t.Errorf("records[0] basics = %+v", r)
			}
			if r.AttendanceAverageAway != 36253 || r.AttendanceAverageHome != 48657 ||
				r.AttendanceAverageYtd != 42494 || r.AttendanceHigh != 54070 ||
				r.AttendanceLow != 15928 || r.AttendanceOpeningAverage != 48657 ||
				r.AttendanceTotal != 6841502 || r.AttendanceTotalAway != 2900251 ||
				r.AttendanceTotalHome != 3941251 {
				t.Errorf("records[0] aggregate fields = %+v", r)
			}
			if r.AttendanceHighDate.Format("2006-01-02") != "2024-07-24" {
				t.Errorf("AttendanceHighDate = %v", r.AttendanceHighDate)
			}
			if r.AttendanceLowDate.Format("2006-01-02") != "2024-03-21" {
				t.Errorf("AttendanceLowDate = %v", r.AttendanceLowDate)
			}
			if r.AttendanceHighGame.GamePk != 746119 ||
				r.AttendanceHighGame.Link != "/api/v1.1/game/746119/feed/live" ||
				r.AttendanceHighGame.Content.Link != "/api/v1/game/746119/content" ||
				r.AttendanceHighGame.DayNight != "day" {
				t.Errorf("AttendanceHighGame = %+v", r.AttendanceHighGame)
			}
			if r.AttendanceLowGame.GamePk != 746175 ||
				r.AttendanceLowGame.DayNight != "day" {
				t.Errorf("AttendanceLowGame = %+v", r.AttendanceLowGame)
			}
			if r.GameType.ID != "R" || r.GameType.Description != "Regular Season" {
				t.Errorf("GameType = %+v", r.GameType)
			}
			if r.Team.ID != LAD || r.Team.Name != "Los Angeles Dodgers" {
				t.Errorf("Team = %+v", r.Team)
			}
			ag := a.AggregateTotals
			if ag.OpeningsTotalAway != 80 || ag.OpeningsTotalHome != 81 ||
				ag.OpeningsTotalLost != 0 || ag.OpeningsTotalYtd != 0 ||
				ag.AttendanceAverageAway != 36253 || ag.AttendanceAverageHome != 48657 ||
				ag.AttendanceAverageYtd != 42494 || ag.AttendanceHigh != 54070 ||
				ag.AttendanceTotal != 6841502 || ag.AttendanceTotalAway != 2900251 ||
				ag.AttendanceTotalHome != 3941251 {
				t.Errorf("AggregateTotals = %+v", ag)
			}
			if ag.AttendanceHighDate.Format("2006-01-02") != "2024-07-24" {
				t.Errorf("AggregateTotals.AttendanceHighDate = %v", ag.AttendanceHighDate)
			}
		})
	}
}
