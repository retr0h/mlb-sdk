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

const liveFeedHappyBody = `{
  "liveData": {
    "plays": {
      "allPlays": [
        {
          "result": {"event": "Strikeout", "eventType": "strikeout", "description": "K swinging."},
          "about": {"inning": 1, "halfInning": "top", "outs": 1}
        },
        {
          "result": {"event": "Grounded Into DP", "eventType": "grounded_into_double_play", "description": "GIDP 6-4-3."},
          "about": {"inning": 4, "halfInning": "bottom", "outs": 2}
        }
      ]
    }
  }
}`

func TestClient_LiveFeed(t *testing.T) {
	cases := []struct {
		name        string
		respStatus  int
		respBody    string
		gamePk      int
		wantLen     int
		wantDPCount int
		wantIs      error
		wantErr     string
	}{
		{
			name:        "happy path: 2 plays incl. one DP",
			respStatus:  200,
			respBody:    liveFeedHappyBody,
			gamePk:      823970,
			wantLen:     2,
			wantDPCount: 1,
		},
		{
			name:       "200 with no liveData yields empty slice",
			respStatus: 200,
			respBody:   `{}`,
			wantLen:    0,
		},
		{
			name:       "200 with liveData but no plays yields empty slice",
			respStatus: 200,
			respBody:   `{"liveData": {}}`,
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
			wantErr:    "liveFeed",
		},
		{
			name:       "network failure is wrapped",
			respStatus: 0,
			wantErr:    "liveFeed",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if !strings.Contains(r.URL.Path, "/api/v1.1/game/") {
					t.Errorf("LiveFeed must hit v1.1, got %q", r.URL.Path)
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(c.respStatus)
				_, _ = w.Write([]byte(c.respBody))
			}))
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
			plays, err := client.LiveFeed(context.Background(), c.gamePk)

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
