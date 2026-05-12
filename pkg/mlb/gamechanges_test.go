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

const gameChangesHappyBody = `{
  "totalItems": 1, "totalEvents": 0, "totalGames": 1, "totalGamesInProgress": 0,
  "dates": [{
    "date": "2024-09-07",
    "games": [{"gamePk": 745455, "gameDate": "2024-09-07T17:35:00Z",
               "status": {"abstractGameState": "Final"}}]
  }]
}`

func TestGameChangesFromGen(t *testing.T) {
	cases := []struct {
		name string
		in   *gen.GameChangesResponse
	}{
		{"nil response", nil},
		{"empty struct", &gen.GameChangesResponse{}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := gameChangesFromGen(c.in)
			if got == nil {
				t.Fatal("returned nil")
			}
		})
	}
}

func TestClient_GameChanges(t *testing.T) {
	cases := []struct {
		name         string
		query        GameChangesQuery
		respStatus   int
		respBody     string
		wantPath     string
		wantQuery    url.Values
		wantErr      string
		wantIs       error
		wantHydrated bool
	}{
		{
			name: "happy path: every filter set",
			query: GameChangesQuery{
				UpdatedSince: "2024-09-07T00:00:00", SportID: 1,
				GameType: "R", Season: 2024, Fields: "dates,games",
			},
			respStatus: 200,
			respBody:   gameChangesHappyBody,
			wantPath:   "/api/v1/game/changes",
			wantQuery: url.Values{
				"updatedSince": {"2024-09-07T00:00:00"}, "sportId": {"1"},
				"gameType": {"R"}, "season": {"2024"}, "fields": {"dates,games"},
			},
			wantHydrated: true,
		},
		{
			name:    "missing UpdatedSince rejected",
			query:   GameChangesQuery{SportID: 1},
			wantErr: "UpdatedSince is required",
			wantIs:  ErrInvalidQuery,
		},
		{
			name:       "200 with empty dates",
			query:      GameChangesQuery{UpdatedSince: "2024-01-01T00:00:00"},
			respStatus: 200,
			respBody:   `{}`,
		},
		{
			name:       "404 returns ErrNotFound",
			query:      GameChangesQuery{UpdatedSince: "2024-01-01T00:00:00"},
			respStatus: 404,
			respBody:   `{}`,
			wantIs:     ErrNotFound,
			wantErr:    "not found",
		},
		{
			name:       "5xx is wrapped",
			query:      GameChangesQuery{UpdatedSince: "2024-01-01T00:00:00"},
			respStatus: 500,
			respBody:   `oops`,
			wantErr:    "unexpected status 500",
		},
		{
			name:       "malformed JSON is wrapped",
			query:      GameChangesQuery{UpdatedSince: "2024-01-01T00:00:00"},
			respStatus: 200,
			respBody:   `not json`,
			wantErr:    "gameChanges",
		},
		{
			name:       "network failure is wrapped",
			query:      GameChangesQuery{UpdatedSince: "2024-01-01T00:00:00"},
			respStatus: 0,
			wantErr:    "gameChanges",
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
			gc, err := client.GameChanges(context.Background(), c.query)

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
			if gc == nil {
				t.Fatal("expected non-nil GameChanges")
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
			if gc.TotalGames != 1 || gc.TotalItems != 1 ||
				gc.TotalGamesInProgress != 0 {
				t.Errorf("totals = %+v", gc)
			}
			if len(gc.Dates) != 1 || gc.Dates[0].Date != "2024-09-07" ||
				len(gc.Dates[0].Games) != 1 || gc.Dates[0].Games[0].GamePk != 745455 {
				t.Errorf("Dates = %+v", gc.Dates)
			}
		})
	}
}
