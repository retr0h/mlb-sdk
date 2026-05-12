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

func TestClient_GameColor(t *testing.T) {
	cases := []struct {
		name       string
		gamePk     int
		query      GameColorQuery
		respStatus int
		respBody   string
		wantErr    string
		wantIs     error
	}{
		{
			name: "happy", gamePk: 1, query: GameColorQuery{Timecode: "20240907_180000", Fields: "items"},
			respStatus: 200, respBody: `{"items":[1]}`,
		},
		{
			name:       "404",
			gamePk:     1,
			respStatus: 404,
			respBody:   `{}`,
			wantIs:     ErrNotFound,
			wantErr:    "not found",
		},
		{
			name:       "5xx",
			gamePk:     1,
			respStatus: 500,
			respBody:   `oops`,
			wantErr:    "unexpected status 500",
		},
		{name: "bad json", gamePk: 1, respStatus: 200, respBody: `x`, wantErr: "gameColor"},
		{name: "network", gamePk: 1, respStatus: 0, wantErr: "gameColor"},
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
			_, err := New(WithBaseURL(u)).GameColor(context.Background(), c.gamePk, c.query)
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

func TestClient_GameColorDiff(t *testing.T) {
	cases := []struct {
		name       string
		gamePk     int
		query      GameColorDiffQuery
		respStatus int
		respBody   string
		wantErr    string
		wantIs     error
	}{
		{
			name: "happy", gamePk: 1, query: GameColorDiffQuery{StartTimecode: "a", EndTimecode: "b"},
			respStatus: 200, respBody: `[{"op":"add"}]`,
		},
		{
			name: "missing timecodes", gamePk: 1, query: GameColorDiffQuery{StartTimecode: "a"},
			wantErr: "StartTimecode and EndTimecode", wantIs: ErrInvalidQuery,
		},
		{
			name: "404", gamePk: 1, query: GameColorDiffQuery{StartTimecode: "a", EndTimecode: "b"},
			respStatus: 404, respBody: `{}`, wantIs: ErrNotFound, wantErr: "not found",
		},
		{
			name: "5xx", gamePk: 1, query: GameColorDiffQuery{StartTimecode: "a", EndTimecode: "b"},
			respStatus: 500, respBody: `oops`, wantErr: "gameColorDiff",
		},
		{
			name: "bad json", gamePk: 1, query: GameColorDiffQuery{StartTimecode: "a", EndTimecode: "b"},
			respStatus: 200, respBody: `not json`, wantErr: "gameColorDiff",
		},
		{
			name: "network", gamePk: 1, query: GameColorDiffQuery{StartTimecode: "a", EndTimecode: "b"},
			respStatus: 0, wantErr: "gameColorDiff",
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
			_, err := New(WithBaseURL(u)).GameColorDiff(context.Background(), c.gamePk, c.query)
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

func TestClient_GameDiff(t *testing.T) {
	cases := []struct {
		name       string
		gamePk     int
		query      GameDiffQuery
		respStatus int
		respBody   string
		wantErr    string
		wantIs     error
	}{
		{
			name: "happy", gamePk: 1, query: GameDiffQuery{StartTimecode: "a", EndTimecode: "b"},
			respStatus: 200, respBody: `[{"op":"add"}]`,
		},
		{
			name: "missing timecodes", gamePk: 1, query: GameDiffQuery{StartTimecode: "a"},
			wantErr: "StartTimecode and EndTimecode", wantIs: ErrInvalidQuery,
		},
		{
			name: "404", gamePk: 1, query: GameDiffQuery{StartTimecode: "a", EndTimecode: "b"},
			respStatus: 404, respBody: `{}`, wantIs: ErrNotFound, wantErr: "not found",
		},
		{
			name: "5xx", gamePk: 1, query: GameDiffQuery{StartTimecode: "a", EndTimecode: "b"},
			respStatus: 500, respBody: `oops`, wantErr: "gameDiff",
		},
		{
			name: "bad json", gamePk: 1, query: GameDiffQuery{StartTimecode: "a", EndTimecode: "b"},
			respStatus: 200, respBody: `not json`, wantErr: "gameDiff",
		},
		{
			name: "network", gamePk: 1, query: GameDiffQuery{StartTimecode: "a", EndTimecode: "b"},
			respStatus: 0, wantErr: "gameDiff",
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
			_, err := New(WithBaseURL(u)).GameDiff(context.Background(), c.gamePk, c.query)
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

func TestClient_GameContent(t *testing.T) {
	cases := []struct {
		name       string
		gamePk     int
		query      GameContentQuery
		respStatus int
		respBody   string
		wantErr    string
		wantIs     error
	}{
		{
			name: "happy", gamePk: 1, query: GameContentQuery{HighlightLimit: 5},
			respStatus: 200, respBody: `{"editorial":{}}`,
		},
		{
			name:       "404",
			gamePk:     1,
			respStatus: 404,
			respBody:   `{}`,
			wantIs:     ErrNotFound,
			wantErr:    "not found",
		},
		{
			name:       "5xx",
			gamePk:     1,
			respStatus: 500,
			respBody:   `oops`,
			wantErr:    "unexpected status 500",
		},
		{name: "bad json", gamePk: 1, respStatus: 200, respBody: `x`, wantErr: "gameContent"},
		{name: "network", gamePk: 1, respStatus: 0, wantErr: "gameContent"},
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
			_, err := New(WithBaseURL(u)).GameContent(context.Background(), c.gamePk, c.query)
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

func TestClient_GameWinProbability(t *testing.T) {
	cases := []struct {
		name       string
		gamePk     int
		query      GameWinProbabilityQuery
		respStatus int
		respBody   string
		wantErr    string
		wantIs     error
	}{
		{
			name: "happy", gamePk: 1, query: GameWinProbabilityQuery{Timecode: "20240907_180000", Fields: "items"},
			respStatus: 200, respBody: `[{"homeTeamWinProbability":52.2}]`,
		},
		{
			name:       "404",
			gamePk:     1,
			respStatus: 404,
			respBody:   `{}`,
			wantIs:     ErrNotFound,
			wantErr:    "not found",
		},
		{name: "5xx", gamePk: 1, respStatus: 500, respBody: `oops`, wantErr: "gameWinProbability"},
		{
			name: "bad json", gamePk: 1, query: GameWinProbabilityQuery{},
			respStatus: 200, respBody: `not json`, wantErr: "gameWinProbability",
		},
		{name: "network", gamePk: 1, respStatus: 0, wantErr: "gameWinProbability"},
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
			_, err := New(
				WithBaseURL(u),
			).GameWinProbability(context.Background(), c.gamePk, c.query)
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

func TestClient_Meta(t *testing.T) {
	cases := []struct {
		name       string
		metaType   string
		respStatus int
		respBody   string
		wantErr    string
		wantIs     error
		wantLen    int
	}{
		{
			name: "happy", metaType: "gameTypes", respStatus: 200,
			respBody: `[{"id":"R","description":"Regular Season"}]`, wantLen: 1,
		},
		{
			name: "404", metaType: "bogus", respStatus: 404, respBody: `{}`,
			wantIs: ErrNotFound, wantErr: "not found",
		},
		{
			name: "5xx", metaType: "gameTypes", respStatus: 500, respBody: `oops`,
			wantErr: "meta",
		},
		{name: "network", metaType: "gameTypes", respStatus: 0, wantErr: "meta"},
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
			res, err := New(WithBaseURL(u)).Meta(context.Background(), c.metaType)
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
