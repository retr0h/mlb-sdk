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
const teamsListHappyBody = `{
  "copyright": "Copyright 2026 MLB Advanced Media, L.P.",
  "teams": [
    {
      "id": 119, "name": "Los Angeles Dodgers",
      "link": "/api/v1/teams/119", "season": 2024,
      "abbreviation": "LAD", "teamName": "Dodgers",
      "locationName": "Los Angeles", "active": true
    },
    {
      "id": 137, "name": "San Francisco Giants",
      "link": "/api/v1/teams/137", "season": 2024,
      "abbreviation": "SF", "teamName": "Giants",
      "locationName": "San Francisco", "active": true
    }
  ]
}`
//
func TestTeams_Team(t *testing.T) {
	empty := teamsFromGen(nil)
	full := &Teams{Teams: []TeamInfo{
		{ID: 119, Name: "Dodgers"},
		{ID: 137, Name: "Giants"},
	}}
	cases := []struct {
		name     string
		t        *Teams
		lookup   int
		wantOk   bool
		wantName string
	}{
		{"nil receiver", nil, 119, false, ""},
		{"empty (no teams)", empty, 119, false, ""},
		{"miss", full, 999, false, ""},
		{"hit Dodgers", full, 119, true, "Dodgers"},
		{"hit Giants", full, 137, true, "Giants"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := c.t.Team(c.lookup)
			if (got != nil) != c.wantOk {
				t.Fatalf("Team(%d) ok=%v, want %v", c.lookup, got != nil, c.wantOk)
			}
			if got != nil && got.Name != c.wantName {
				t.Errorf("Team(%d).Name = %q, want %q", c.lookup, got.Name, c.wantName)
			}
		})
	}
}
//
func TestClient_Teams(t *testing.T) {
	cases := []struct {
		name       string
		query      TeamsQuery
		respStatus int
		respBody   string
		wantPath   string
		wantQuery  url.Values
		wantErr    string
		wantIs     error
		wantLen    int
		wantFirst  string
	}{
		{
			name: "happy path: every filter set",
			query: TeamsQuery{
				Season:       2024,
				ActiveStatus: "ACTIVE",
				LeagueIDs:    "103,104",
				SportID:      1,
				SportIDs:     "1,11",
				GameType:     "R",
				Hydrate:      "league",
				Fields:       "teams,id,name",
			},
			respStatus: 200,
			respBody:   teamsListHappyBody,
			wantPath:   "/api/v1/teams",
			wantQuery: url.Values{
				"season":       {"2024"},
				"activeStatus": {"ACTIVE"},
				"leagueIds":    {"103,104"},
				"sportId":      {"1"},
				"sportIds":     {"1,11"},
				"gameType":     {"R"},
				"hydrate":      {"league"},
				"fields":       {"teams,id,name"},
			},
			wantLen:   2,
			wantFirst: "Los Angeles Dodgers",
		},
		{
			name:       "empty query yields no query params",
			query:      TeamsQuery{},
			respStatus: 200,
			respBody:   teamsListHappyBody,
			wantPath:   "/api/v1/teams",
			wantQuery:  url.Values{},
			wantLen:    2,
			wantFirst:  "Los Angeles Dodgers",
		},
		{
			name:       "200 with no teams yields empty slice",
			respStatus: 200,
			respBody:   `{}`,
			wantLen:    0,
		},
		{
			name:       "200 with explicit empty teams array",
			respStatus: 200,
			respBody:   `{"teams": []}`,
			wantLen:    0,
		},
		{
			name:       "404 returns ErrNotFound",
			respStatus: 404,
			respBody:   `{}`,
			wantIs:     ErrNotFound,
			wantErr:    "not found",
		},
		{
			name:       "5xx is wrapped",
			respStatus: 500,
			respBody:   `oops`,
			wantErr:    "unexpected status 500",
		},
		{
			name:       "malformed JSON is wrapped",
			respStatus: 200,
			respBody:   `not json`,
			wantErr:    "teams",
		},
		{
			name:       "network failure is wrapped",
			respStatus: 0,
			wantErr:    "teams",
		},
	}
//
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
//
			client := New(WithBaseURL(urlStr))
			ts, err := client.Teams(context.Background(), c.query)
//
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
			if ts == nil {
				t.Fatal("expected non-nil Teams")
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
			if got := len(ts.Teams); got != c.wantLen {
				t.Errorf("len(Teams) = %d, want %d", got, c.wantLen)
			}
			if c.wantFirst != "" {
				if ts.Teams[0].Name != c.wantFirst {
					t.Errorf("Teams[0].Name = %q, want %q", ts.Teams[0].Name, c.wantFirst)
				}
			}
		})
	}
}
