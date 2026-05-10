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
	"strings"
	"testing"
)

const playByPlayHappyBody = `{
  "allPlays": [
    {
      "result": {"event": "Strikeout", "eventType": "strikeout", "description": "Smith strikes out swinging."},
      "about": {"inning": 1, "halfInning": "top", "outs": 1}
    },
    {
      "result": {"event": "Grounded Into DP", "eventType": "grounded_into_double_play", "description": "Jones grounds into a double play, 6-4-3."},
      "about": {"inning": 3, "halfInning": "bottom", "outs": 2}
    },
    {
      "result": {"event": "Single", "eventType": "single", "description": "Brown singles on a line drive."},
      "about": {"inning": 5, "halfInning": "top", "outs": 0}
    }
  ]
}`

func TestPlay_IsDoublePlay(t *testing.T) {
	cases := []struct {
		name string
		play Play
		want bool
	}{
		{"zero value", Play{}, false},
		{"strikeout", Play{EventType: EventStrikeout}, false},
		{"single", Play{EventType: EventSingle}, false},
		{
			"force out (records 2 outs but not officially a DP)",
			Play{EventType: EventForceOut},
			false,
		},
		{"grounded into double play", Play{EventType: EventGroundedIntoDoublePlay}, true},
		{"unknown event type string", Play{EventType: EventType("triple_play")}, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := c.play.IsDoublePlay(); got != c.want {
				t.Errorf("IsDoublePlay() = %v, want %v", got, c.want)
			}
		})
	}
}

func TestClient_PlayByPlay(t *testing.T) {
	cases := []struct {
		name        string
		respStatus  int
		respBody    string
		gamePk      int
		wantLen     int
		wantDPCount int // count of plays where IsDoublePlay() == true
		wantIs      error
		wantErr     string
	}{
		{
			name:        "happy path returns 3 plays incl. one DP",
			respStatus:  200,
			respBody:    playByPlayHappyBody,
			gamePk:      823970,
			wantLen:     3,
			wantDPCount: 1,
		},
		{
			name:       "200 with no plays yields empty slice",
			respStatus: 200,
			respBody:   `{}`,
			wantLen:    0,
		},
		{
			name:       "404 returns ErrNotFound",
			respStatus: 404,
			respBody:   `{}`,
			wantIs:     ErrNotFound,
			wantErr:    "not found",
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
			wantErr:    "playByPlay",
		},
		{
			name:       "network failure is wrapped",
			respStatus: 0,
			wantErr:    "playByPlay",
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
			urlStr := srv.URL
			if c.respStatus == 0 {
				srv.Close()
			} else {
				defer srv.Close()
			}

			client, err := New(WithBaseURL(urlStr))
			if err != nil {
				t.Fatalf("New: %v", err)
			}
			plays, err := client.PlayByPlay(context.Background(), c.gamePk)

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
			if got := len(plays); got != c.wantLen {
				t.Fatalf("len(plays) = %d, want %d", got, c.wantLen)
			}
			gotDPs := 0
			for _, p := range plays {
				if p.IsDoublePlay() {
					gotDPs++
				}
			}
			if gotDPs != c.wantDPCount {
				t.Errorf("DP count = %d, want %d", gotDPs, c.wantDPCount)
			}
		})
	}
}
