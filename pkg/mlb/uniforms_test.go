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
func TestClient_GameUniforms(t *testing.T) {
	cases := []struct {
		name       string
		query      GameUniformsQuery
		respStatus int
		respBody   string
		wantErr    string
		wantIs     error
		wantLen    int
	}{
		{
			name: "happy", query: GameUniformsQuery{GamePks: "745455", Fields: "uniforms"},
			respStatus: 200, respBody: `{"uniforms":[{"gamePk":745455}]}`, wantLen: 1,
		},
		{
			name: "missing GamePks", query: GameUniformsQuery{},
			wantErr: "GamePks is required", wantIs: ErrInvalidQuery,
		},
		{
			name: "empty", query: GameUniformsQuery{GamePks: "1"},
			respStatus: 200, respBody: `{}`, wantLen: 0,
		},
		{
			name: "404", query: GameUniformsQuery{GamePks: "1"},
			respStatus: 404, respBody: `{}`, wantIs: ErrNotFound, wantErr: "not found",
		},
		{
			name: "5xx", query: GameUniformsQuery{GamePks: "1"},
			respStatus: 500, respBody: `oops`, wantErr: "unexpected status 500",
		},
		{
			name: "bad json", query: GameUniformsQuery{GamePks: "1"},
			respStatus: 200, respBody: `x`, wantErr: "gameUniforms",
		},
		{
			name: "network", query: GameUniformsQuery{GamePks: "1"},
			respStatus: 0, wantErr: "gameUniforms",
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
			res, err := New(WithBaseURL(u)).GameUniforms(context.Background(), c.query)
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
			if len(res) != c.wantLen {
				t.Errorf("len = %d, want %d", len(res), c.wantLen)
			}
		})
	}
}
//
func TestClient_TeamUniforms(t *testing.T) {
	cases := []struct {
		name       string
		query      TeamUniformsQuery
		respStatus int
		respBody   string
		wantErr    string
		wantIs     error
		wantLen    int
	}{
		{
			name: "happy", query: TeamUniformsQuery{TeamIDs: "119", Season: 2024, Fields: "uniforms"},
			respStatus: 200, respBody: `{"uniforms":[{"teamId":119}]}`, wantLen: 1,
		},
		{
			name: "missing TeamIDs", query: TeamUniformsQuery{},
			wantErr: "TeamIDs is required", wantIs: ErrInvalidQuery,
		},
		{
			name: "empty", query: TeamUniformsQuery{TeamIDs: "119"},
			respStatus: 200, respBody: `{}`, wantLen: 0,
		},
		{
			name: "404", query: TeamUniformsQuery{TeamIDs: "119"},
			respStatus: 404, respBody: `{}`, wantIs: ErrNotFound, wantErr: "not found",
		},
		{
			name: "5xx", query: TeamUniformsQuery{TeamIDs: "119"},
			respStatus: 500, respBody: `oops`, wantErr: "unexpected status 500",
		},
		{
			name: "bad json", query: TeamUniformsQuery{TeamIDs: "119"},
			respStatus: 200, respBody: `x`, wantErr: "teamUniforms",
		},
		{
			name: "network", query: TeamUniformsQuery{TeamIDs: "119"},
			respStatus: 0, wantErr: "teamUniforms",
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
			res, err := New(WithBaseURL(u)).TeamUniforms(context.Background(), c.query)
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
			if len(res) != c.wantLen {
				t.Errorf("len = %d, want %d", len(res), c.wantLen)
			}
		})
	}
}
