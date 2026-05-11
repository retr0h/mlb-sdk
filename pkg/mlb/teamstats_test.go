// Copyright (c) 2026 John Dewey

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to
// deal in the Software without restriction, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
// sell copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER
// DEALINGS IN THE SOFTWARE.

package mlb

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

const teamStatsHappyBody = `{
  "stats": [{
    "type":  {"displayName": "season"},
    "group": {"displayName": "fielding"},
    "splits": [
      {"season": "2026", "stat": {"doublePlays": 23, "errors": 12, "fielding": ".985", "rangeFactorPerGame": 4.5}},
      {"season": "2025", "stat": {"doublePlays": 19, "errors": 18}}
    ]
  }]
}`

const teamStatsMultiGroupBody = `{
  "stats": [
    {"type": {"displayName":"season"}, "group": {"displayName":"hitting"},  "splits": [{"season":"2026","stat":{"homeRuns":210}}]},
    {"type": {"displayName":"season"}, "group": {"displayName":"fielding"}, "splits": [{"season":"2026","stat":{"doublePlays":23}}]}
  ]
}`

func TestTeamStatsSplit_Int(t *testing.T) {
	split := &TeamStatsSplit{Stat: map[string]any{
		"a": float64(7),
		"b": int64(99),
		"c": int(5),
		"d": "not a number",
	}}
	cases := []struct {
		name  string
		split *TeamStatsSplit
		key   string
		want  int
	}{
		{"nil receiver", nil, "a", 0},
		{"missing key", split, "missing", 0},
		{"float64 (JSON number)", split, "a", 7},
		{"int64", split, "b", 99},
		{"int", split, "c", 5},
		{"non-numeric type", split, "d", 0},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := c.split.Int(c.key); got != c.want {
				t.Errorf("Int(%q) = %d, want %d", c.key, got, c.want)
			}
		})
	}
}

func TestTeamStatsSplit_Float(t *testing.T) {
	split := &TeamStatsSplit{Stat: map[string]any{
		"a": float64(4.5),
		"b": "string not number",
		"c": int(3),
	}}
	cases := []struct {
		name  string
		split *TeamStatsSplit
		key   string
		want  float64
	}{
		{"nil receiver", nil, "a", 0},
		{"missing key", split, "missing", 0},
		{"float64", split, "a", 4.5},
		{"non-float type (string)", split, "b", 0},
		{"non-float type (int — Float only accepts float64)", split, "c", 0},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := c.split.Float(c.key); got != c.want {
				t.Errorf("Float(%q) = %v, want %v", c.key, got, c.want)
			}
		})
	}
}

func TestTeamStatsSplit_String(t *testing.T) {
	split := &TeamStatsSplit{Stat: map[string]any{
		"a": "good",
		"b": float64(7),
	}}
	cases := []struct {
		name  string
		split *TeamStatsSplit
		key   string
		want  string
	}{
		{"nil receiver", nil, "a", ""},
		{"missing key", split, "missing", ""},
		{"string", split, "a", "good"},
		{"non-string type", split, "b", ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := c.split.String(c.key); got != c.want {
				t.Errorf("String(%q) = %q, want %q", c.key, got, c.want)
			}
		})
	}
}

func TestTeamStatsSplit_DoublePlays(t *testing.T) {
	cases := []struct {
		name  string
		split *TeamStatsSplit
		want  int
	}{
		{"nil receiver", nil, 0},
		{"missing field", &TeamStatsSplit{Stat: map[string]any{"errors": 3.0}}, 0},
		{"present field", &TeamStatsSplit{Stat: map[string]any{"doublePlays": float64(23)}}, 23},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := c.split.DoublePlays(); got != c.want {
				t.Errorf("DoublePlays() = %d, want %d", got, c.want)
			}
		})
	}
}

func TestTeamStatGroupResult_Season(t *testing.T) {
	g := &TeamStatGroupResult{Splits: []TeamStatsSplit{
		{Season: "2025"}, {Season: "2026"},
	}}
	cases := []struct {
		name   string
		group  *TeamStatGroupResult
		season string
		wantOk bool
	}{
		{"nil receiver", nil, "2026", false},
		{"miss", g, "1999", false},
		{"hit", g, "2026", true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := c.group.Season(c.season)
			if (got != nil) != c.wantOk {
				t.Errorf("Season(%q) ok=%v, want %v", c.season, got != nil, c.wantOk)
			}
			if got != nil && got.Season != c.season {
				t.Errorf("Season(%q).Season = %q, want %q", c.season, got.Season, c.season)
			}
		})
	}
}

func TestTeamStats_Group(t *testing.T) {
	ts := &TeamStats{Groups: []TeamStatGroupResult{
		{Group: "hitting"},
		{Group: "fielding"},
	}}
	cases := []struct {
		name   string
		ts     *TeamStats
		group  TeamStatGroup
		wantOk bool
	}{
		{"nil receiver", nil, TeamStatGroupFielding, false},
		{"miss", ts, TeamStatGroupPitching, false},
		{"hit (exact case)", ts, TeamStatGroupFielding, true},
		{
			"hit (case-insensitive)",
			&TeamStats{Groups: []TeamStatGroupResult{{Group: "Fielding"}}},
			TeamStatGroupFielding,
			true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := c.ts.Group(c.group)
			if (got != nil) != c.wantOk {
				t.Errorf("Group(%v) ok=%v, want %v", c.group, got != nil, c.wantOk)
			}
		})
	}
}

func TestClient_TeamStats(t *testing.T) {
	cases := []struct {
		name          string
		query         TeamStatsQuery
		respStatus    int
		respBody      string
		wantErr       string
		wantIs        error
		wantPath      string
		wantQuery     url.Values
		wantNumGroups int
		wantDPs       int // checked when wantNumGroups > 0; 0 means don't check
		wantHomeRuns  int // checked when wantNumGroups > 0
	}{
		{
			name: "happy path: fielding stats with two season splits",
			query: TeamStatsQuery{
				Team: LAD, Season: 2026,
				Type: TeamStatTypeSeason, Group: TeamStatGroupFielding,
			},
			respStatus: 200,
			respBody:   teamStatsHappyBody,
			wantPath:   "/api/v1/teams/119/stats",
			wantQuery: url.Values{
				"season": {"2026"},
				"stats":  {"season"},
				"group":  {"fielding"},
			},
			wantNumGroups: 1,
			wantDPs:       23,
		},
		{
			name:          "happy path: multi-group response is preserved",
			query:         TeamStatsQuery{Team: LAD},
			respStatus:    200,
			respBody:      teamStatsMultiGroupBody,
			wantPath:      "/api/v1/teams/119/stats",
			wantNumGroups: 2,
			wantDPs:       23,
			wantHomeRuns:  210,
		},
		{
			name:       "200 with no stats yields empty groups",
			query:      TeamStatsQuery{Team: LAD},
			respStatus: 200,
			respBody:   `{}`,
		},
		{
			name:    "missing required Team is rejected before HTTP call",
			query:   TeamStatsQuery{Season: 2026},
			wantErr: "Team is required",
			wantIs:  ErrInvalidQuery,
		},
		{
			name:       "404 returns ErrNotFound",
			query:      TeamStatsQuery{Team: LAD},
			respStatus: 404,
			respBody:   `{}`,
			wantIs:     ErrNotFound,
			wantErr:    "not found",
		},
		{
			name:       "5xx is wrapped",
			query:      TeamStatsQuery{Team: LAD},
			respStatus: 500,
			respBody:   `oops`,
			wantErr:    "unexpected status 500",
		},
		{
			name:       "malformed JSON is wrapped",
			query:      TeamStatsQuery{Team: LAD},
			respStatus: 200,
			respBody:   `not json`,
			wantErr:    "teamStats",
		},
		{
			name:       "network failure is wrapped",
			query:      TeamStatsQuery{Team: LAD},
			respStatus: 0,
			wantErr:    "teamStats",
		},
	}

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

			client := New(WithBaseURL(urlStr))
			ts, err := client.TeamStats(context.Background(), c.query)

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
			if ts == nil {
				t.Fatal("expected non-nil TeamStats")
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
			if got := len(ts.Groups); got != c.wantNumGroups {
				t.Errorf("len(Groups) = %d, want %d", got, c.wantNumGroups)
			}
			if c.wantDPs > 0 {
				got := ts.Group(TeamStatGroupFielding).Season("2026").DoublePlays()
				if got != c.wantDPs {
					t.Errorf("fielding 2026 DoublePlays() = %d, want %d", got, c.wantDPs)
				}
			}
			if c.wantHomeRuns > 0 {
				got := ts.Group(TeamStatGroupHitting).Season("2026").Int("homeRuns")
				if got != c.wantHomeRuns {
					t.Errorf("hitting 2026 homeRuns = %d, want %d", got, c.wantHomeRuns)
				}
			}
		})
	}
}
