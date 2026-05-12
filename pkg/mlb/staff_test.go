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

const staffHappyBody = `{
  "roster": [{
    "person": {"id": 150353, "fullName": "Dave Roberts", "link": "/api/v1/people/150353"},
    "jerseyNumber": "30", "job": "Manager", "jobId": "MNGR", "title": "Manager"
  }],
  "link": "/api/v1/teams/119/coaches", "teamId": 119, "rosterType": "coach"
}`

func TestStaffFromGen(t *testing.T) {
	cases := []struct {
		name string
		in   *gen.StaffResponse
	}{
		{"nil response", nil},
		{"empty struct", &gen.StaffResponse{}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := staffFromGen(c.in)
			if got == nil {
				t.Fatal("returned nil")
			}
		})
	}
}

func TestClient_Coaches(t *testing.T) {
	on := time.Date(2024, 7, 15, 0, 0, 0, 0, time.UTC)
	cases := []struct {
		name         string
		teamID       int
		query        CoachesQuery
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
			name: "happy path", teamID: 119,
			query:      CoachesQuery{Season: 2024, On: on, Fields: "roster"},
			respStatus: 200, respBody: staffHappyBody,
			wantPath: "/api/v1/teams/119/coaches",
			wantQuery: url.Values{
				"season": {"2024"}, "date": {"2024-07-15"}, "fields": {"roster"},
			},
			wantLen: 1, wantHydrated: true,
		},
		{name: "empty", teamID: 119, respStatus: 200, respBody: `{}`, wantLen: 0},
		{
			name: "404", teamID: 9999, respStatus: 404, respBody: `{}`,
			wantIs: ErrNotFound, wantErr: "not found",
		},
		{
			name: "5xx", teamID: 119, respStatus: 500, respBody: `oops`,
			wantErr: "unexpected status 500",
		},
		{
			name: "bad json", teamID: 119, respStatus: 200, respBody: `x`,
			wantErr: "coaches",
		},
		{name: "network", teamID: 119, respStatus: 0, wantErr: "coaches"},
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
			s, err := client.Coaches(context.Background(), c.teamID, c.query)
			if c.wantErr != "" {
				if err == nil {
					t.Fatalf("expected error containing %q", c.wantErr)
				}
				if !strings.Contains(err.Error(), c.wantErr) {
					t.Errorf("err = %v, want %q", err, c.wantErr)
				}
				if c.wantIs != nil && !errors.Is(err, c.wantIs) {
					t.Errorf("errors.Is = false")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected: %v", err)
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
			if len(s.Roster) != c.wantLen {
				t.Errorf("len = %d, want %d", len(s.Roster), c.wantLen)
			}
			if c.wantHydrated {
				e := s.Roster[0]
				if e.Person.FullName != "Dave Roberts" || e.Job != "Manager" ||
					e.JobID != "MNGR" || e.JerseyNumber != "30" {
					t.Errorf("entry = %+v", e)
				}
				if s.TeamID != 119 || s.RosterType != "coach" {
					t.Errorf("meta = teamId=%d type=%q", s.TeamID, s.RosterType)
				}
			}
		})
	}
}

func TestClient_Personnel(t *testing.T) {
	cases := []struct {
		name       string
		teamID     int
		query      PersonnelQuery
		respStatus int
		respBody   string
		wantErr    string
		wantIs     error
		wantLen    int
	}{
		{name: "happy", teamID: 119, respStatus: 200, respBody: staffHappyBody, wantLen: 1},
		{
			name: "with date filter", teamID: 119,
			query: PersonnelQuery{
				On:     time.Date(2024, 7, 15, 0, 0, 0, 0, time.UTC),
				Fields: "roster",
			},
			respStatus: 200, respBody: staffHappyBody, wantLen: 1,
		},
		{name: "empty", teamID: 119, respStatus: 200, respBody: `{}`, wantLen: 0},
		{
			name: "404", teamID: 9999, respStatus: 404, respBody: `{}`,
			wantIs: ErrNotFound, wantErr: "not found",
		},
		{
			name: "5xx", teamID: 119, respStatus: 500, respBody: `oops`,
			wantErr: "unexpected status 500",
		},
		{
			name: "bad json", teamID: 119, respStatus: 200, respBody: `x`,
			wantErr: "personnel",
		},
		{name: "network", teamID: 119, respStatus: 0, wantErr: "personnel"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			srv := httptest.NewServer(
				http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
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
			s, err := client.Personnel(context.Background(), c.teamID, c.query)
			if c.wantErr != "" {
				if err == nil {
					t.Fatalf("expected error containing %q", c.wantErr)
				}
				if !strings.Contains(err.Error(), c.wantErr) {
					t.Errorf("err = %v", err)
				}
				if c.wantIs != nil && !errors.Is(err, c.wantIs) {
					t.Errorf("errors.Is = false")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected: %v", err)
			}
			if len(s.Roster) != c.wantLen {
				t.Errorf("len = %d, want %d", len(s.Roster), c.wantLen)
			}
		})
	}
}

func TestClient_Umpires(t *testing.T) {
	cases := []struct {
		name       string
		query      UmpiresQuery
		respStatus int
		respBody   string
		wantPath   string
		wantQuery  url.Values
		wantErr    string
		wantIs     error
		wantLen    int
	}{
		{
			name:       "happy",
			query:      UmpiresQuery{SportID: 1, Fields: "roster"},
			respStatus: 200, respBody: staffHappyBody,
			wantPath:  "/api/v1/jobs/umpires",
			wantQuery: url.Values{"sportId": {"1"}, "fields": {"roster"}},
			wantLen:   1,
		},
		{
			name:       "with date filter",
			query:      UmpiresQuery{On: time.Date(2024, 7, 15, 0, 0, 0, 0, time.UTC)},
			respStatus: 200, respBody: staffHappyBody, wantLen: 1,
		},
		{name: "empty", respStatus: 200, respBody: `{}`, wantLen: 0},
		{
			name: "404", respStatus: 404, respBody: `{}`,
			wantIs: ErrNotFound, wantErr: "not found",
		},
		{
			name: "5xx", respStatus: 500, respBody: `oops`,
			wantErr: "unexpected status 500",
		},
		{name: "bad json", respStatus: 200, respBody: `x`, wantErr: "umpires"},
		{name: "network", respStatus: 0, wantErr: "umpires"},
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
			s, err := client.Umpires(context.Background(), c.query)
			if c.wantErr != "" {
				if err == nil {
					t.Fatalf("expected error containing %q", c.wantErr)
				}
				if !strings.Contains(err.Error(), c.wantErr) {
					t.Errorf("err = %v", err)
				}
				if c.wantIs != nil && !errors.Is(err, c.wantIs) {
					t.Errorf("errors.Is = false")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected: %v", err)
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
			if len(s.Roster) != c.wantLen {
				t.Errorf("len = %d, want %d", len(s.Roster), c.wantLen)
			}
		})
	}
}
