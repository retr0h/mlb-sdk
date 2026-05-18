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
func TestClient_GameTimestamps(t *testing.T) {
	cases := []struct {
		name       string
		gamePk     int
		respStatus int
		respBody   string
		wantErr    string
		wantIs     error
		wantLen    int
	}{
		{
			name: "happy path", gamePk: 745455, respStatus: 200,
			respBody: `{"timestamps":["20240907_173500","20240907_173600"]}`,
			wantLen:  2,
		},
		{name: "empty", gamePk: 1, respStatus: 200, respBody: `{}`, wantLen: 0},
		{
			name: "404", gamePk: 9999, respStatus: 404, respBody: `{}`,
			wantIs: ErrNotFound, wantErr: "not found",
		},
		{
			name: "5xx", gamePk: 1, respStatus: 500, respBody: `oops`,
			wantErr: "unexpected status 500",
		},
		{
			name: "bad json", gamePk: 1, respStatus: 200, respBody: `not json`,
			wantErr: "gameTimestamps",
		},
		{name: "network", gamePk: 1, respStatus: 0, wantErr: "gameTimestamps"},
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
			urlStr := srv.URL
			if c.respStatus == 0 {
				srv.Close()
			} else {
				defer srv.Close()
			}
			client := New(WithBaseURL(urlStr))
			ts, err := client.GameTimestamps(context.Background(), c.gamePk)
			if c.wantErr != "" {
				if err == nil {
					t.Fatalf("expected error containing %q", c.wantErr)
				}
				if !strings.Contains(err.Error(), c.wantErr) {
					t.Errorf("err = %v, want substring %q", err, c.wantErr)
				}
				if c.wantIs != nil && !errors.Is(err, c.wantIs) {
					t.Errorf("errors.Is = false")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(ts) != c.wantLen {
				t.Errorf("len = %d, want %d", len(ts), c.wantLen)
			}
		})
	}
}
//
func TestClient_GameColorTimestamps(t *testing.T) {
	cases := []struct {
		name       string
		gamePk     int
		respStatus int
		respBody   string
		wantErr    string
		wantIs     error
		wantLen    int
	}{
		{
			name: "happy path", gamePk: 745455, respStatus: 200,
			respBody: `{"timestamps":["20240907_173500"]}`,
			wantLen:  1,
		},
		{name: "empty", gamePk: 1, respStatus: 200, respBody: `{}`, wantLen: 0},
		{
			name: "404", gamePk: 9999, respStatus: 404, respBody: `{}`,
			wantIs: ErrNotFound, wantErr: "not found",
		},
		{
			name: "5xx", gamePk: 1, respStatus: 500, respBody: `oops`,
			wantErr: "unexpected status 500",
		},
		{
			name: "bad json", gamePk: 1, respStatus: 200, respBody: `not json`,
			wantErr: "gameColorTimestamps",
		},
		{name: "network", gamePk: 1, respStatus: 0, wantErr: "gameColorTimestamps"},
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
			urlStr := srv.URL
			if c.respStatus == 0 {
				srv.Close()
			} else {
				defer srv.Close()
			}
			client := New(WithBaseURL(urlStr))
			ts, err := client.GameColorTimestamps(context.Background(), c.gamePk)
			if c.wantErr != "" {
				if err == nil {
					t.Fatalf("expected error containing %q", c.wantErr)
				}
				if !strings.Contains(err.Error(), c.wantErr) {
					t.Errorf("err = %v, want substring %q", err, c.wantErr)
				}
				if c.wantIs != nil && !errors.Is(err, c.wantIs) {
					t.Errorf("errors.Is = false")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(ts) != c.wantLen {
				t.Errorf("len = %d, want %d", len(ts), c.wantLen)
			}
		})
	}
}
