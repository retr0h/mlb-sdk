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

	"github.com/retr0h/mlb-sdk/internal/gen"
)

const rosterHappyBody = `{
  "copyright": "Copyright 2026 MLB Advanced Media, L.P.",
  "link": "/api/v1/teams/119/roster",
  "roster": [{
    "person": {"id": 605113, "fullName": "Nick Ahmed", "link": "/api/v1/people/605113"},
    "jerseyNumber": "12",
    "position": {"code": "6", "name": "Shortstop", "type": "Infielder", "abbreviation": "SS"},
    "status":   {"code": "MIN", "description": "Minor League Contract"}
  }]
}`

func TestRosterFromGen(t *testing.T) {
	cases := []struct {
		name    string
		in      *gen.RosterResponse
		wantLen int
	}{
		{"nil response yields empty Roster", nil, 0},
		{"empty struct yields empty Roster", &gen.RosterResponse{}, 0},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := rosterFromGen(c.in)
			if got == nil {
				t.Fatal("rosterFromGen returned nil")
			}
			if len(got.Roster) != c.wantLen {
				t.Errorf("len(Roster) = %d, want %d", len(got.Roster), c.wantLen)
			}
		})
	}
}

func TestClient_Roster(t *testing.T) {
	on := time.Date(2024, 7, 15, 0, 0, 0, 0, time.UTC)
	cases := []struct {
		name         string
		teamID       int
		query        RosterQuery
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
			name:   "happy path: every filter set",
			teamID: 119,
			query: RosterQuery{
				RosterType: "active",
				Season:     2024,
				On:         on,
				Hydrate:    "person",
				Fields:     "roster,person",
			},
			respStatus: 200,
			respBody:   rosterHappyBody,
			wantPath:   "/api/v1/teams/119/roster",
			wantQuery: url.Values{
				"rosterType": {"active"},
				"season":     {"2024"},
				"date":       {"2024-07-15"},
				"hydrate":    {"person"},
				"fields":     {"roster,person"},
			},
			wantLen:      1,
			wantHydrated: true,
		},
		{
			name:       "200 with no roster yields empty slice",
			teamID:     119,
			respStatus: 200,
			respBody:   `{}`,
			wantLen:    0,
		},
		{
			name:       "200 with explicit empty array",
			teamID:     119,
			respStatus: 200,
			respBody:   `{"roster": []}`,
			wantLen:    0,
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
			wantErr:    "roster",
		},
		{
			name:       "network failure is wrapped",
			teamID:     119,
			respStatus: 0,
			wantErr:    "roster",
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
			ro, err := client.Roster(context.Background(), c.teamID, c.query)

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
			if ro == nil {
				t.Fatal("expected non-nil Roster")
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
			if got := len(ro.Roster); got != c.wantLen {
				t.Errorf("len(Roster) = %d, want %d", got, c.wantLen)
			}
			if !c.wantHydrated {
				return
			}
			if ro.Link != "/api/v1/teams/119/roster" {
				t.Errorf("Link = %q", ro.Link)
			}
			e := ro.Roster[0]
			if e.Person.ID != 605113 || e.Person.FullName != "Nick Ahmed" {
				t.Errorf("Person = %+v", e.Person)
			}
			if e.JerseyNumber != "12" {
				t.Errorf("JerseyNumber = %q", e.JerseyNumber)
			}
			if e.Position.Abbreviation != "SS" || e.Position.Code != "6" {
				t.Errorf("Position = %+v", e.Position)
			}
			if e.Status.Code != "MIN" ||
				e.Status.Description != "Minor League Contract" {
				t.Errorf("Status = %+v", e.Status)
			}
		})
	}
}
