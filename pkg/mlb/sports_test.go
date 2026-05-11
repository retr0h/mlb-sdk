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

const sportsHappyBody = `{
  "copyright": "Copyright 2026 MLB Advanced Media, L.P.",
  "sports": [
    {
      "id": 1, "code": "mlb", "link": "/api/v1/sports/1",
      "name": "Major League Baseball", "abbreviation": "MLB",
      "sortOrder": 11, "activeStatus": true
    },
    {
      "id": 11, "code": "aaa", "link": "/api/v1/sports/11",
      "name": "Triple-A", "abbreviation": "AAA",
      "sortOrder": 101, "activeStatus": true
    }
  ]
}`

func TestSports_Sport(t *testing.T) {
	empty := sportsFromGen(nil)
	full := &Sports{Sports: []Sport{
		{ID: 1, Name: "MLB"},
		{ID: 11, Name: "AAA"},
	}}
	cases := []struct {
		name     string
		s        *Sports
		lookup   int
		wantOk   bool
		wantName string
	}{
		{"nil receiver", nil, 1, false, ""},
		{"empty (no sports)", empty, 1, false, ""},
		{"miss", full, 999, false, ""},
		{"hit MLB", full, 1, true, "MLB"},
		{"hit AAA", full, 11, true, "AAA"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := c.s.Sport(c.lookup)
			if (got != nil) != c.wantOk {
				t.Fatalf("Sport(%d) ok=%v, want %v", c.lookup, got != nil, c.wantOk)
			}
			if got != nil && got.Name != c.wantName {
				t.Errorf("Sport(%d).Name = %q, want %q", c.lookup, got.Name, c.wantName)
			}
		})
	}
}

func TestClient_Sports(t *testing.T) {
	cases := []struct {
		name       string
		query      SportsQuery
		respStatus int
		respBody   string
		wantPath   string
		wantQuery  url.Values
		wantErr    string
		wantIs     error
		wantLen    int
		wantFirst  *Sport
	}{
		{
			name: "happy path: MLB + AAA",
			query: SportsQuery{
				SportID: 1,
				Fields:  "sports,id,name",
			},
			respStatus: 200,
			respBody:   sportsHappyBody,
			wantPath:   "/api/v1/sports",
			wantQuery: url.Values{
				"sportId": {"1"},
				"fields":  {"sports,id,name"},
			},
			wantLen: 2,
			wantFirst: &Sport{
				ID: 1, Code: "mlb", Link: "/api/v1/sports/1",
				Name: "Major League Baseball", Abbreviation: "MLB",
				SortOrder: 11, ActiveStatus: true,
			},
		},
		{
			name:       "200 with no sports yields empty slice",
			respStatus: 200,
			respBody:   `{}`,
			wantLen:    0,
		},
		{
			name:       "200 with explicit empty sports array",
			respStatus: 200,
			respBody:   `{"sports": []}`,
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
			wantErr:    "sports",
		},
		{
			name:       "network failure is wrapped",
			respStatus: 0,
			wantErr:    "sports",
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
			s, err := client.Sports(context.Background(), c.query)

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
			if s == nil {
				t.Fatal("expected non-nil Sports")
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
			if got := len(s.Sports); got != c.wantLen {
				t.Errorf("len(Sports) = %d, want %d", got, c.wantLen)
			}
			if c.wantFirst != nil {
				if got := s.Sports[0]; got != *c.wantFirst {
					t.Errorf("Sports[0]:\n got = %+v\nwant = %+v", got, *c.wantFirst)
				}
			}
		})
	}
}
