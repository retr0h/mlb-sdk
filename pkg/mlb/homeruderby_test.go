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

func TestClient_HomeRunDerby(t *testing.T) {
	cases := []struct {
		name       string
		gamePk     int
		query      HomeRunDerbyQuery
		respStatus int
		respBody   string
		wantErr    string
		wantIs     error
	}{
		{
			name: "happy", gamePk: 1, query: HomeRunDerbyQuery{Fields: "info"},
			respStatus: 200, respBody: `{"info":{}}`,
		},
		{
			name: "404", gamePk: 9999, respStatus: 404, respBody: `{}`,
			wantIs: ErrNotFound, wantErr: "not found",
		},
		{
			name: "5xx", gamePk: 1, respStatus: 500, respBody: `oops`,
			wantErr: "unexpected status 500",
		},
		{
			name: "bad json", gamePk: 1, respStatus: 200, respBody: `x`,
			wantErr: "homeRunDerby",
		},
		{name: "network", gamePk: 1, respStatus: 0, wantErr: "homeRunDerby"},
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
			_, err := New(WithBaseURL(u)).HomeRunDerby(context.Background(), c.gamePk, c.query)
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
