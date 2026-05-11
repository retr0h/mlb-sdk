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

const teamHappyBody = `{
  "copyright": "Copyright 2026 MLB Advanced Media, L.P.",
  "teams": [{
    "id": 119, "name": "Los Angeles Dodgers",
    "link": "/api/v1/teams/119", "season": 2026,
    "allStarStatus": "N", "active": true,
    "teamCode": "lan", "fileCode": "la", "abbreviation": "LAD",
    "teamName": "Dodgers", "locationName": "Los Angeles",
    "firstYearOfPlay": "1884",
    "shortName": "LA Dodgers", "franchiseName": "Los Angeles", "clubName": "Dodgers",
    "venue": {
      "id": 22, "name": "Dodger Stadium",
      "link": "/api/v1/venues/22", "active": true, "season": "2026"
    },
    "springLeague": {
      "id": 114, "name": "Cactus League",
      "link": "/api/v1/league/114", "abbreviation": "CL"
    },
    "springVenue": {"id": 3809, "link": "/api/v1/venues/3809"},
    "league": {
      "id": 104, "name": "National League", "link": "/api/v1/league/104",
      "abbreviation": "NL", "nameShort": "National", "seasonState": "inseason",
      "hasWildCard": true, "hasSplitSeason": false, "numGames": 162,
      "hasPlayoffPoints": false, "numTeams": 15, "numWildcardTeams": 3,
      "season": "2026", "orgCode": "NL",
      "conferencesInUse": false, "divisionsInUse": true,
      "sport": {"id": 1, "link": "/api/v1/sports/1"},
      "sortOrder": 31, "active": true,
      "seasonDateInfo": {
        "seasonId": "2026",
        "regularSeasonStartDate": "2026-03-25",
        "regularSeasonEndDate":   "2026-09-27"
      }
    },
    "division": {
      "id": 203, "name": "National League West", "season": "2026",
      "nameShort": "NL West", "link": "/api/v1/divisions/203",
      "abbreviation": "NLW",
      "league": {"id": 104, "link": "/api/v1/league/104"},
      "sport":  {"id": 1,   "link": "/api/v1/sports/1"},
      "hasWildcard": false, "sortOrder": 34, "numPlayoffTeams": 1, "active": true
    },
    "sport": {
      "id": 1, "code": "mlb", "link": "/api/v1/sports/1",
      "name": "Major League Baseball", "abbreviation": "MLB",
      "sortOrder": 11, "activeStatus": true
    }
  }]
}`

func TestClient_Team(t *testing.T) {
	cases := []struct {
		name       string
		teamID     int
		query      TeamQuery
		respStatus int
		respBody   string
		wantPath   string
		wantQuery  url.Values
		wantErr    string
		wantIs     error
		// Set when we expect the fully-hydrated happy body.
		wantHydrated bool
		wantName     string
	}{
		{
			name:   "happy path: hydrated team",
			teamID: 119,
			query: TeamQuery{
				Season:  2026,
				SportID: 1,
				Hydrate: "league,division,sport,springLeague,venue",
				Fields:  "teams,id,name",
			},
			respStatus: 200,
			respBody:   teamHappyBody,
			wantPath:   "/api/v1/teams/119",
			wantQuery: url.Values{
				"season":  {"2026"},
				"sportId": {"1"},
				"hydrate": {"league,division,sport,springLeague,venue"},
				"fields":  {"teams,id,name"},
			},
			wantHydrated: true,
			wantName:     "Los Angeles Dodgers",
		},
		{
			name:       "200 with minimal team parses cleanly",
			teamID:     108,
			respStatus: 200,
			respBody:   `{"teams":[{"id":108,"name":"Los Angeles Angels"}]}`,
			wantName:   "Los Angeles Angels",
		},
		{
			name:       "200 with empty teams slice maps to ErrNotFound",
			teamID:     9999,
			respStatus: 200,
			respBody:   `{"teams":[]}`,
			wantIs:     ErrNotFound,
			wantErr:    "not found",
		},
		{
			name:       "200 with missing teams key maps to ErrNotFound",
			teamID:     9999,
			respStatus: 200,
			respBody:   `{}`,
			wantIs:     ErrNotFound,
			wantErr:    "not found",
		},
		{
			name:       "404 returns ErrNotFound",
			teamID:     9999,
			respStatus: 404,
			respBody:   `{}`,
			wantIs:     ErrNotFound,
			wantErr:    "not found",
		},
		{
			name:       "5xx is wrapped",
			teamID:     119,
			respStatus: 500,
			respBody:   `oops`,
			wantErr:    "unexpected status 500",
		},
		{
			name:       "malformed JSON is wrapped",
			teamID:     119,
			respStatus: 200,
			respBody:   `not json`,
			wantErr:    "team",
		},
		{
			name:       "network failure is wrapped",
			teamID:     119,
			respStatus: 0,
			wantErr:    "team",
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
			ti, err := client.Team(context.Background(), c.teamID, c.query)

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
			if ti == nil {
				t.Fatal("expected non-nil TeamInfo")
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
			if ti.Name != c.wantName {
				t.Errorf("Name = %q, want %q", ti.Name, c.wantName)
			}
			if !c.wantHydrated {
				return
			}

			// Verify the full field-promotion audit on the hydrated row.
			if ti.ID != 119 || ti.Link != "/api/v1/teams/119" ||
				ti.Season != 2026 || ti.AllStarStatus != "N" || !ti.Active {
				t.Errorf("core fields = id=%d link=%q season=%d allStarStatus=%q active=%v",
					ti.ID, ti.Link, ti.Season, ti.AllStarStatus, ti.Active)
			}
			if ti.TeamCode != "lan" || ti.FileCode != "la" || ti.Abbreviation != "LAD" {
				t.Errorf("codes = teamCode=%q fileCode=%q abbr=%q",
					ti.TeamCode, ti.FileCode, ti.Abbreviation)
			}
			if ti.TeamName != "Dodgers" || ti.LocationName != "Los Angeles" ||
				ti.FirstYearOfPlay != "1884" {
				t.Errorf("name parts = teamName=%q location=%q firstYear=%q",
					ti.TeamName, ti.LocationName, ti.FirstYearOfPlay)
			}
			if ti.ShortName != "LA Dodgers" || ti.FranchiseName != "Los Angeles" ||
				ti.ClubName != "Dodgers" {
				t.Errorf("display names = short=%q franchise=%q club=%q",
					ti.ShortName, ti.FranchiseName, ti.ClubName)
			}
			if ti.Venue.ID != 22 || ti.Venue.Name != "Dodger Stadium" {
				t.Errorf("Venue = %+v", ti.Venue)
			}
			if ti.SpringLeague.ID != 114 || ti.SpringLeague.Name != "Cactus League" ||
				ti.SpringLeague.Abbreviation != "CL" {
				t.Errorf("SpringLeague = %+v", ti.SpringLeague)
			}
			if ti.SpringVenue.ID != 3809 || ti.SpringVenue.Link != "/api/v1/venues/3809" {
				t.Errorf("SpringVenue = %+v", ti.SpringVenue)
			}
			if ti.League.ID != 104 || ti.League.Name != "National League" ||
				ti.League.Abbreviation != "NL" || ti.League.NameShort != "National" ||
				ti.League.SeasonState != "inseason" || !ti.League.HasWildCard ||
				ti.League.HasSplitSeason || ti.League.NumGames != 162 ||
				ti.League.HasPlayoffPoints || ti.League.NumTeams != 15 ||
				ti.League.NumWildcardTeams != 3 || ti.League.Season != "2026" ||
				ti.League.OrgCode != "NL" || ti.League.ConferencesInUse ||
				!ti.League.DivisionsInUse || ti.League.SortOrder != 31 ||
				!ti.League.Active {
				t.Errorf("League = %+v", ti.League)
			}
			if ti.League.Sport.ID != 1 || ti.League.Sport.Link != "/api/v1/sports/1" {
				t.Errorf("League.Sport = %+v", ti.League.Sport)
			}
			if ti.League.SeasonDateInfo.SeasonID != "2026" ||
				ti.League.SeasonDateInfo.RegularSeasonStartDate.Format(seasonDateFmt) !=
					"2026-03-25" {
				t.Errorf("League.SeasonDateInfo = %+v", ti.League.SeasonDateInfo)
			}
			if ti.Division.ID != 203 || ti.Division.Name != "National League West" ||
				ti.Division.NameShort != "NL West" || ti.Division.Abbreviation != "NLW" ||
				ti.Division.NumPlayoffTeams != 1 || !ti.Division.Active {
				t.Errorf("Division = %+v", ti.Division)
			}
			if ti.Sport.ID != 1 || ti.Sport.Code != "mlb" ||
				ti.Sport.Abbreviation != "MLB" || !ti.Sport.ActiveStatus {
				t.Errorf("Sport = %+v", ti.Sport)
			}
		})
	}
}
