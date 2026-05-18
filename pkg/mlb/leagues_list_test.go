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
const leaguesHappyBody = `{
  "copyright": "Copyright 2026 MLB Advanced Media, L.P.",
  "leagues": [
    {
      "id": 103, "name": "American League", "link": "/api/v1/league/103",
      "abbreviation": "AL", "nameShort": "American", "seasonState": "inseason",
      "hasWildCard": true, "hasSplitSeason": false, "numGames": 162,
      "hasPlayoffPoints": false, "numTeams": 15, "numWildcardTeams": 3,
      "season": "2026", "orgCode": "AL",
      "conferencesInUse": false, "divisionsInUse": true,
      "sport": {"id": 1, "link": "/api/v1/sports/1"},
      "sortOrder": 21, "active": true,
      "seasonDateInfo": {
        "seasonId": "2026",
        "regularSeasonStartDate": "2026-03-25",
        "regularSeasonEndDate":   "2026-09-27"
      }
    },
    {
      "id": 104, "name": "National League", "link": "/api/v1/league/104"
    }
  ]
}`
//
func TestLeagues_League(t *testing.T) {
	empty := leaguesFromGen(nil)
	full := &Leagues{Leagues: []LeagueInfo{
		{ID: 103, Name: "AL"},
		{ID: 104, Name: "NL"},
	}}
	cases := []struct {
		name     string
		l        *Leagues
		lookup   int
		wantOk   bool
		wantName string
	}{
		{"nil receiver", nil, 103, false, ""},
		{"empty (no leagues)", empty, 103, false, ""},
		{"miss", full, 999, false, ""},
		{"hit AL", full, 103, true, "AL"},
		{"hit NL", full, 104, true, "NL"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := c.l.League(c.lookup)
			if (got != nil) != c.wantOk {
				t.Fatalf("League(%d) ok=%v, want %v", c.lookup, got != nil, c.wantOk)
			}
			if got != nil && got.Name != c.wantName {
				t.Errorf("League(%d).Name = %q, want %q", c.lookup, got.Name, c.wantName)
			}
		})
	}
}
//
func TestClient_Leagues(t *testing.T) {
	cases := []struct {
		name       string
		query      LeaguesQuery
		respStatus int
		respBody   string
		wantPath   string
		wantQuery  url.Values
		wantErr    string
		wantIs     error
		wantLen    int
		// Set when we expect the fully-hydrated first row to be audited.
		wantHydrated bool
	}{
		{
			name: "happy path: SportID + Seasons + Fields",
			query: LeaguesQuery{
				SportID: 1,
				Seasons: "2026",
				Fields:  "leagues,id,name",
			},
			respStatus: 200,
			respBody:   leaguesHappyBody,
			wantPath:   "/api/v1/league",
			wantQuery: url.Values{
				"sportId": {"1"},
				"seasons": {"2026"},
				"fields":  {"leagues,id,name"},
			},
			wantLen:      2,
			wantHydrated: true,
		},
		{
			name: "filter with LeagueIDs instead of SportID",
			query: LeaguesQuery{
				LeagueIDs: "103,104",
			},
			respStatus: 200,
			respBody:   `{"leagues":[{"id":103,"name":"AL"}]}`,
			wantPath:   "/api/v1/league",
			wantQuery: url.Values{
				"leagueIds": {"103,104"},
			},
			wantLen: 1,
		},
		{
			name:    "missing required combo rejected before HTTP",
			query:   LeaguesQuery{Seasons: "2026"},
			wantErr: "one of SportID or LeagueIDs",
			wantIs:  ErrInvalidQuery,
		},
		{
			name:       "200 with no leagues yields empty slice",
			query:      LeaguesQuery{SportID: 1},
			respStatus: 200,
			respBody:   `{}`,
			wantLen:    0,
		},
		{
			name:       "200 with explicit empty leagues array",
			query:      LeaguesQuery{SportID: 1},
			respStatus: 200,
			respBody:   `{"leagues": []}`,
			wantLen:    0,
		},
		{
			name:       "404 returns ErrNotFound",
			query:      LeaguesQuery{SportID: 1},
			respStatus: 404,
			respBody:   `{}`,
			wantIs:     ErrNotFound,
			wantErr:    "not found",
		},
		{
			name:       "5xx is wrapped",
			query:      LeaguesQuery{SportID: 1},
			respStatus: 500,
			respBody:   `oops`,
			wantErr:    "unexpected status 500",
		},
		{
			name:       "malformed JSON is wrapped",
			query:      LeaguesQuery{SportID: 1},
			respStatus: 200,
			respBody:   `not json`,
			wantErr:    "leagues",
		},
		{
			name:       "network failure is wrapped",
			query:      LeaguesQuery{SportID: 1},
			respStatus: 0,
			wantErr:    "leagues",
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
			ls, err := client.Leagues(context.Background(), c.query)
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
			if ls == nil {
				t.Fatal("expected non-nil Leagues")
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
			if got := len(ls.Leagues); got != c.wantLen {
				t.Errorf("len(Leagues) = %d, want %d", got, c.wantLen)
			}
			if !c.wantHydrated {
				return
			}
//
			first := ls.Leagues[0]
			if first.ID != 103 || first.Name != "American League" ||
				first.Abbreviation != "AL" || first.NameShort != "American" ||
				first.SeasonState != "inseason" || !first.HasWildCard ||
				first.HasSplitSeason || first.NumGames != 162 ||
				first.HasPlayoffPoints || first.NumTeams != 15 ||
				first.NumWildcardTeams != 3 || first.Season != "2026" ||
				first.OrgCode != "AL" || first.ConferencesInUse ||
				!first.DivisionsInUse || first.SortOrder != 21 ||
				!first.Active || first.Link != "/api/v1/league/103" {
				t.Errorf("Leagues[0] = %+v", first)
			}
			if first.Sport.ID != 1 || first.Sport.Link != "/api/v1/sports/1" {
				t.Errorf("Leagues[0].Sport = %+v", first.Sport)
			}
			if first.SeasonDateInfo.SeasonID != "2026" ||
				first.SeasonDateInfo.RegularSeasonStartDate.Format(seasonDateFmt) !=
					"2026-03-25" {
				t.Errorf("SeasonDateInfo = %+v", first.SeasonDateInfo)
			}
		})
	}
}
