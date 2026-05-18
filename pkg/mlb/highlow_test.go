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
const highLowHappyBody = `{
  "highLowResults": [{
    "group": {"displayName": "hitting"},
    "totalSplits": 303,
    "splits": [{
      "season": "2024",
      "stat": {"homeRuns": 3},
      "team": {"id": 117, "name": "Houston Astros"},
      "player": {"id": 670541, "fullName": "Yordan Alvarez", "link": "/api/v1/people/670541"}
    }]
  }]
}`
//
func TestClient_HighLow(t *testing.T) {
	cases := []struct {
		name         string
		orgType      string
		query        HighLowQuery
		respStatus   int
		respBody     string
		wantPath     string
		wantQuery    url.Values
		wantErr      string
		wantIs       error
		wantLen      int
		wantHydrated bool
	}{
		{
			name:    "happy path: every filter set",
			orgType: "player",
			query: HighLowQuery{
				SortStat: "homeRuns", Season: 2024,
				GameType: "R", TeamID: 117, LeagueID: 103,
				SportIDs: "1", StatGroup: "hitting", Limit: 3,
				Fields: "highLowResults,splits",
			},
			respStatus: 200,
			respBody:   highLowHappyBody,
			wantPath:   "/api/v1/highLow/player",
			wantQuery: url.Values{
				"sortStat": {"homeRuns"}, "season": {"2024"},
				"gameType": {"R"}, "teamId": {"117"}, "leagueId": {"103"},
				"sportIds": {"1"}, "statGroup": {"hitting"}, "limit": {"3"},
				"fields": {"highLowResults,splits"},
			},
			wantLen:      1,
			wantHydrated: true,
		},
		{
			name:    "missing SortStat rejected",
			orgType: "player",
			query:   HighLowQuery{Season: 2024},
			wantErr: "SortStat and Season are both required",
			wantIs:  ErrInvalidQuery,
		},
		{
			name:    "missing Season rejected",
			orgType: "player",
			query:   HighLowQuery{SortStat: "homeRuns"},
			wantErr: "SortStat and Season are both required",
			wantIs:  ErrInvalidQuery,
		},
		{
			name:       "200 with no results yields empty slice",
			orgType:    "player",
			query:      HighLowQuery{SortStat: "homeRuns", Season: 2024},
			respStatus: 200,
			respBody:   `{}`,
			wantLen:    0,
		},
		{
			name:       "404 returns ErrNotFound",
			orgType:    "player",
			query:      HighLowQuery{SortStat: "homeRuns", Season: 2024},
			respStatus: 404,
			respBody:   `{}`,
			wantIs:     ErrNotFound,
			wantErr:    "not found",
		},
		{
			name:       "5xx is wrapped",
			orgType:    "player",
			query:      HighLowQuery{SortStat: "homeRuns", Season: 2024},
			respStatus: 500,
			respBody:   `oops`,
			wantErr:    "unexpected status 500",
		},
		{
			name:       "malformed JSON is wrapped",
			orgType:    "player",
			query:      HighLowQuery{SortStat: "homeRuns", Season: 2024},
			respStatus: 200,
			respBody:   `not json`,
			wantErr:    "highLow",
		},
		{
			name:       "network failure is wrapped",
			orgType:    "player",
			query:      HighLowQuery{SortStat: "homeRuns", Season: 2024},
			respStatus: 0,
			wantErr:    "highLow",
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
			hl, err := client.HighLow(context.Background(), c.orgType, c.query)
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
			if hl == nil {
				t.Fatal("expected non-nil HighLow")
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
			if got := len(hl.Results); got != c.wantLen {
				t.Errorf("len(Results) = %d, want %d", got, c.wantLen)
			}
			if !c.wantHydrated {
				return
			}
			g := hl.Results[0]
			if g.Group != "hitting" || g.TotalSplits != 303 || len(g.Splits) != 1 {
				t.Errorf("Results[0] = group=%q total=%d splits=%d",
					g.Group, g.TotalSplits, len(g.Splits))
			}
			s := g.Splits[0]
			if s.Season != "2024" || s.Player.ID != 670541 ||
				s.Player.FullName != "Yordan Alvarez" ||
				s.Team.ID != TeamID(117) || s.Team.Name != "Houston Astros" {
				t.Errorf("Splits[0] = %+v", s)
			}
			if hr, ok := s.Stat["homeRuns"]; !ok || hr != float64(3) {
				t.Errorf("Stat = %+v", s.Stat)
			}
		})
	}
}
