// Copyright (c) 2026 John Dewey
//
// SPDX-License-Identifier: MIT

package mlb

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestClient_SchedulePostseason(t *testing.T) {
	body := `{"dates":[{"date":"2024-10-01","games":[{"gamePk":1}]}]}`
	cases := []struct {
		name       string
		query      SchedulePostseasonQuery
		respStatus int
		respBody   string
		wantErr    string
		wantIs     error
		wantLen    int
	}{
		{
			name: "happy",
			query: SchedulePostseasonQuery{
				Season:       2024,
				GameTypes:    "D,L,W",
				SeriesNumber: 1,
				TeamID:       119,
				SportID:      1,
				Hydrate:      "team",
				Fields:       "dates",
			},
			respStatus: 200, respBody: body, wantLen: 1,
		},
		{name: "empty", respStatus: 200, respBody: `{}`, wantLen: 0},
		{name: "404", respStatus: 404, respBody: `{}`, wantIs: ErrNotFound, wantErr: "not found"},
		{name: "5xx", respStatus: 500, respBody: `oops`, wantErr: "unexpected status 500"},
		{name: "bad json", respStatus: 200, respBody: `x`, wantErr: "schedulePostseason"},
		{name: "network", respStatus: 0, wantErr: "schedulePostseason"},
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
			u := srv.URL
			if c.respStatus == 0 {
				srv.Close()
			} else {
				defer srv.Close()
			}
			games, err := New(WithBaseURL(u)).SchedulePostseason(context.Background(), c.query)
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
			if len(games) != c.wantLen {
				t.Errorf("len = %d, want %d", len(games), c.wantLen)
			}
		})
	}
}

func TestClient_SchedulePostseasonTuneIn(t *testing.T) {
	body := `{"dates":[{"date":"2024-10-01","games":[{"gamePk":2}]}]}`
	cases := []struct {
		name       string
		query      SchedulePostseasonTuneInQuery
		respStatus int
		respBody   string
		wantErr    string
		wantIs     error
		wantLen    int
	}{
		{
			name: "happy", query: SchedulePostseasonTuneInQuery{Season: 2024, TeamID: 119, SportID: 1, Hydrate: "team", Fields: "dates"},
			respStatus: 200, respBody: body, wantLen: 1,
		},
		{name: "empty", respStatus: 200, respBody: `{}`, wantLen: 0},
		{name: "404", respStatus: 404, respBody: `{}`, wantIs: ErrNotFound, wantErr: "not found"},
		{name: "5xx", respStatus: 500, respBody: `oops`, wantErr: "unexpected status 500"},
		{name: "bad json", respStatus: 200, respBody: `x`, wantErr: "schedulePostseasonTuneIn"},
		{name: "network", respStatus: 0, wantErr: "schedulePostseasonTuneIn"},
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
			u := srv.URL
			if c.respStatus == 0 {
				srv.Close()
			} else {
				defer srv.Close()
			}
			games, err := New(
				WithBaseURL(u),
			).SchedulePostseasonTuneIn(context.Background(), c.query)
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
			if len(games) != c.wantLen {
				t.Errorf("len = %d, want %d", len(games), c.wantLen)
			}
		})
	}
}
