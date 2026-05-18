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
const statsLeadersHappyBody = `{
  "leagueLeaders": [{
    "leaderCategory": "homeRuns", "season": "2024",
    "gameType": {"displayName": "Regular Season"},
    "statGroup": "hitting", "totalSplits": 300,
    "leaders": [{
      "rank": 1, "value": "58",
      "team":   {"id": 147, "name": "New York Yankees"},
      "league": {"id": 103, "link": "/api/v1/league/103"},
      "person": {"id": 592450, "fullName": "Aaron Judge", "link": "/api/v1/people/592450"},
      "sport":  {"id": 1, "link": "/api/v1/sports/1"}
    }]
  }]
}`
//
func TestClient_StatsLeaders(t *testing.T) {
	cases := []struct {
		name         string
		query        StatsLeadersQuery
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
			name: "happy path: every filter set",
			query: StatsLeadersQuery{
				LeaderCategories: "homeRuns", Season: 2024, SportID: 1,
				LeagueID: 103, StatGroup: "hitting", PlayerPool: "all",
				LeaderGameTypes: "R", StatType: "statsSingleSeason",
				Hydrate: "person", Limit: 5, Fields: "leagueLeaders,leaders",
			},
			respStatus: 200,
			respBody:   statsLeadersHappyBody,
			wantPath:   "/api/v1/stats/leaders",
			wantQuery: url.Values{
				"leaderCategories": {"homeRuns"}, "season": {"2024"},
				"sportId": {"1"}, "leagueId": {"103"},
				"statGroup": {"hitting"}, "playerPool": {"all"},
				"leaderGameTypes": {"R"}, "statType": {"statsSingleSeason"},
				"hydrate": {"person"}, "limit": {"5"},
				"fields": {"leagueLeaders,leaders"},
			},
			wantLen:      1,
			wantHydrated: true,
		},
		{
			name:    "missing LeaderCategories rejected",
			query:   StatsLeadersQuery{Season: 2024},
			wantErr: "LeaderCategories is required",
			wantIs:  ErrInvalidQuery,
		},
		{
			name:       "200 with no leagueLeaders yields empty slice",
			query:      StatsLeadersQuery{LeaderCategories: "homeRuns"},
			respStatus: 200,
			respBody:   `{}`,
			wantLen:    0,
		},
		{
			name:       "404 returns ErrNotFound",
			query:      StatsLeadersQuery{LeaderCategories: "homeRuns"},
			respStatus: 404,
			respBody:   `{}`,
			wantIs:     ErrNotFound,
			wantErr:    "not found",
		},
		{
			name:       "5xx is wrapped",
			query:      StatsLeadersQuery{LeaderCategories: "homeRuns"},
			respStatus: 500,
			respBody:   `oops`,
			wantErr:    "unexpected status 500",
		},
		{
			name:       "malformed JSON is wrapped",
			query:      StatsLeadersQuery{LeaderCategories: "homeRuns"},
			respStatus: 200,
			respBody:   `not json`,
			wantErr:    "statsLeaders",
		},
		{
			name:       "network failure is wrapped",
			query:      StatsLeadersQuery{LeaderCategories: "homeRuns"},
			respStatus: 0,
			wantErr:    "statsLeaders",
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
			sl, err := client.StatsLeaders(context.Background(), c.query)
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
			if sl == nil {
				t.Fatal("expected non-nil StatsLeaders")
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
			if got := len(sl.LeagueLeaders); got != c.wantLen {
				t.Errorf("len(LeagueLeaders) = %d, want %d", got, c.wantLen)
			}
			if !c.wantHydrated {
				return
			}
			cat := sl.LeagueLeaders[0]
			if cat.LeaderCategory != "homeRuns" || cat.Season != "2024" ||
				cat.GameType != "Regular Season" || cat.StatGroup != "hitting" ||
				cat.TotalSplits != 300 || len(cat.Leaders) != 1 {
				t.Errorf("LeagueLeaders[0] = %+v", cat)
			}
			l := cat.Leaders[0]
			if l.Rank != 1 || l.Value != "58" ||
				l.Team.ID != NYY || l.Team.Name != "New York Yankees" ||
				l.League.ID != 103 || l.Player.ID != 592450 ||
				l.Player.FullName != "Aaron Judge" || l.Sport.ID != 1 {
				t.Errorf("Leaders[0] = %+v", l)
			}
		})
	}
}
