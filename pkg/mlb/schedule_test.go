// Copyright (c) 2026 John Dewey
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to
// deal in the Software without restriction, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
// sell copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
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
	"time"
)

const scheduleHappyBody = `{
  "totalGames": 2,
  "dates": [{
    "date": "2026-05-08",
    "games": [
      {
        "gamePk": 823957,
        "gameDate": "2026-05-08T19:10:00Z",
        "status": {"abstractGameState": "Final", "detailedState": "Final"},
        "teams": {
          "home": {"team": {"id": 119, "name": "Los Angeles Dodgers"}, "score": 3},
          "away": {"team": {"id": 144, "name": "Atlanta Braves"},      "score": 1}
        }
      },
      {
        "gamePk": 823958,
        "gameDate": "2026-05-08T22:00:00Z",
        "status": {"abstractGameState": "Live"},
        "teams": {
          "home": {"team": {"id": 137, "name": "San Francisco Giants"}, "score": 0},
          "away": {"team": {"id": 119, "name": "Los Angeles Dodgers"},  "score": 0}
        }
      }
    ]
  }]
}`

func TestClient_Schedule(t *testing.T) {
	parsed := time.Date(2026, 5, 8, 19, 10, 0, 0, time.UTC)

	cases := []struct {
		name        string
		query       ScheduleQuery
		respStatus  int // 0 to simulate network failure
		respBody    string
		wantQuery   url.Values // expected query string the server saw (subset; nil to skip)
		wantLen     int
		wantFirstPk int
		wantFirst   *Game // optional deep-check
		wantIs      error
		wantErr     string
	}{
		{
			name:        "happy path with no filters",
			respStatus:  200,
			respBody:    scheduleHappyBody,
			wantLen:     2,
			wantFirstPk: 823957,
			wantFirst: &Game{
				GamePk:         823957,
				Date:           parsed,
				Status:         StatusFinal,
				DetailedStatus: "Final",
				Home:           TeamScore{ID: LAD, Name: "Los Angeles Dodgers", Score: 3},
				Away:           TeamScore{ID: ATL, Name: "Atlanta Braves", Score: 1},
			},
		},
		{
			name:       "team filter is forwarded as teamId query param",
			query:      ScheduleQuery{Team: LAD},
			respStatus: 200,
			respBody:   `{"dates":[]}`,
			wantQuery:  url.Values{"sportId": {"1"}, "teamId": {"119"}},
		},
		{
			name:       "On filter is forwarded as date query param",
			query:      ScheduleQuery{On: time.Date(2026, 5, 8, 0, 0, 0, 0, time.UTC)},
			respStatus: 200,
			respBody:   `{"dates":[]}`,
			wantQuery:  url.Values{"sportId": {"1"}, "date": {"2026-05-08"}},
		},
		{
			name: "From/To pair is forwarded as startDate+endDate",
			query: ScheduleQuery{
				From: time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC),
				To:   time.Date(2026, 5, 7, 0, 0, 0, 0, time.UTC),
			},
			respStatus: 200,
			respBody:   `{"dates":[]}`,
			wantQuery: url.Values{
				"sportId":   {"1"},
				"startDate": {"2026-05-01"},
				"endDate":   {"2026-05-07"},
			},
		},
		{
			name:       "200 with no dates yields empty slice",
			respStatus: 200,
			respBody:   `{}`,
			wantLen:    0,
		},
		{
			name:       "200 with date that has nil games is skipped",
			respStatus: 200,
			respBody:   `{"dates":[{"date":"2026-05-08"}]}`, // games key omitted → nil
			wantLen:    0,
		},
		{
			name:       "game with nil away side still parses",
			respStatus: 200,
			respBody: `{"dates":[{"games":[{
				"gamePk": 1, "status": {"abstractGameState": "Final"},
				"teams": {"home": {"team": {"id": 119, "name": "Los Angeles Dodgers"}, "score": 3}}
			}]}]}`, // away omitted → teamScoreFromGen(nil) returns zero TeamScore
			wantLen:     1,
			wantFirstPk: 1,
		},
		{
			name:       "404 is wrapped (Schedule does not map to ErrNotFound)",
			respStatus: 404,
			respBody:   `{}`,
			wantErr:    "unexpected status 404",
		},
		{
			name:       "5xx is wrapped",
			respStatus: 500,
			respBody:   `oops`,
			wantErr:    "unexpected status 500",
		},
		{
			name:       "malformed JSON is wrapped",
			respStatus: 200,
			respBody:   `not json`,
			wantErr:    "schedule",
		},
		{
			name:       "network failure is wrapped",
			respStatus: 0,
			wantErr:    "schedule",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var seenQuery url.Values
			srv := httptest.NewServer(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
			games, err := client.Schedule(context.Background(), c.query)

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
			if got := len(games); got != c.wantLen {
				t.Errorf("len(games) = %d, want %d", got, c.wantLen)
			}
			if c.wantFirstPk > 0 && (len(games) == 0 || games[0].GamePk != c.wantFirstPk) {
				t.Errorf("games[0].GamePk = %v, want %d", firstPk(games), c.wantFirstPk)
			}
			if c.wantFirst != nil && len(games) > 0 {
				if games[0] != *c.wantFirst {
					t.Errorf("games[0] = %+v, want %+v", games[0], *c.wantFirst)
				}
			}
			if c.wantQuery != nil {
				for k, want := range c.wantQuery {
					if got := seenQuery.Get(k); got != want[0] {
						t.Errorf("query[%q] = %q, want %q", k, got, want[0])
					}
				}
			}
		})
	}
}

func firstPk(games []Game) any {
	if len(games) == 0 {
		return "<empty>"
	}
	return games[0].GamePk
}
