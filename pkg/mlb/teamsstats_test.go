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
func TestClient_TeamsStats(t *testing.T) {
	body := `{"stats":[{"type":{"displayName":"season"},"group":{"displayName":"hitting"},"splits":[]}]}`
	cases := []struct {
		name       string
		query      TeamsStatsQuery
		respStatus int
		respBody   string
		wantErr    string
		wantIs     error
	}{
		{
			name: "happy",
			query: TeamsStatsQuery{
				Season: 2024, Group: "hitting", Stats: "season", SportIDs: "1", GameType: "R",
				Order: "desc", SortStat: "homeRuns", StartDate: "2024-04-01", EndDate: "2024-09-30", Fields: "stats",
			},
			respStatus: 200, respBody: body,
		},
		{
			name: "missing required", query: TeamsStatsQuery{Season: 2024, Group: "hitting"},
			wantErr: "Season, Group, and Stats are all required", wantIs: ErrInvalidQuery,
		},
		{
			name: "404", query: TeamsStatsQuery{Season: 2024, Group: "hitting", Stats: "season"},
			respStatus: 404, respBody: `{}`, wantIs: ErrNotFound, wantErr: "not found",
		},
		{
			name: "5xx", query: TeamsStatsQuery{Season: 2024, Group: "hitting", Stats: "season"},
			respStatus: 500, respBody: `oops`, wantErr: "unexpected status 500",
		},
		{
			name: "bad json", query: TeamsStatsQuery{Season: 2024, Group: "hitting", Stats: "season"},
			respStatus: 200, respBody: `x`, wantErr: "teamsStats",
		},
		{
			name: "network", query: TeamsStatsQuery{Season: 2024, Group: "hitting", Stats: "season"},
			respStatus: 0, wantErr: "teamsStats",
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
			_, err := New(WithBaseURL(u)).TeamsStats(context.Background(), c.query)
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
		})
	}
}
