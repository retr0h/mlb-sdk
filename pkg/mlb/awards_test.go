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
const awardsHappyBody = `{
  "copyright": "Copyright 2026 MLB Advanced Media, L.P.",
  "awards": [
    {
      "id": "MLBHOF", "name": "Hall Of Fame",
      "date": "2026-07-26", "season": "2026",
      "team": {"id": 137, "link": "/api/v1/teams/137"},
      "player": {
        "id": 116999, "link": "/api/v1/people/116999",
        "fullName": "Jeff Kent", "nameFirstLast": "Jeff Kent",
        "primaryPosition": {
          "code": "4", "name": "Second Base",
          "type": "Infielder", "abbreviation": "2B"
        }
      },
      "notes": "Contemporary Era Ballot Selection"
    },
    {
      "id": "MLBHOF", "name": "Hall Of Fame",
      "date": "2026-01-20", "season": "2026",
      "team": {"id": 144, "link": "/api/v1/teams/144"},
      "player": {
        "id": 116662, "link": "/api/v1/people/116662",
        "nameFirstLast": "Andruw Jones",
        "primaryPosition": {
          "code": "O", "name": "Outfield",
          "type": "Outfielder", "abbreviation": "OF"
        }
      }
    }
  ]
}`
//
func TestClient_AwardRecipients(t *testing.T) {
	cases := []struct {
		name         string
		awardID      string
		query        AwardRecipientsQuery
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
			awardID: "MLBHOF",
			query: AwardRecipientsQuery{
				SportID:  1,
				LeagueID: 103,
				Season:   2026,
				Hydrate:  "person.primaryPosition",
				Fields:   "awards,id,player",
			},
			respStatus: 200,
			respBody:   awardsHappyBody,
			wantPath:   "/api/v1/awards/MLBHOF/recipients",
			wantQuery: url.Values{
				"sportId":  {"1"},
				"leagueId": {"103"},
				"season":   {"2026"},
				"hydrate":  {"person.primaryPosition"},
				"fields":   {"awards,id,player"},
			},
			wantLen:      2,
			wantHydrated: true,
		},
		{
			name:       "200 with no awards yields empty slice",
			awardID:    "MLBHOF",
			respStatus: 200,
			respBody:   `{}`,
			wantLen:    0,
		},
		{
			name:       "200 with explicit empty awards array",
			awardID:    "MLBHOF",
			respStatus: 200,
			respBody:   `{"awards": []}`,
			wantLen:    0,
		},
		{
			name:       "404 returns ErrNotFound",
			awardID:    "BOGUS",
			respStatus: 404,
			respBody:   `{}`,
			wantIs:     ErrNotFound,
			wantErr:    "not found",
		},
		{
			name:       "5xx is wrapped",
			awardID:    "MLBHOF",
			respStatus: 500,
			respBody:   `oops`,
			wantErr:    "unexpected status 500",
		},
		{
			name:       "malformed JSON is wrapped",
			awardID:    "MLBHOF",
			respStatus: 200,
			respBody:   `not json`,
			wantErr:    "awardRecipients",
		},
		{
			name:       "network failure is wrapped",
			awardID:    "MLBHOF",
			respStatus: 0,
			wantErr:    "awardRecipients",
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
			a, err := client.AwardRecipients(context.Background(), c.awardID, c.query)
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
			if a == nil {
				t.Fatal("expected non-nil Awards")
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
			if got := len(a.Recipients); got != c.wantLen {
				t.Errorf("len(Recipients) = %d, want %d", got, c.wantLen)
			}
			if !c.wantHydrated {
				return
			}
			r := a.Recipients[0]
			if r.ID != "MLBHOF" || r.Name != "Hall Of Fame" || r.Date != "2026-07-26" ||
				r.Season != "2026" || r.Notes != "Contemporary Era Ballot Selection" {
				t.Errorf("recipient[0] = %+v", r)
			}
			if r.Team.ID != 137 || r.Team.Link != "/api/v1/teams/137" {
				t.Errorf("recipient[0].Team = %+v", r.Team)
			}
			if r.Player.ID != 116999 || r.Player.FullName != "Jeff Kent" ||
				r.Player.NameFirstLast != "Jeff Kent" ||
				r.Player.Link != "/api/v1/people/116999" {
				t.Errorf("recipient[0].Player = %+v", r.Player)
			}
			pp := r.Player.PrimaryPosition
			if pp.Code != "4" || pp.Name != "Second Base" ||
				pp.Type != "Infielder" || pp.Abbreviation != "2B" {
				t.Errorf("recipient[0].Player.PrimaryPosition = %+v", pp)
			}
		})
	}
}
