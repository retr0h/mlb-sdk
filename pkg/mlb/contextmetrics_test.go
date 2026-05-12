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

	"github.com/retr0h/mlb-sdk/internal/gen"
)

const contextMetricsHappyBody = `{
  "game": {"id": 745455, "link": "/api/v1.1/game/745455/feed/live"},
  "leftFieldSacFlyProbability": {},
  "centerFieldSacFlyProbability": {},
  "rightFieldSacFlyProbability": {},
  "homeWinProbability": 0.0,
  "awayWinProbability": 100.0
}`

func TestContextMetricsFromGen(t *testing.T) {
	cases := []struct {
		name string
		in   *gen.ContextMetricsResponse
	}{
		{"nil response", nil},
		{"empty struct", &gen.ContextMetricsResponse{}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := contextMetricsFromGen(c.in)
			if got == nil {
				t.Fatal("returned nil")
			}
		})
	}
}

func TestClient_ContextMetrics(t *testing.T) {
	cases := []struct {
		name         string
		gamePk       int
		query        ContextMetricsQuery
		respStatus   int
		respBody     string
		wantPath     string
		wantQuery    url.Values
		wantErr      string
		wantIs       error
		wantHydrated bool
	}{
		{
			name:       "happy path",
			gamePk:     745455,
			query:      ContextMetricsQuery{Timecode: "20240907_180000", Fields: "game"},
			respStatus: 200,
			respBody:   contextMetricsHappyBody,
			wantPath:   "/api/v1/game/745455/contextMetrics",
			wantQuery: url.Values{
				"timecode": {"20240907_180000"},
				"fields":   {"game"},
			},
			wantHydrated: true,
		},
		{
			name:       "200 with minimal body",
			gamePk:     1,
			respStatus: 200,
			respBody:   `{"homeWinProbability": 50.0}`,
		},
		{
			name:       "404 returns ErrNotFound",
			gamePk:     9999,
			respStatus: 404,
			respBody:   `{}`,
			wantIs:     ErrNotFound,
			wantErr:    "not found",
		},
		{
			name:       "5xx is wrapped",
			gamePk:     745455,
			respStatus: 500,
			respBody:   `oops`,
			wantErr:    "unexpected status 500",
		},
		{
			name:       "malformed JSON is wrapped",
			gamePk:     745455,
			respStatus: 200,
			respBody:   `not json`,
			wantErr:    "contextMetrics",
		},
		{
			name:       "network failure is wrapped",
			gamePk:     745455,
			respStatus: 0,
			wantErr:    "contextMetrics",
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
			cm, err := client.ContextMetrics(context.Background(), c.gamePk, c.query)

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
			if cm == nil {
				t.Fatal("expected non-nil ContextMetrics")
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
			if cm.Game.Link != "/api/v1.1/game/745455/feed/live" {
				t.Errorf("Game = %+v", cm.Game)
			}
			if cm.HomeWinProbability != 0.0 || cm.AwayWinProbability != 100.0 {
				t.Errorf("prob = home=%v away=%v",
					cm.HomeWinProbability, cm.AwayWinProbability)
			}
		})
	}
}
