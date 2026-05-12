// Copyright (c) 2026 John Dewey
//
// SPDX-License-Identifier: MIT

package mlb

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var statsBody = `{"stats":[{"type":{"displayName":"season"},"group":{"displayName":"hitting"},"splits":[]}]}`

func TestClient_Stats(t *testing.T) {
	cases := []struct {
		name       string
		query      StatsQuery
		respStatus int
		respBody   string
		wantErr    string
		wantIs     error
	}{
		{name: "happy", query: StatsQuery{
			Stats: "season", Group: "hitting", Season: 2024, SportIDs: "1",
			GameType: "R", PlayerPool: "all", Position: "OF", TeamID: 119,
			LeagueID: 104, PersonID: 660271, Limit: 10, Offset: 5,
			SortStat: "homeRuns", Order: "desc", Metrics: "exitVelocity",
			StartDate: "2024-04-01", EndDate: "2024-09-30",
			Hydrate: "person", Fields: "stats",
		}, respStatus: 200, respBody: statsBody},
		{
			name: "missing required", query: StatsQuery{Stats: "season"},
			wantErr: "Stats and Group are both required", wantIs: ErrInvalidQuery,
		},
		{
			name: "404", query: StatsQuery{Stats: "season", Group: "hitting"},
			respStatus: 404, respBody: `{}`, wantIs: ErrNotFound, wantErr: "not found",
		},
		{
			name: "5xx", query: StatsQuery{Stats: "season", Group: "hitting"},
			respStatus: 500, respBody: `oops`, wantErr: "unexpected status 500",
		},
		{
			name: "bad json", query: StatsQuery{Stats: "season", Group: "hitting"},
			respStatus: 200, respBody: `x`, wantErr: "stats",
		},
		{
			name: "network", query: StatsQuery{Stats: "season", Group: "hitting"},
			respStatus: 0, wantErr: "stats",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			srv := httptest.NewServer(
				http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(c.respStatus)
					_, _ = w.Write([]byte(c.respBody))
				}),
			)
			u := srv.URL
			if c.respStatus == 0 {
				srv.Close()
			} else {
				defer srv.Close()
			}
			_, err := New(WithBaseURL(u)).Stats(context.Background(), c.query)
			if c.wantErr != "" {
				if err == nil {
					t.Fatalf("expected %q", c.wantErr)
				}
				if !strings.Contains(err.Error(), c.wantErr) {
					t.Errorf("err = %v", err)
				}
				if c.wantIs != nil && !errors.Is(err, c.wantIs) {
					t.Errorf("errors.Is = false")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected: %v", err)
			}
		})
	}
}

func TestClient_SchedulePostseasonSeries(t *testing.T) {
	body := `{"dates":[{"date":"2024-10-01","games":[{"gamePk":1}]}]}`
	cases := []struct {
		name       string
		query      SchedulePostseasonQuery
		respStatus int
		respBody   string
		wantErr    string
		wantIs     error
		wantLen    int
	}{
		{
			name: "happy", query: SchedulePostseasonQuery{
				Season: 2024, GameTypes: "D,L,W", SeriesNumber: 1,
				TeamID: 119, SportID: 1, Fields: "dates",
			},
			respStatus: 200, respBody: body, wantLen: 1,
		},
		{name: "empty", respStatus: 200, respBody: `{}`, wantLen: 0},
		{name: "404", respStatus: 404, respBody: `{}`, wantIs: ErrNotFound, wantErr: "not found"},
		{name: "5xx", respStatus: 500, respBody: `oops`, wantErr: "unexpected status 500"},
		{name: "bad json", respStatus: 200, respBody: `x`, wantErr: "schedulePostseasonSeries"},
		{name: "network", respStatus: 0, wantErr: "schedulePostseasonSeries"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			srv := httptest.NewServer(
				http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(c.respStatus)
					_, _ = w.Write([]byte(c.respBody))
				}),
			)
			u := srv.URL
			if c.respStatus == 0 {
				srv.Close()
			} else {
				defer srv.Close()
			}
			games, err := New(
				WithBaseURL(u),
			).SchedulePostseasonSeries(context.Background(), c.query)
			if c.wantErr != "" {
				if err == nil {
					t.Fatalf("expected %q", c.wantErr)
				}
				if !strings.Contains(err.Error(), c.wantErr) {
					t.Errorf("err = %v", err)
				}
				if c.wantIs != nil && !errors.Is(err, c.wantIs) {
					t.Errorf("errors.Is = false")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected: %v", err)
			}
			if len(games) != c.wantLen {
				t.Errorf("len = %d, want %d", len(games), c.wantLen)
			}
		})
	}
}

func TestClient_AllStarBallot(t *testing.T) {
	cases := []struct {
		name       string
		leagueID   int
		query      AllStarBallotQuery
		respStatus int
		respBody   string
		wantErr    string
		wantIs     error
		wantLen    int
	}{
		{
			name: "happy", leagueID: 103, query: AllStarBallotQuery{Season: 2024, Fields: "people"},
			respStatus: 200, respBody: `{"people":[{"id":1}]}`, wantLen: 1,
		},
		{
			name: "missing Season", leagueID: 103, query: AllStarBallotQuery{},
			wantErr: "Season is required", wantIs: ErrInvalidQuery,
		},
		{
			name: "empty", leagueID: 103, query: AllStarBallotQuery{Season: 2024},
			respStatus: 200, respBody: `{}`, wantLen: 0,
		},
		{
			name: "404", leagueID: 103, query: AllStarBallotQuery{Season: 2024},
			respStatus: 404, respBody: `{}`, wantIs: ErrNotFound, wantErr: "not found",
		},
		{
			name: "5xx", leagueID: 103, query: AllStarBallotQuery{Season: 2024},
			respStatus: 500, respBody: `oops`, wantErr: "unexpected status 500",
		},
		{
			name: "bad json", leagueID: 103, query: AllStarBallotQuery{Season: 2024},
			respStatus: 200, respBody: `x`, wantErr: "allStarBallot",
		},
		{
			name: "network", leagueID: 103, query: AllStarBallotQuery{Season: 2024},
			respStatus: 0, wantErr: "allStarBallot",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			srv := httptest.NewServer(
				http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(c.respStatus)
					_, _ = w.Write([]byte(c.respBody))
				}),
			)
			u := srv.URL
			if c.respStatus == 0 {
				srv.Close()
			} else {
				defer srv.Close()
			}
			pp, err := New(WithBaseURL(u)).AllStarBallot(context.Background(), c.leagueID, c.query)
			if c.wantErr != "" {
				if err == nil {
					t.Fatalf("expected %q", c.wantErr)
				}
				if !strings.Contains(err.Error(), c.wantErr) {
					t.Errorf("err = %v", err)
				}
				if c.wantIs != nil && !errors.Is(err, c.wantIs) {
					t.Errorf("errors.Is = false")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected: %v", err)
			}
			if len(pp) != c.wantLen {
				t.Errorf("len = %d, want %d", len(pp), c.wantLen)
			}
		})
	}
}

func TestClient_AllStarFinalVote(t *testing.T) {
	cases := []struct {
		name       string
		leagueID   int
		query      AllStarBallotQuery
		respStatus int
		respBody   string
		wantErr    string
		wantIs     error
		wantLen    int
	}{
		{
			name: "happy", leagueID: 103, query: AllStarBallotQuery{Season: 2024, Fields: "people"},
			respStatus: 200, respBody: `{"people":[{"id":1}]}`, wantLen: 1,
		},
		{
			name: "missing Season", leagueID: 103, query: AllStarBallotQuery{},
			wantErr: "Season is required", wantIs: ErrInvalidQuery,
		},
		{
			name: "404", leagueID: 103, query: AllStarBallotQuery{Season: 2024},
			respStatus: 404, respBody: `{}`, wantIs: ErrNotFound, wantErr: "not found",
		},
		{
			name: "5xx", leagueID: 103, query: AllStarBallotQuery{Season: 2024},
			respStatus: 500, respBody: `oops`, wantErr: "unexpected status 500",
		},
		{
			name: "bad json", leagueID: 103, query: AllStarBallotQuery{Season: 2024},
			respStatus: 200, respBody: `x`, wantErr: "allStarFinalVote",
		},
		{
			name: "network", leagueID: 103, query: AllStarBallotQuery{Season: 2024},
			respStatus: 0, wantErr: "allStarFinalVote",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			srv := httptest.NewServer(
				http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(c.respStatus)
					_, _ = w.Write([]byte(c.respBody))
				}),
			)
			u := srv.URL
			if c.respStatus == 0 {
				srv.Close()
			} else {
				defer srv.Close()
			}
			pp, err := New(
				WithBaseURL(u),
			).AllStarFinalVote(context.Background(), c.leagueID, c.query)
			if c.wantErr != "" {
				if err == nil {
					t.Fatalf("expected %q", c.wantErr)
				}
				if !strings.Contains(err.Error(), c.wantErr) {
					t.Errorf("err = %v", err)
				}
				if c.wantIs != nil && !errors.Is(err, c.wantIs) {
					t.Errorf("errors.Is = false")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected: %v", err)
			}
			if len(pp) != c.wantLen {
				t.Errorf("len = %d, want %d", len(pp), c.wantLen)
			}
		})
	}
}

func TestClient_AllStarWriteIns(t *testing.T) {
	cases := []struct {
		name       string
		leagueID   int
		query      AllStarBallotQuery
		respStatus int
		respBody   string
		wantErr    string
		wantIs     error
		wantLen    int
	}{
		{
			name: "happy", leagueID: 103, query: AllStarBallotQuery{Season: 2024, Fields: "people"},
			respStatus: 200, respBody: `{"people":[{"id":1}]}`, wantLen: 1,
		},
		{
			name: "missing Season", leagueID: 103, query: AllStarBallotQuery{},
			wantErr: "Season is required", wantIs: ErrInvalidQuery,
		},
		{
			name: "404", leagueID: 103, query: AllStarBallotQuery{Season: 2024},
			respStatus: 404, respBody: `{}`, wantIs: ErrNotFound, wantErr: "not found",
		},
		{
			name: "5xx", leagueID: 103, query: AllStarBallotQuery{Season: 2024},
			respStatus: 500, respBody: `oops`, wantErr: "unexpected status 500",
		},
		{
			name: "bad json", leagueID: 103, query: AllStarBallotQuery{Season: 2024},
			respStatus: 200, respBody: `x`, wantErr: "allStarWriteIns",
		},
		{
			name: "network", leagueID: 103, query: AllStarBallotQuery{Season: 2024},
			respStatus: 0, wantErr: "allStarWriteIns",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			srv := httptest.NewServer(
				http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(c.respStatus)
					_, _ = w.Write([]byte(c.respBody))
				}),
			)
			u := srv.URL
			if c.respStatus == 0 {
				srv.Close()
			} else {
				defer srv.Close()
			}
			pp, err := New(
				WithBaseURL(u),
			).AllStarWriteIns(context.Background(), c.leagueID, c.query)
			if c.wantErr != "" {
				if err == nil {
					t.Fatalf("expected %q", c.wantErr)
				}
				if !strings.Contains(err.Error(), c.wantErr) {
					t.Errorf("err = %v", err)
				}
				if c.wantIs != nil && !errors.Is(err, c.wantIs) {
					t.Errorf("errors.Is = false")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected: %v", err)
			}
			if len(pp) != c.wantLen {
				t.Errorf("len = %d, want %d", len(pp), c.wantLen)
			}
		})
	}
}
