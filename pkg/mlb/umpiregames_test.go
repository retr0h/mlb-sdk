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

func TestClient_UmpireGames(t *testing.T) {
	body := `{"dates":[{"date":"2024-09-07","games":[{"gamePk":745455}]}]}`
	cases := []struct {
		name       string
		umpireID   int
		query      UmpireGamesQuery
		respStatus int
		respBody   string
		wantErr    string
		wantIs     error
		wantLen    int
	}{
		{
			name: "happy", umpireID: 596809, query: UmpireGamesQuery{Season: 2024, Fields: "dates"},
			respStatus: 200, respBody: body, wantLen: 1,
		},
		{
			name: "missing Season", umpireID: 596809, query: UmpireGamesQuery{},
			wantErr: "Season is required", wantIs: ErrInvalidQuery,
		},
		{
			name: "empty", umpireID: 596809, query: UmpireGamesQuery{Season: 2024},
			respStatus: 200, respBody: `{}`, wantLen: 0,
		},
		{
			name: "404", umpireID: 9999, query: UmpireGamesQuery{Season: 2024},
			respStatus: 404, respBody: `{}`, wantIs: ErrNotFound, wantErr: "not found",
		},
		{
			name: "5xx", umpireID: 1, query: UmpireGamesQuery{Season: 2024},
			respStatus: 500, respBody: `oops`, wantErr: "unexpected status 500",
		},
		{
			name: "bad json", umpireID: 1, query: UmpireGamesQuery{Season: 2024},
			respStatus: 200, respBody: `x`, wantErr: "umpireGames",
		},
		{
			name: "network", umpireID: 1, query: UmpireGamesQuery{Season: 2024},
			respStatus: 0, wantErr: "umpireGames",
		},
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
			games, err := New(WithBaseURL(u)).UmpireGames(context.Background(), c.umpireID, c.query)
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
