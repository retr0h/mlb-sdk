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
	"strings"
	"testing"
)
//
func TestClient_SportsPlayers(t *testing.T) {
	cases := []struct {
		name       string
		sportID    int
		query      SportsPlayersQuery
		respStatus int
		respBody   string
		wantErr    string
		wantIs     error
		wantLen    int
	}{
		{
			name: "happy", sportID: 1, query: SportsPlayersQuery{Season: 2024, GameType: "R", Fields: "people"},
			respStatus: 200, respBody: `{"people":[{"id":1,"fullName":"Test"}]}`, wantLen: 1,
		},
		{
			name: "missing Season", sportID: 1, query: SportsPlayersQuery{},
			wantErr: "Season is required", wantIs: ErrInvalidQuery,
		},
		{
			name: "empty", sportID: 1, query: SportsPlayersQuery{Season: 2024},
			respStatus: 200, respBody: `{}`, wantLen: 0,
		},
		{
			name: "404", sportID: 1, query: SportsPlayersQuery{Season: 2024},
			respStatus: 404, respBody: `{}`, wantIs: ErrNotFound, wantErr: "not found",
		},
		{
			name: "5xx", sportID: 1, query: SportsPlayersQuery{Season: 2024},
			respStatus: 500, respBody: `oops`, wantErr: "unexpected status 500",
		},
		{
			name: "bad json", sportID: 1, query: SportsPlayersQuery{Season: 2024},
			respStatus: 200, respBody: `x`, wantErr: "sportsPlayers",
		},
		{
			name: "network", sportID: 1, query: SportsPlayersQuery{Season: 2024},
			respStatus: 0, wantErr: "sportsPlayers",
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
			pp, err := New(WithBaseURL(u)).SportsPlayers(context.Background(), c.sportID, c.query)
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
