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
)

const gamePaceHappyBody = `{
  "sports": [{
    "season": "2024",
    "sport": {"id": 1, "link": "/api/v1/sports/1"},
    "hitsPer9Inn": 16.63, "runsPer9Inn": 8.91,
    "pitchesPer9Inn": 296.2, "plateAppearancesPer9Inn": 76.17,
    "hitsPerGame": 16.39, "runsPerGame": 8.79,
    "inningsPlayedPerGame": 8.9, "pitchesPerGame": 292.1,
    "pitchersPerGame": 8.52, "plateAppearancesPerGame": 75.11,
    "hitsPerRun": 1.87, "pitchesPerPitcher": 34.27,
    "totalGameTime": "6426:37:00", "totalInningsPlayed": 21626.5,
    "totalHits": 39823, "totalRuns": 21343,
    "totalPlateAppearances": 182449, "totalPitchers": 20695,
    "totalPitches": 709511, "totalGames": 2429,
    "total7InnGames": 4, "total9InnGames": 2209,
    "total9InnGamesCompletedEarly": 1, "total9InnGamesScheduled": 2430,
    "total9InnGamesWithoutExtraInn": 2209,
    "totalExtraInnGames": 216, "totalExtraInnTime": "123:45:00",
    "timePerGame": "02:38:44", "timePerPitch": "00:00:32",
    "timePerHit": "00:09:40", "timePerRun": "00:18:04",
    "timePerPlateAppearance": "00:02:06",
    "timePer9Inn": "02:40:58",
    "timePer77PlateAppearances": "02:42:44",
    "timePer7InnGameWithoutExtraInn": "02:05:30"
  }]
}`

func TestClient_GamePace(t *testing.T) {
	cases := []struct {
		name         string
		query        GamePaceQuery
		respStatus   int
		respBody     string
		wantPath     string
		wantQuery    url.Values
		wantErr      string
		wantIs       error
		wantHydrated bool
	}{
		{
			name: "happy path: every filter set",
			query: GamePaceQuery{
				Season: 2024, SportID: 1, TeamIDs: "119",
				LeagueIDs: "103", LeagueListID: "mlb",
				GameType: "R", StartDate: "2024-04-01", EndDate: "2024-09-30",
				VenueIDs: "22", OrgType: "T", IncludeChildren: true,
				Fields: "sports,season",
			},
			respStatus: 200,
			respBody:   gamePaceHappyBody,
			wantPath:   "/api/v1/gamePace",
			wantQuery: url.Values{
				"season": {"2024"}, "sportId": {"1"}, "teamIds": {"119"},
				"leagueIds": {"103"}, "leagueListId": {"mlb"},
				"gameType": {"R"}, "startDate": {"2024-04-01"},
				"endDate": {"2024-09-30"}, "venueIds": {"22"},
				"orgType": {"T"}, "includeChildren": {"true"},
				"fields": {"sports,season"},
			},
			wantHydrated: true,
		},
		{
			name:    "missing Season rejected",
			query:   GamePaceQuery{SportID: 1},
			wantErr: "Season is required",
			wantIs:  ErrInvalidQuery,
		},
		{
			name:       "200 with empty sports maps to ErrNotFound",
			query:      GamePaceQuery{Season: 2024},
			respStatus: 200,
			respBody:   `{"sports":[]}`,
			wantIs:     ErrNotFound,
			wantErr:    "not found",
		},
		{
			name:       "200 with missing sports key maps to ErrNotFound",
			query:      GamePaceQuery{Season: 2024},
			respStatus: 200,
			respBody:   `{}`,
			wantIs:     ErrNotFound,
			wantErr:    "not found",
		},
		{
			name:       "404 returns ErrNotFound",
			query:      GamePaceQuery{Season: 2024},
			respStatus: 404,
			respBody:   `{}`,
			wantIs:     ErrNotFound,
			wantErr:    "not found",
		},
		{
			name:       "5xx is wrapped",
			query:      GamePaceQuery{Season: 2024},
			respStatus: 500,
			respBody:   `oops`,
			wantErr:    "unexpected status 500",
		},
		{
			name:       "malformed JSON is wrapped",
			query:      GamePaceQuery{Season: 2024},
			respStatus: 200,
			respBody:   `not json`,
			wantErr:    "gamePace",
		},
		{
			name:       "network failure is wrapped",
			query:      GamePaceQuery{Season: 2024},
			respStatus: 0,
			wantErr:    "gamePace",
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
			gp, err := client.GamePace(context.Background(), c.query)

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
			if gp == nil {
				t.Fatal("expected non-nil GamePace")
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
			if !c.wantHydrated {
				return
			}
			if gp.Season != "2024" || gp.Sport.ID != 1 {
				t.Errorf("Season=%q Sport=%+v", gp.Season, gp.Sport)
			}
			if gp.TotalGames != 2429 || gp.TotalHits != 39823 ||
				gp.TotalRuns != 21343 || gp.TotalPitches != 709511 {
				t.Errorf("totals = games=%d hits=%d runs=%d pitches=%d",
					gp.TotalGames, gp.TotalHits, gp.TotalRuns, gp.TotalPitches)
			}
			if gp.TimePerGame != "02:38:44" || gp.TimePer9Inn != "02:40:58" {
				t.Errorf("time = perGame=%q per9Inn=%q", gp.TimePerGame, gp.TimePer9Inn)
			}
			if gp.HitsPer9Inn != 16.63 || gp.RunsPerGame != 8.79 {
				t.Errorf("rates = hitsPer9Inn=%v runsPerGame=%v",
					gp.HitsPer9Inn, gp.RunsPerGame)
			}
			if gp.Total7InnGames != 4 || gp.Total9InnGames != 2209 ||
				gp.TotalExtraInnGames != 216 || gp.Total9InnGamesCompletedEarly != 1 ||
				gp.Total9InnGamesScheduled != 2430 ||
				gp.Total9InnGamesWithoutExtraInn != 2209 {
				t.Errorf("game breakdowns wrong")
			}
			if gp.TotalExtraInnTime != "123:45:00" ||
				gp.TimePer77PlateAppearances != "02:42:44" ||
				gp.TimePer7InnGameWithoutExtraInn != "02:05:30" {
				t.Errorf("extended time fields wrong")
			}
			if gp.HitsPerRun != 1.87 || gp.PitchesPerPitcher != 34.27 {
				t.Errorf("derived rates = hitsPerRun=%v pitchesPerPitcher=%v",
					gp.HitsPerRun, gp.PitchesPerPitcher)
			}
		})
	}
}
