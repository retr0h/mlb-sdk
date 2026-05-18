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
const draftHappyBody = `{
  "copyright": "Copyright 2026 MLB Advanced Media, L.P.",
  "drafts": {
    "draftYear": 2024,
    "rounds": [{
      "round": "1",
      "picks": [{
        "bisPlayerId": 5027335,
        "pickRound": "1", "pickNumber": 1,
        "displayPickNumber": 1, "roundPickNumber": 1, "rank": 1,
        "pickValue": "10570600", "signingBonus": "8950000",
        "home": {"city": "Sydney", "country": "AUS"},
        "scoutingReport": "https://mlb.com/video/bazzana",
        "school": {"name": "Oregon State", "schoolClass": "4YR JR",
                   "city": "Corvallis", "country": "USA", "state": "OR"},
        "blurb": "Australian star",
        "headshotLink": "https://img.mlbstatic.com/people/683953",
        "person": {"id": 683953, "fullName": "Travis Bazzana",
                   "link": "/api/v1/people/683953"},
        "team": {"id": 114, "name": "Cleveland Guardians",
                 "link": "/api/v1/teams/114"},
        "draftType": {"code": "JR", "description": "Rule 4 / June Amateur Draft"},
        "isDrafted": true, "isPass": false, "year": "2024"
      }]
    }]
  }
}`
//
func TestClient_Draft(t *testing.T) {
	cases := []struct {
		name         string
		year         int
		query        DraftQuery
		respStatus   int
		respBody     string
		wantPath     string
		wantQuery    url.Values
		wantErr      string
		wantIs       error
		wantHydrated bool
	}{
		{
			name: "happy path: 2024 round 1",
			year: 2024,
			query: DraftQuery{
				Round:  "1",
				Fields: "drafts,rounds,picks",
			},
			respStatus: 200,
			respBody:   draftHappyBody,
			wantPath:   "/api/v1/draft/2024",
			wantQuery: url.Values{
				"round":  {"1"},
				"fields": {"drafts,rounds,picks"},
			},
			wantHydrated: true,
		},
		{
			name:       "200 with empty drafts",
			year:       1900,
			respStatus: 200,
			respBody:   `{}`,
		},
		{
			name:       "404 returns ErrNotFound",
			year:       9999,
			respStatus: 404,
			respBody:   `{}`,
			wantIs:     ErrNotFound,
			wantErr:    "not found",
		},
		{
			name:       "5xx is wrapped",
			year:       2024,
			respStatus: 500,
			respBody:   `oops`,
			wantErr:    "unexpected status 500",
		},
		{
			name:       "malformed JSON is wrapped",
			year:       2024,
			respStatus: 200,
			respBody:   `not json`,
			wantErr:    "draft",
		},
		{
			name:       "network failure is wrapped",
			year:       2024,
			respStatus: 0,
			wantErr:    "draft",
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
			d, err := client.Draft(context.Background(), c.year, c.query)
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
			if d == nil {
				t.Fatal("expected non-nil DraftData")
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
			if d.DraftYear != 2024 {
				t.Errorf("DraftYear = %d", d.DraftYear)
			}
			if len(d.Rounds) != 1 || d.Rounds[0].Round != "1" {
				t.Fatalf("Rounds = %+v", d.Rounds)
			}
			if len(d.Rounds[0].Picks) != 1 {
				t.Fatalf("Picks = %d", len(d.Rounds[0].Picks))
			}
			p := d.Rounds[0].Picks[0]
			if p.BisPlayerID != 5027335 || p.PickRound != "1" || p.PickNumber != 1 ||
				p.DisplayPickNumber != 1 || p.RoundPickNumber != 1 || p.Rank != 1 ||
				p.PickValue != "10570600" || p.SigningBonus != "8950000" {
				t.Errorf("pick basics = %+v", p)
			}
			if p.Home.City != "Sydney" || p.Home.Country != "AUS" {
				t.Errorf("Home = %+v", p.Home)
			}
			if p.School.Name != "Oregon State" || p.School.SchoolClass != "4YR JR" ||
				p.School.State != "OR" {
				t.Errorf("School = %+v", p.School)
			}
			if p.Person.ID != 683953 || p.Person.FullName != "Travis Bazzana" {
				t.Errorf("Person = %+v", p.Person)
			}
			if p.Team.ID != 114 || p.Team.Name != "Cleveland Guardians" {
				t.Errorf("Team = %+v", p.Team)
			}
			if p.DraftType.Code != "JR" {
				t.Errorf("DraftType = %+v", p.DraftType)
			}
			if !p.IsDrafted || p.IsPass || p.Year != "2024" {
				t.Errorf("flags = drafted=%v pass=%v year=%q", p.IsDrafted, p.IsPass, p.Year)
			}
			if p.Blurb != "Australian star" || p.ScoutingReport != "https://mlb.com/video/bazzana" {
				t.Errorf("text = blurb=%q report=%q", p.Blurb, p.ScoutingReport)
			}
		})
	}
}
