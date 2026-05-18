// Copyright (c) 2026 John Dewey
//
// SPDX-License-Identifier: MIT
//
package mlb
//
import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)
//
func TestClient_TeamsAffiliates(t *testing.T) {
	cases := []struct {
		name       string
		query      TeamsAffiliatesQuery
		respStatus int
		respBody   string
		wantPath   string
		wantQuery  url.Values
		wantErr    string
		wantIs     error
		wantLen    int
	}{
		{
			name: "happy path",
			query: TeamsAffiliatesQuery{
				TeamIDs: "119", SportID: 1, Season: 2024,
				Hydrate: "league", Fields: "teams,id",
			},
			respStatus: 200,
			respBody:   `{"teams":[{"id":119,"name":"Dodgers"},{"id":260,"name":"OKC"}]}`,
			wantPath:   "/api/v1/teams/affiliates",
			wantQuery: url.Values{
				"teamIds": {"119"}, "sportId": {"1"}, "season": {"2024"},
				"hydrate": {"league"}, "fields": {"teams,id"},
			},
			wantLen: 2,
		},
		{
			name: "missing TeamIDs", query: TeamsAffiliatesQuery{},
			wantErr: "TeamIDs is required", wantIs: ErrInvalidQuery,
		},
		{
			name: "empty", query: TeamsAffiliatesQuery{TeamIDs: "119"},
			respStatus: 200, respBody: `{}`, wantLen: 0,
		},
		{
			name: "404", query: TeamsAffiliatesQuery{TeamIDs: "119"},
			respStatus: 404, respBody: `{}`, wantIs: ErrNotFound, wantErr: "not found",
		},
		{
			name: "5xx", query: TeamsAffiliatesQuery{TeamIDs: "119"},
			respStatus: 500, respBody: `oops`, wantErr: "unexpected status 500",
		},
		{
			name: "bad json", query: TeamsAffiliatesQuery{TeamIDs: "119"},
			respStatus: 200, respBody: `x`, wantErr: "teamsAffiliates",
		},
		{
			name: "network", query: TeamsAffiliatesQuery{TeamIDs: "119"},
			respStatus: 0, wantErr: "teamsAffiliates",
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
			ts, err := client.TeamsAffiliates(context.Background(), c.query)
			if c.wantErr != "" {
				if err == nil {
					t.Fatalf("expected error containing %q", c.wantErr)
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
			if len(ts.Teams) != c.wantLen {
				t.Errorf("len = %d, want %d", len(ts.Teams), c.wantLen)
			}
		})
	}
}
//
func TestClient_TeamsHistory(t *testing.T) {
	cases := []struct {
		name       string
		query      TeamsHistoryQuery
		respStatus int
		respBody   string
		wantPath   string
		wantQuery  url.Values
		wantErr    string
		wantIs     error
		wantLen    int
	}{
		{
			name: "happy path",
			query: TeamsHistoryQuery{
				TeamIDs: "119", StartSeason: 1960, EndSeason: 2024,
				Fields: "teams,id,name",
			},
			respStatus: 200,
			respBody:   `{"teams":[{"id":119,"name":"Dodgers","season":1962}]}`,
			wantPath:   "/api/v1/teams/history",
			wantQuery: url.Values{
				"teamIds": {"119"}, "startSeason": {"1960"},
				"endSeason": {"2024"}, "fields": {"teams,id,name"},
			},
			wantLen: 1,
		},
		{
			name: "missing TeamIDs", query: TeamsHistoryQuery{},
			wantErr: "TeamIDs is required", wantIs: ErrInvalidQuery,
		},
		{
			name: "empty", query: TeamsHistoryQuery{TeamIDs: "119"},
			respStatus: 200, respBody: `{}`, wantLen: 0,
		},
		{
			name: "404", query: TeamsHistoryQuery{TeamIDs: "119"},
			respStatus: 404, respBody: `{}`, wantIs: ErrNotFound, wantErr: "not found",
		},
		{
			name: "5xx", query: TeamsHistoryQuery{TeamIDs: "119"},
			respStatus: 500, respBody: `oops`, wantErr: "unexpected status 500",
		},
		{
			name: "bad json", query: TeamsHistoryQuery{TeamIDs: "119"},
			respStatus: 200, respBody: `x`, wantErr: "teamsHistory",
		},
		{
			name: "network", query: TeamsHistoryQuery{TeamIDs: "119"},
			respStatus: 0, wantErr: "teamsHistory",
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
			ts, err := client.TeamsHistory(context.Background(), c.query)
			if c.wantErr != "" {
				if err == nil {
					t.Fatalf("expected error containing %q", c.wantErr)
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
			if len(ts.Teams) != c.wantLen {
				t.Errorf("len = %d, want %d", len(ts.Teams), c.wantLen)
			}
		})
	}
}
//
func TestClient_ScheduleTied(t *testing.T) {
	cases := []struct {
		name       string
		query      ScheduleTiedQuery
		respStatus int
		respBody   string
		wantPath   string
		wantQuery  url.Values
		wantErr    string
		wantIs     error
		wantLen    int
	}{
		{
			name: "happy path",
			query: ScheduleTiedQuery{
				Season: 2024, GameTypes: "R", Hydrate: "team",
				Fields: "dates,games,gamePk",
			},
			respStatus: 200,
			respBody:   `{"dates":[{"date":"2024-05-01","games":[{"gamePk":1}]}]}`,
			wantPath:   "/api/v1/schedule/games/tied",
			wantQuery: url.Values{
				"season": {"2024"}, "gameTypes": {"R"},
				"hydrate": {"team"}, "fields": {"dates,games,gamePk"},
			},
			wantLen: 1,
		},
		{
			name: "missing Season", query: ScheduleTiedQuery{},
			wantErr: "Season is required", wantIs: ErrInvalidQuery,
		},
		{
			name: "empty", query: ScheduleTiedQuery{Season: 2024},
			respStatus: 200, respBody: `{}`, wantLen: 0,
		},
		{
			name: "404", query: ScheduleTiedQuery{Season: 2024},
			respStatus: 404, respBody: `{}`, wantIs: ErrNotFound, wantErr: "not found",
		},
		{
			name: "5xx", query: ScheduleTiedQuery{Season: 2024},
			respStatus: 500, respBody: `oops`, wantErr: "unexpected status 500",
		},
		{
			name: "bad json", query: ScheduleTiedQuery{Season: 2024},
			respStatus: 200, respBody: `x`, wantErr: "scheduleTied",
		},
		{
			name: "network", query: ScheduleTiedQuery{Season: 2024},
			respStatus: 0, wantErr: "scheduleTied",
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
			games, err := client.ScheduleTied(context.Background(), c.query)
			if c.wantErr != "" {
				if err == nil {
					t.Fatalf("expected error containing %q", c.wantErr)
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
			if len(games) != c.wantLen {
				t.Errorf("len = %d, want %d", len(games), c.wantLen)
			}
		})
	}
}
