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

const standingsHappyBody = `{
  "records": [{
    "standingsType": "regularSeason",
    "league":   {"id": 104, "link": "/api/v1/league/104"},
    "division": {"id": 204, "link": "/api/v1/divisions/204"},
    "sport":    {"id": 1,   "link": "/api/v1/sports/1"},
    "lastUpdated": "2026-02-05T00:42:15.121Z",
    "teamRecords": [
      {
        "team": {"id": 119, "name": "Los Angeles Dodgers"},
        "streak": {"streakCode": "W3", "streakType": "wins", "streakNumber": 3},
        "wins": 96, "losses": 66, "gamesPlayed": 162,
        "runsScored": 800, "runsAllowed": 670, "runDifferential": 130,
        "winningPercentage": ".593",
        "gamesBack": "-", "wildCardGamesBack": "-",
        "divisionRank": "1", "leagueRank": "1", "sportRank": "1",
        "eliminationNumber": "-", "magicNumber": "-",
        "clinched": true, "divisionLeader": true, "divisionChamp": true, "hasWildcard": false,
        "season": "2026", "lastUpdated": "2026-02-05T00:42:15.121Z"
      },
      {
        "team": {"id": 137, "name": "San Francisco Giants"},
        "streak": {"streakCode": "L1", "streakType": "losses", "streakNumber": 1},
        "wins": 80, "losses": 82, "gamesPlayed": 162,
        "runsScored": 700, "runsAllowed": 720, "runDifferential": -20,
        "winningPercentage": ".494",
        "gamesBack": "16.0", "wildCardGamesBack": "-",
        "divisionRank": "2", "leagueRank": "8", "sportRank": "16",
        "eliminationNumber": "E", "magicNumber": "-",
        "clinched": false, "divisionLeader": false, "divisionChamp": false, "hasWildcard": false,
        "season": "2026", "lastUpdated": "2026-02-05T00:42:15.121Z"
      }
    ]
  }]
}`

func TestStandings_Division(t *testing.T) {
	st := standingsFromGen(nil) // nil → empty Standings
	full := &Standings{Records: []DivisionStandings{
		{Division: Ref{ID: 200}},
		{Division: Ref{ID: 204}},
	}}

	cases := []struct {
		name       string
		standings  *Standings
		divisionID int
		wantOk     bool
		wantID     int
	}{
		{"nil receiver", nil, 204, false, 0},
		{"empty Standings (no records)", st, 204, false, 0},
		{"miss", full, 999, false, 0},
		{"hit", full, 204, true, 204},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := c.standings.Division(c.divisionID)
			if (got != nil) != c.wantOk {
				t.Fatalf("Division(%d) ok=%v, want %v", c.divisionID, got != nil, c.wantOk)
			}
			if got != nil && got.Division.ID != c.wantID {
				t.Errorf(
					"Division(%d).Division.ID = %d, want %d",
					c.divisionID,
					got.Division.ID,
					c.wantID,
				)
			}
		})
	}
}

func TestDivisionStandings_Team(t *testing.T) {
	d := &DivisionStandings{TeamRecords: []TeamRecord{
		{Team: TeamRef{ID: LAD, Name: "Los Angeles Dodgers"}},
		{Team: TeamRef{ID: SF, Name: "San Francisco Giants"}},
	}}
	cases := []struct {
		name     string
		division *DivisionStandings
		lookup   TeamID
		wantOk   bool
		wantName string
	}{
		{"nil receiver", nil, LAD, false, ""},
		{"miss", d, NYY, false, ""},
		{"hit Dodgers", d, LAD, true, "Los Angeles Dodgers"},
		{"hit Giants", d, SF, true, "San Francisco Giants"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := c.division.Team(c.lookup)
			if (got != nil) != c.wantOk {
				t.Fatalf("Team(%v) ok=%v, want %v", c.lookup, got != nil, c.wantOk)
			}
			if got != nil && got.Team.Name != c.wantName {
				t.Errorf("Team(%v).Team.Name = %q, want %q", c.lookup, got.Team.Name, c.wantName)
			}
		})
	}
}

func TestClient_Standings(t *testing.T) {
	cases := []struct {
		name        string
		query       StandingsQuery
		respStatus  int
		respBody    string
		wantPath    string
		wantQuery   url.Values
		wantErr     string
		wantIs      error
		wantNumDivs int
		wantWins    int // for happy rows; checks first division's first team Wins
		wantStreak  string
	}{
		{
			name: "happy path: NL East with 2 teams",
			query: StandingsQuery{
				League:         NL,
				Season:         2026,
				StandingsTypes: "regularSeason",
				On:             time.Date(2026, 9, 28, 0, 0, 0, 0, time.UTC),
				Hydrate:        "team",
			},
			respStatus: 200,
			respBody:   standingsHappyBody,
			wantPath:   "/api/v1/standings",
			wantQuery: url.Values{
				"leagueId":       {"104"},
				"season":         {"2026"},
				"standingsTypes": {"regularSeason"},
				"date":           {"2026-09-28"},
				"hydrate":        {"team"},
			},
			wantNumDivs: 1,
			wantWins:    96,
			wantStreak:  "W3",
		},
		{
			name:        "200 with no records yields empty slice",
			query:       StandingsQuery{League: NL},
			respStatus:  200,
			respBody:    `{}`,
			wantNumDivs: 0,
		},
		{
			name:    "missing required League is rejected before HTTP call",
			query:   StandingsQuery{Season: 2026},
			wantErr: "League is required",
			wantIs:  ErrInvalidQuery,
		},
		{
			name:       "404 returns ErrNotFound",
			query:      StandingsQuery{League: NL},
			respStatus: 404,
			respBody:   `{}`,
			wantIs:     ErrNotFound,
			wantErr:    "not found",
		},
		{
			name:       "5xx is wrapped",
			query:      StandingsQuery{League: NL},
			respStatus: 500,
			respBody:   `oops`,
			wantErr:    "unexpected status 500",
		},
		{
			name:       "malformed JSON is wrapped",
			query:      StandingsQuery{League: NL},
			respStatus: 200,
			respBody:   `not json`,
			wantErr:    "standings",
		},
		{
			name:       "network failure is wrapped",
			query:      StandingsQuery{League: NL},
			respStatus: 0,
			wantErr:    "standings",
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
			st, err := client.Standings(context.Background(), c.query)

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
			if st == nil {
				t.Fatal("expected non-nil Standings")
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
			if got := len(st.Records); got != c.wantNumDivs {
				t.Errorf("len(Records) = %d, want %d", got, c.wantNumDivs)
			}
			if c.wantWins > 0 {
				rec := st.Records[0]
				if rec.League.ID != 104 || rec.League.Link != "/api/v1/league/104" {
					t.Errorf("League = %+v, want {104 /api/v1/league/104}", rec.League)
				}
				if rec.Division.ID != 204 || rec.Division.Link != "/api/v1/divisions/204" {
					t.Errorf("Division = %+v, want {204 /api/v1/divisions/204}", rec.Division)
				}
				if rec.Sport.ID != 1 || rec.Sport.Link != "/api/v1/sports/1" {
					t.Errorf("Sport = %+v, want {1 /api/v1/sports/1}", rec.Sport)
				}
				tr := rec.Team(LAD)
				if tr == nil {
					t.Fatal("expected Dodgers in first division")
				}
				if tr.Wins != c.wantWins {
					t.Errorf("Dodgers.Wins = %d, want %d", tr.Wins, c.wantWins)
				}
				if tr.Streak.Code != c.wantStreak {
					t.Errorf("Dodgers.Streak.Code = %q, want %q", tr.Streak.Code, c.wantStreak)
				}
			}
		})
	}
}
