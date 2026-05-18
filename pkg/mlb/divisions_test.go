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
const divisionsHappyBody = `{
  "copyright": "Copyright 2026 MLB Advanced Media, L.P.",
  "divisions": [
    {
      "id": 200, "name": "American League West", "season": "2026", "nameShort": "AL West",
      "link": "/api/v1/divisions/200", "abbreviation": "ALW",
      "league": {"id": 103, "link": "/api/v1/league/103"},
      "sport":  {"id": 1,   "link": "/api/v1/sports/1"},
      "hasWildcard": false, "sortOrder": 24, "numPlayoffTeams": 1, "active": true
    },
    {
      "id": 204, "name": "National League East", "season": "2026", "nameShort": "NL East",
      "link": "/api/v1/divisions/204", "abbreviation": "NLE",
      "league": {"id": 104, "link": "/api/v1/league/104"},
      "sport":  {"id": 1,   "link": "/api/v1/sports/1"},
      "hasWildcard": true, "sortOrder": 30, "numPlayoffTeams": 1, "active": true
    }
  ]
}`
//
func TestDivisions_Division(t *testing.T) {
	empty := divisionsFromGen(nil) // nil → empty Divisions
	full := &Divisions{Divisions: []Division{
		{ID: 200, Name: "AL West"},
		{ID: 204, Name: "NL East"},
	}}
	cases := []struct {
		name     string
		div      *Divisions
		lookup   int
		wantOk   bool
		wantName string
	}{
		{"nil receiver", nil, 200, false, ""},
		{"empty (no divisions)", empty, 200, false, ""},
		{"miss", full, 999, false, ""},
		{"hit AL West", full, 200, true, "AL West"},
		{"hit NL East", full, 204, true, "NL East"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := c.div.Division(c.lookup)
			if (got != nil) != c.wantOk {
				t.Fatalf("Division(%d) ok=%v, want %v", c.lookup, got != nil, c.wantOk)
			}
			if got != nil && got.Name != c.wantName {
				t.Errorf("Division(%d).Name = %q, want %q", c.lookup, got.Name, c.wantName)
			}
		})
	}
}
//
func TestClient_Divisions(t *testing.T) {
	cases := []struct {
		name       string
		query      DivisionsQuery
		respStatus int
		respBody   string
		wantPath   string
		wantQuery  url.Values
		wantErr    string
		wantIs     error
		wantLen    int
		wantFirst  *Division
	}{
		{
			name: "happy path: AL West + NL East",
			query: DivisionsQuery{
				DivisionID: 200,
				LeagueID:   103,
				SportID:    1,
				Season:     2026,
			},
			respStatus: 200,
			respBody:   divisionsHappyBody,
			wantPath:   "/api/v1/divisions",
			wantQuery: url.Values{
				"divisionId": {"200"},
				"leagueId":   {"103"},
				"sportId":    {"1"},
				"season":     {"2026"},
			},
			wantLen: 2,
			wantFirst: &Division{
				ID: 200, Name: "American League West", Season: "2026",
				NameShort: "AL West", Link: "/api/v1/divisions/200",
				Abbreviation: "ALW",
				League:       Ref{ID: 103, Link: "/api/v1/league/103"},
				Sport:        Ref{ID: 1, Link: "/api/v1/sports/1"},
				HasWildcard:  false, SortOrder: 24, NumPlayoffTeams: 1, Active: true,
			},
		},
		{
			name:       "200 with no divisions yields empty slice",
			respStatus: 200,
			respBody:   `{}`,
			wantLen:    0,
		},
		{
			name:       "200 with explicit empty divisions array",
			respStatus: 200,
			respBody:   `{"divisions": []}`,
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
			wantErr:    "divisions",
		},
		{
			name:       "network failure is wrapped",
			respStatus: 0,
			wantErr:    "divisions",
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
			d, err := client.Divisions(context.Background(), c.query)
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
				t.Fatal("expected non-nil Divisions")
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
			if got := len(d.Divisions); got != c.wantLen {
				t.Errorf("len(Divisions) = %d, want %d", got, c.wantLen)
			}
			if c.wantFirst != nil {
				if got := d.Divisions[0]; got != *c.wantFirst {
					t.Errorf("Divisions[0]:\n got = %+v\nwant = %+v", got, *c.wantFirst)
				}
			}
		})
	}
}
