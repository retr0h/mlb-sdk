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

func TestClient_Jobs(t *testing.T) {
	cases := []struct {
		name       string
		query      JobsQuery
		respStatus int
		respBody   string
		wantPath   string
		wantQuery  url.Values
		wantErr    string
		wantIs     error
		wantLen    int
	}{
		{
			name: "happy path",
			query: JobsQuery{
				JobType: "UMPR", SportID: 1, Date: "2024-09-07",
				Fields: "roster",
			},
			respStatus: 200, respBody: staffHappyBody,
			wantPath: "/api/v1/jobs",
			wantQuery: url.Values{
				"jobType": {"UMPR"}, "sportId": {"1"},
				"date": {"2024-09-07"}, "fields": {"roster"},
			},
			wantLen: 1,
		},
		{
			name: "missing JobType", query: JobsQuery{},
			wantErr: "JobType is required", wantIs: ErrInvalidQuery,
		},
		{
			name: "empty", query: JobsQuery{JobType: "UMPR"},
			respStatus: 200, respBody: `{}`, wantLen: 0,
		},
		{
			name: "404", query: JobsQuery{JobType: "UMPR"},
			respStatus: 404, respBody: `{}`, wantIs: ErrNotFound, wantErr: "not found",
		},
		{
			name: "5xx", query: JobsQuery{JobType: "UMPR"},
			respStatus: 500, respBody: `oops`, wantErr: "unexpected status 500",
		},
		{
			name: "bad json", query: JobsQuery{JobType: "UMPR"},
			respStatus: 200, respBody: `x`, wantErr: "jobs",
		},
		{
			name: "network", query: JobsQuery{JobType: "UMPR"},
			respStatus: 0, wantErr: "jobs",
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
			s, err := client.Jobs(context.Background(), c.query)
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

func TestClient_Datacasters(t *testing.T) {
	cases := []struct {
		name       string
		query      DatacastersQuery
		respStatus int
		respBody   string
		wantErr    string
		wantIs     error
		wantLen    int
	}{
		{
			name: "happy", query: DatacastersQuery{SportID: 1, Date: "2024-09-07", Fields: "roster"},
			respStatus: 200, respBody: staffHappyBody, wantLen: 1,
		},
		{name: "empty", respStatus: 200, respBody: `{}`, wantLen: 0},
		{name: "404", respStatus: 404, respBody: `{}`, wantIs: ErrNotFound, wantErr: "not found"},
		{name: "5xx", respStatus: 500, respBody: `oops`, wantErr: "unexpected status 500"},
		{name: "bad json", respStatus: 200, respBody: `x`, wantErr: "datacasters"},
		{name: "network", respStatus: 0, wantErr: "datacasters"},
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
			s, err := client.Datacasters(context.Background(), c.query)
			if c.wantErr != "" {
				if err == nil {
					t.Fatalf("expected %q", c.wantErr)
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

func TestClient_OfficialScorers(t *testing.T) {
	cases := []struct {
		name       string
		query      OfficialScorersQuery
		respStatus int
		respBody   string
		wantErr    string
		wantIs     error
		wantLen    int
	}{
		{
			name: "happy", query: OfficialScorersQuery{Timecode: "20240907_180000", Fields: "roster"},
			respStatus: 200, respBody: staffHappyBody, wantLen: 1,
		},
		{name: "empty", respStatus: 200, respBody: `{}`, wantLen: 0},
		{name: "404", respStatus: 404, respBody: `{}`, wantIs: ErrNotFound, wantErr: "not found"},
		{name: "5xx", respStatus: 500, respBody: `oops`, wantErr: "unexpected status 500"},
		{name: "bad json", respStatus: 200, respBody: `x`, wantErr: "officialScorers"},
		{name: "network", respStatus: 0, wantErr: "officialScorers"},
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
			s, err := client.OfficialScorers(context.Background(), c.query)
			if c.wantErr != "" {
				if err == nil {
					t.Fatalf("expected %q", c.wantErr)
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

func TestClient_PeopleChanges(t *testing.T) {
	cases := []struct {
		name       string
		query      PeopleChangesQuery
		respStatus int
		respBody   string
		wantErr    string
		wantIs     error
		wantLen    int
	}{
		{
			name: "happy", query: PeopleChangesQuery{UpdatedSince: "2024-09-07T00:00:00", Fields: "people"},
			respStatus: 200, respBody: `{"people":[{"id":1,"fullName":"Test"}]}`, wantLen: 1,
		},
		{name: "empty", respStatus: 200, respBody: `{}`, wantLen: 0},
		{name: "404", respStatus: 404, respBody: `{}`, wantIs: ErrNotFound, wantErr: "not found"},
		{name: "5xx", respStatus: 500, respBody: `oops`, wantErr: "unexpected status 500"},
		{name: "bad json", respStatus: 200, respBody: `x`, wantErr: "peopleChanges"},
		{name: "network", respStatus: 0, wantErr: "peopleChanges"},
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
			pp, err := client.PeopleChanges(context.Background(), c.query)
			if c.wantErr != "" {
				if err == nil {
					t.Fatalf("expected %q", c.wantErr)
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
			if len(pp) != c.wantLen {
				t.Errorf("len = %d, want %d", len(pp), c.wantLen)
			}
		})
	}
}
