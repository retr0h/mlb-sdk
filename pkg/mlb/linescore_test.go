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
	"net/url"
	"strings"
	"testing"
//
	"github.com/retr0h/mlb-sdk/internal/gen"
)
//
const linescoreHappyBody = `{
  "copyright": "Copyright 2026 MLB Advanced Media, L.P.",
  "currentInning": 9, "currentInningOrdinal": "9th",
  "inningState": "Bottom", "inningHalf": "Bottom",
  "isTopInning": false, "scheduledInnings": 9,
  "innings": [
    {"num": 1, "ordinalNum": "1st",
     "home": {"runs": 0, "hits": 0, "errors": 0, "leftOnBase": 1},
     "away": {"runs": 0, "hits": 1, "errors": 0, "leftOnBase": 1}},
    {"num": 2, "ordinalNum": "2nd",
     "home": {"runs": 0, "hits": 0, "errors": 0, "leftOnBase": 1},
     "away": {"runs": 1, "hits": 1, "errors": 0, "leftOnBase": 1}}
  ],
  "teams": {
    "home": {"runs": 3, "hits": 5, "errors": 1, "leftOnBase": 5, "isWinner": false},
    "away": {"runs": 5, "hits": 7, "errors": 0, "leftOnBase": 7, "isWinner": true}
  },
  "defense": {
    "pitcher":   {"id": 640448, "fullName": "Kyle Finnegan", "link": "/api/v1/people/640448"},
    "catcher":   {"id": 660688, "fullName": "Keibert Ruiz",  "link": "/api/v1/people/660688"},
    "first":     {"id": 100001, "fullName": "First Base"},
    "second":    {"id": 100002, "fullName": "Second Base"},
    "third":     {"id": 100003, "fullName": "Third Base"},
    "shortstop": {"id": 100004, "fullName": "Shortstop"},
    "left":      {"id": 100005, "fullName": "Left Field"},
    "center":    {"id": 100006, "fullName": "Center Field"},
    "right":     {"id": 100007, "fullName": "Right Field"}
  },
  "offense": {
    "batter": {"id": 650559, "fullName": "Bryan De La Cruz"},
    "onDeck": {"id": 669707, "fullName": "Jared Triolo"},
    "inHole": {"id": 100008, "fullName": "In Hole"},
    "first":  {"id": 100009, "fullName": "Runner 1B"},
    "second": {"id": 100010, "fullName": "Runner 2B"},
    "third":  {"id": 100011, "fullName": "Runner 3B"}
  },
  "balls": 0, "strikes": 1, "outs": 3
}`
//
func TestLinescoreFromGen(t *testing.T) {
	cases := []struct {
		name       string
		in         *gen.LinescoreResponse
		wantInning int
	}{
		{"nil response yields empty Linescore", nil, 0},
		{"empty struct yields empty Linescore", &gen.LinescoreResponse{}, 0},
		{
			"innings populated",
			&gen.LinescoreResponse{Innings: &[]gen.LinescoreInning{{}}},
			0,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := linescoreFromGen(c.in)
			if got == nil {
				t.Fatal("linescoreFromGen returned nil")
			}
			if got.CurrentInning != c.wantInning {
				t.Errorf("CurrentInning = %d, want %d", got.CurrentInning, c.wantInning)
			}
		})
	}
}
//
func TestClient_Linescore(t *testing.T) {
	cases := []struct {
		name         string
		gamePk       int
		query        LinescoreQuery
		respStatus   int
		respBody     string
		wantPath     string
		wantQuery    url.Values
		wantErr      string
		wantIs       error
		wantHydrated bool
	}{
		{
			name:   "happy path: hydrated linescore",
			gamePk: 745455,
			query: LinescoreQuery{
				Timecode: "20240701_180000",
				Fields:   "innings,teams",
			},
			respStatus: 200,
			respBody:   linescoreHappyBody,
			wantPath:   "/api/v1/game/745455/linescore",
			wantQuery: url.Values{
				"timecode": {"20240701_180000"},
				"fields":   {"innings,teams"},
			},
			wantHydrated: true,
		},
		{
			name:       "200 with minimal body parses cleanly",
			gamePk:     1,
			respStatus: 200,
			respBody:   `{"currentInning": 3}`,
		},
		{
			name:       "404 returns ErrNotFound",
			gamePk:     9999,
			respStatus: 404,
			respBody:   `{}`,
			wantIs:     ErrNotFound,
			wantErr:    "not found",
		},
		{
			name:       "5xx is wrapped",
			gamePk:     745455,
			respStatus: 500,
			respBody:   `oops`,
			wantErr:    "unexpected status 500",
		},
		{
			name:       "malformed JSON is wrapped",
			gamePk:     745455,
			respStatus: 200,
			respBody:   `not json`,
			wantErr:    "linescore",
		},
		{
			name:       "network failure is wrapped",
			gamePk:     745455,
			respStatus: 0,
			wantErr:    "linescore",
		},
	}
//
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var seenPath string
			var seenQuery url.Values
			srv := httptest.NewServer(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					seenPath = r.URL.Path
					seenQuery = r.URL.Query()
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
//
			client := New(WithBaseURL(urlStr))
			ls, err := client.Linescore(context.Background(), c.gamePk, c.query)
//
			if c.wantErr != "" {
				if err == nil {
					t.Fatalf("expected error containing %q, got nil", c.wantErr)
				}
				if !strings.Contains(err.Error(), c.wantErr) {
					t.Errorf("err = %v, want substring %q", err, c.wantErr)
				}
				if c.wantIs != nil && !errors.Is(err, c.wantIs) {
					t.Errorf("errors.Is(err, %v) = false, want true", c.wantIs)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if ls == nil {
				t.Fatal("expected non-nil Linescore")
			}
			if c.wantPath != "" && seenPath != c.wantPath {
				t.Errorf("path = %q, want %q", seenPath, c.wantPath)
			}
			if c.wantQuery != nil {
				for k, want := range c.wantQuery {
					if got := seenQuery.Get(k); got != want[0] {
						t.Errorf("query[%q] = %q, want %q", k, got, want[0])
					}
				}
			}
			if !c.wantHydrated {
				return
			}
			if ls.CurrentInning != 9 || ls.CurrentInningOrdinal != "9th" ||
				ls.InningState != "Bottom" || ls.InningHalf != "Bottom" ||
				ls.IsTopInning != false || ls.ScheduledInnings != 9 {
				t.Errorf("top-level = inning=%d ord=%q state=%q half=%q top=%v sched=%d",
					ls.CurrentInning, ls.CurrentInningOrdinal,
					ls.InningState, ls.InningHalf,
					ls.IsTopInning, ls.ScheduledInnings)
			}
			if ls.Balls != 0 || ls.Strikes != 1 || ls.Outs != 3 {
				t.Errorf("count = balls=%d strikes=%d outs=%d",
					ls.Balls, ls.Strikes, ls.Outs)
			}
			if len(ls.Innings) != 2 {
				t.Fatalf("len(Innings) = %d, want 2", len(ls.Innings))
			}
			inn1 := ls.Innings[0]
			if inn1.Num != 1 || inn1.OrdinalNum != "1st" ||
				inn1.Away.Runs != 0 || inn1.Away.Hits != 1 ||
				inn1.Home.LeftOnBase != 1 {
				t.Errorf("Innings[0] = %+v", inn1)
			}
			if ls.Teams.Away.Runs != 5 || ls.Teams.Away.IsWinner != true ||
				ls.Teams.Home.Runs != 3 || ls.Teams.Home.Errors != 1 ||
				ls.Teams.Home.LeftOnBase != 5 || ls.Teams.Away.LeftOnBase != 7 {
				t.Errorf("Teams = %+v", ls.Teams)
			}
			if ls.Defense.Pitcher.ID != 640448 ||
				ls.Defense.Pitcher.FullName != "Kyle Finnegan" ||
				ls.Defense.Catcher.ID != 660688 ||
				ls.Defense.First.ID != 100001 ||
				ls.Defense.Second.ID != 100002 ||
				ls.Defense.Third.ID != 100003 ||
				ls.Defense.Shortstop.ID != 100004 ||
				ls.Defense.Left.ID != 100005 ||
				ls.Defense.Center.ID != 100006 ||
				ls.Defense.Right.ID != 100007 {
				t.Errorf("Defense = %+v", ls.Defense)
			}
			if ls.Offense.Batter.ID != 650559 ||
				ls.Offense.OnDeck.ID != 669707 ||
				ls.Offense.InHole.ID != 100008 ||
				ls.Offense.First.ID != 100009 ||
				ls.Offense.Second.ID != 100010 ||
				ls.Offense.Third.ID != 100011 {
				t.Errorf("Offense = %+v", ls.Offense)
			}
		})
	}
}
