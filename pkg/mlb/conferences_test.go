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

const conferencesHappyBody = `{
  "conferences": [{
    "id": 301, "name": "PCL American Conference",
    "link": "/api/v1/conferences/301", "abbreviation": "PCLA",
    "nameShort": "PCL American", "hasWildcard": false,
    "league": {"id": 112, "link": "/api/v1/league/112"},
    "sport":  {"id": 11,  "link": "/api/v1/sports/11"}
  }]
}`

func TestConferences_Conference(t *testing.T) {
	empty := conferencesFromGen(nil)
	full := &Conferences{Conferences: []Conference{
		{ID: 301, Name: "PCL American"},
	}}
	cases := []struct {
		name   string
		c      *Conferences
		lookup int
		wantOk bool
	}{
		{"nil receiver", nil, 301, false},
		{"empty", empty, 301, false},
		{"miss", full, 999, false},
		{"hit", full, 301, true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := c.c.Conference(c.lookup)
			if (got != nil) != c.wantOk {
				t.Fatalf("Conference(%d) ok=%v, want %v", c.lookup, got != nil, c.wantOk)
			}
		})
	}
}

func TestClient_Conferences(t *testing.T) {
	cases := []struct {
		name       string
		query      ConferencesQuery
		respStatus int
		respBody   string
		wantPath   string
		wantQuery  url.Values
		wantErr    string
		wantIs     error
		wantLen    int
		wantFirst  *Conference
	}{
		{
			name: "happy path: every filter set",
			query: ConferencesQuery{
				ConferenceID: 301,
				Season:       2024,
				Fields:       "conferences,id",
			},
			respStatus: 200,
			respBody:   conferencesHappyBody,
			wantPath:   "/api/v1/conferences",
			wantQuery: url.Values{
				"conferenceId": {"301"},
				"season":       {"2024"},
				"fields":       {"conferences,id"},
			},
			wantLen: 1,
			wantFirst: &Conference{
				ID: 301, Name: "PCL American Conference",
				Link: "/api/v1/conferences/301", Abbreviation: "PCLA",
				NameShort: "PCL American", HasWildcard: false,
				League: Ref{ID: 112, Link: "/api/v1/league/112"},
				Sport:  Ref{ID: 11, Link: "/api/v1/sports/11"},
			},
		},
		{
			name:       "200 with no conferences",
			respStatus: 200,
			respBody:   `{}`,
			wantLen:    0,
		},
		{
			name:       "200 with empty array",
			respStatus: 200,
			respBody:   `{"conferences": []}`,
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
			wantErr:    "conferences",
		},
		{
			name:       "network failure is wrapped",
			respStatus: 0,
			wantErr:    "conferences",
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
			co, err := client.Conferences(context.Background(), c.query)

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
			if co == nil {
				t.Fatal("expected non-nil Conferences")
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
			if got := len(co.Conferences); got != c.wantLen {
				t.Errorf("len(Conferences) = %d, want %d", got, c.wantLen)
			}
			if c.wantFirst != nil {
				if got := co.Conferences[0]; got != *c.wantFirst {
					t.Errorf("Conferences[0]:\n got = %+v\nwant = %+v", got, *c.wantFirst)
				}
			}
		})
	}
}
