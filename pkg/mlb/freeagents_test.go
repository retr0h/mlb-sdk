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
const freeAgentsHappyBody = `{
  "freeAgents": [{
    "player": {"id": 642715, "fullName": "Willy Adames", "link": "/api/v1/people/642715"},
    "originalTeam": {"id": 158, "name": "Milwaukee Brewers"},
    "newTeam":      {"id": 137, "name": "San Francisco Giants"},
    "notes": "Seven-Year Contract",
    "dateSigned":   "2024-12-10",
    "dateDeclared": "2024-10-31",
    "position": {"code": "6", "name": "Shortstop", "type": "Infielder", "abbreviation": "SS"}
  }]
}`
//
func TestClient_FreeAgents(t *testing.T) {
	cases := []struct {
		name         string
		query        FreeAgentsQuery
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
			query: FreeAgentsQuery{
				Season:  2024,
				Order:   "asc",
				Hydrate: "person",
				Fields:  "freeAgents,player",
			},
			respStatus: 200,
			respBody:   freeAgentsHappyBody,
			wantPath:   "/api/v1/people/freeAgents",
			wantQuery: url.Values{
				"season":  {"2024"},
				"order":   {"asc"},
				"hydrate": {"person"},
				"fields":  {"freeAgents,player"},
			},
			wantLen:      1,
			wantHydrated: true,
		},
		{
			name:       "200 with no freeAgents yields empty slice",
			respStatus: 200,
			respBody:   `{}`,
			wantLen:    0,
		},
		{
			name:       "200 with explicit empty array",
			respStatus: 200,
			respBody:   `{"freeAgents": []}`,
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
			wantErr:    "freeAgents",
		},
		{
			name:       "network failure is wrapped",
			respStatus: 0,
			wantErr:    "freeAgents",
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
			fa, err := client.FreeAgents(context.Background(), c.query)
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
			if fa == nil {
				t.Fatal("expected non-nil FreeAgents")
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
			if got := len(fa.FreeAgents); got != c.wantLen {
				t.Errorf("len(FreeAgents) = %d, want %d", got, c.wantLen)
			}
			if !c.wantHydrated {
				return
			}
			f := fa.FreeAgents[0]
			if f.Player.ID != 642715 || f.Player.FullName != "Willy Adames" {
				t.Errorf("Player = %+v", f.Player)
			}
			if f.OriginalTeam.ID != TeamID(158) ||
				f.OriginalTeam.Name != "Milwaukee Brewers" {
				t.Errorf("OriginalTeam = %+v", f.OriginalTeam)
			}
			if f.NewTeam.ID != TeamID(137) ||
				f.NewTeam.Name != "San Francisco Giants" {
				t.Errorf("NewTeam = %+v", f.NewTeam)
			}
			if f.Notes != "Seven-Year Contract" || f.DateSigned != "2024-12-10" ||
				f.DateDeclared != "2024-10-31" {
				t.Errorf("metadata = notes=%q signed=%q declared=%q",
					f.Notes, f.DateSigned, f.DateDeclared)
			}
			if f.Position.Abbreviation != "SS" || f.Position.Name != "Shortstop" {
				t.Errorf("Position = %+v", f.Position)
			}
		})
	}
}
