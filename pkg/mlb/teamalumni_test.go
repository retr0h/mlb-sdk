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
func TestClient_TeamAlumni(t *testing.T) {
	cases := []struct {
		name       string
		teamID     int
		query      TeamAlumniQuery
		respStatus int
		respBody   string
		wantErr    string
		wantIs     error
		wantLen    int
	}{
		{
			name: "happy", teamID: 119, query: TeamAlumniQuery{Season: 2024, Group: "hitting", Hydrate: "person", Fields: "people"},
			respStatus: 200, respBody: `{"people":[{"id":1,"fullName":"Test"}]}`, wantLen: 1,
		},
		{
			name: "missing required", teamID: 119, query: TeamAlumniQuery{Season: 2024},
			wantErr: "Season and Group are both required", wantIs: ErrInvalidQuery,
		},
		{
			name: "empty", teamID: 119, query: TeamAlumniQuery{Season: 2024, Group: "hitting"},
			respStatus: 200, respBody: `{}`, wantLen: 0,
		},
		{
			name: "404", teamID: 119, query: TeamAlumniQuery{Season: 2024, Group: "hitting"},
			respStatus: 404, respBody: `{}`, wantIs: ErrNotFound, wantErr: "not found",
		},
		{
			name: "5xx", teamID: 119, query: TeamAlumniQuery{Season: 2024, Group: "hitting"},
			respStatus: 500, respBody: `oops`, wantErr: "unexpected status 500",
		},
		{
			name: "bad json", teamID: 119, query: TeamAlumniQuery{Season: 2024, Group: "hitting"},
			respStatus: 200, respBody: `x`, wantErr: "teamAlumni",
		},
		{
			name: "network", teamID: 119, query: TeamAlumniQuery{Season: 2024, Group: "hitting"},
			respStatus: 0, wantErr: "teamAlumni",
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
			pp, err := New(WithBaseURL(u)).TeamAlumni(context.Background(), c.teamID, c.query)
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
