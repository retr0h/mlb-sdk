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
func TestClient_TeamLeaders(t *testing.T) {
	body := `{"teamLeaders":[{"leaderCategory":"homeRuns","season":"2024","totalSplits":10,"leaders":[{"rank":1,"value":"54","person":{"id":1,"fullName":"Test"}}]}]}`
	cases := []struct {
		name       string
		teamID     int
		query      TeamLeadersQuery
		respStatus int
		respBody   string
		wantErr    string
		wantIs     error
		wantLen    int
	}{
		{
			name: "happy", teamID: 119,
			query: TeamLeadersQuery{
				LeaderCategories: "homeRuns",
				Season:           2024,
				LeaderGameTypes:  "R",
				Hydrate:          "person",
				Limit:            5,
				Fields:           "teamLeaders",
			},
			respStatus: 200, respBody: body, wantLen: 1,
		},
		{
			name: "missing required", teamID: 119, query: TeamLeadersQuery{Season: 2024},
			wantErr: "LeaderCategories and Season are both required", wantIs: ErrInvalidQuery,
		},
		{
			name: "empty", teamID: 119, query: TeamLeadersQuery{LeaderCategories: "homeRuns", Season: 2024},
			respStatus: 200, respBody: `{}`, wantLen: 0,
		},
		{
			name: "404", teamID: 119, query: TeamLeadersQuery{LeaderCategories: "homeRuns", Season: 2024},
			respStatus: 404, respBody: `{}`, wantIs: ErrNotFound, wantErr: "not found",
		},
		{
			name: "5xx", teamID: 119, query: TeamLeadersQuery{LeaderCategories: "homeRuns", Season: 2024},
			respStatus: 500, respBody: `oops`, wantErr: "unexpected status 500",
		},
		{
			name: "bad json", teamID: 119, query: TeamLeadersQuery{LeaderCategories: "homeRuns", Season: 2024},
			respStatus: 200, respBody: `x`, wantErr: "teamLeaders",
		},
		{
			name: "network", teamID: 119, query: TeamLeadersQuery{LeaderCategories: "homeRuns", Season: 2024},
			respStatus: 0, wantErr: "teamLeaders",
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
			sl, err := New(WithBaseURL(u)).TeamLeaders(context.Background(), c.teamID, c.query)
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
			if len(sl.LeagueLeaders) != c.wantLen {
				t.Errorf("len = %d, want %d", len(sl.LeagueLeaders), c.wantLen)
			}
		})
	}
}
