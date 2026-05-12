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

func TestClient_StatsStreaks(t *testing.T) {
	q := StatsStreaksQuery{
		StreakType: "hittingStreak", StreakSpan: "season",
		Season: 2024, SportID: 1, Limit: 3,
	}
	cases := []struct {
		name       string
		query      StatsStreaksQuery
		respStatus int
		respBody   string
		wantErr    string
		wantIs     error
	}{
		{
			name: "happy",
			query: StatsStreaksQuery{
				StreakType: q.StreakType, StreakSpan: q.StreakSpan,
				Season: q.Season, SportID: q.SportID, Limit: q.Limit,
				GameType: "R", Hydrate: "person", Fields: "streaks",
			},
			respStatus: 200, respBody: `{"streaks":[]}`,
		},
		{
			name: "missing required", query: StatsStreaksQuery{StreakType: "x"},
			wantErr: "StreakType, StreakSpan, Season, SportID, and Limit are all required",
			wantIs:  ErrInvalidQuery,
		},
		{
			name: "404", query: q, respStatus: 404, respBody: `{}`,
			wantIs: ErrNotFound, wantErr: "not found",
		},
		{
			name: "5xx", query: q, respStatus: 500, respBody: `oops`,
			wantErr: "unexpected status 500",
		},
		{
			name: "bad json", query: q, respStatus: 200, respBody: `x`,
			wantErr: "statsStreaks",
		},
		{name: "network", query: q, respStatus: 0, wantErr: "statsStreaks"},
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
			_, err := New(WithBaseURL(u)).StatsStreaks(context.Background(), c.query)
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
