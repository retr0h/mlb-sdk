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

	"github.com/retr0h/mlb-sdk/internal/gen"
)

func TestParseDPCount(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  int
	}{
		{"empty", "", 0},
		{"whitespace only", "   ", 0},
		{"single DP, no leading number", "(Freeman, F-Rojas, M).", 1},
		{"two DPs with leading 2", "2 (Freeland, A-Betts-Freeman, F; Betts-Freeman, F).", 2},
		{"three DPs with leading 3", "3 (2 Rocchio-Arias, G-Hoskins; Hoskins-Arias, G-Hoskins).", 3},
		{"leading whitespace then number", "  2 (X-Y).", 2},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := parseDPCount(c.input); got != c.want {
				t.Errorf("parseDPCount(%q) = %d, want %d", c.input, got, c.want)
			}
		})
	}
}

func TestBoxscoreTeam_DoublePlaysTurned(t *testing.T) {
	mkTeam := func(sections ...gen.InfoSection) *BoxscoreTeam {
		side := &gen.BoxscoreSide{Info: &sections}
		return teamFromGen(side)
	}

	cases := []struct {
		name string
		team *BoxscoreTeam
		want int
	}{
		{"nil receiver", nil, 0},
		{"team with nil raw", &BoxscoreTeam{}, 0},
		{
			"no info sections",
			mkTeam(),
			0,
		},
		{
			"fielding section without DP entry",
			mkTeam(gen.InfoSection{
				Title:     strPtr("FIELDING"),
				FieldList: &[]gen.InfoItem{{Label: strPtr("E"), Value: strPtr("Jarvis (1, throw).")}},
			}),
			0,
		},
		{
			"DP in non-fielding section is ignored",
			mkTeam(gen.InfoSection{
				Title:     strPtr("BATTING"),
				FieldList: &[]gen.InfoItem{{Label: strPtr("DP"), Value: strPtr("2 (X-Y).")}},
			}),
			0,
		},
		{
			"single DP entry",
			mkTeam(gen.InfoSection{
				Title:     strPtr("FIELDING"),
				FieldList: &[]gen.InfoItem{{Label: strPtr("DP"), Value: strPtr("(Freeman, F-Rojas, M).")}},
			}),
			1,
		},
		{
			"two DPs",
			mkTeam(gen.InfoSection{
				Title:     strPtr("FIELDING"),
				FieldList: &[]gen.InfoItem{{Label: strPtr("DP"), Value: strPtr("2 (Freeland; Betts-Freeman).")}},
			}),
			2,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := c.team.DoublePlaysTurned(); got != c.want {
				t.Errorf("DoublePlaysTurned() = %d, want %d", got, c.want)
			}
		})
	}
}

func TestBoxscore_Team(t *testing.T) {
	resp := &gen.BoxscoreResponse{
		Teams: &gen.BoxscoreTeams{
			Home: &gen.BoxscoreSide{Team: &gen.Team{Id: intPtr(int(LAD)), Name: strPtr("Los Angeles Dodgers")}},
			Away: &gen.BoxscoreSide{Team: &gen.Team{Id: intPtr(int(ATL)), Name: strPtr("Atlanta Braves")}},
		},
	}
	box := boxscoreFromGen(resp)

	cases := []struct {
		name     string
		box      *Boxscore
		lookup   TeamID
		wantNil  bool
		wantID   TeamID
		wantName string
	}{
		{"nil receiver", nil, LAD, true, 0, ""},
		{"home team match", box, LAD, false, LAD, "Los Angeles Dodgers"},
		{"away team match", box, ATL, false, ATL, "Atlanta Braves"},
		{"team not in boxscore", box, NYY, true, 0, ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := c.box.Team(c.lookup)
			if (got == nil) != c.wantNil {
				t.Fatalf("Team(%v) nil=%v, wantNil=%v", c.lookup, got == nil, c.wantNil)
			}
			if got == nil {
				return
			}
			if got.ID != c.wantID || got.Name != c.wantName {
				t.Errorf("Team(%v) = {%v %q}, want {%v %q}", c.lookup, got.ID, got.Name, c.wantID, c.wantName)
			}
		})
	}
}

func TestClient_Boxscore(t *testing.T) {
	const happyBody = `{
		"teams": {
			"home": {
				"team": {"id": 119, "name": "Los Angeles Dodgers"},
				"info": [{"title":"FIELDING","fieldList":[{"label":"DP","value":"3 (a; b; c)."}]}]
			},
			"away": {"team": {"id": 144, "name": "Atlanta Braves"}}
		}
	}`

	cases := []struct {
		name       string
		respStatus int    // 0 means close server before request to simulate net failure
		respBody   string
		gamePk     int
		wantNil    bool
		wantDPs    int   // for happy rows; checked when wantErr == ""
		wantIs     error // errors.Is target; nil means no errors.Is check
		wantErr    string
	}{
		{
			name:       "happy path returns parsed boxscore",
			respStatus: 200,
			respBody:   happyBody,
			gamePk:     823957,
			wantDPs:    3,
		},
		{
			name:       "200 with empty body yields zero-value boxscore",
			respStatus: 200,
			respBody:   `{}`,
			gamePk:     1,
		},
		{
			name:       "404 returns ErrNotFound",
			respStatus: 404,
			respBody:   `{}`,
			gamePk:     1,
			wantNil:    true,
			wantIs:     ErrNotFound,
			wantErr:    "not found",
		},
		{
			name:       "5xx is wrapped",
			respStatus: 500,
			respBody:   `oops`,
			gamePk:     1,
			wantNil:    true,
			wantErr:    "unexpected status 500",
		},
		{
			name:       "malformed JSON is wrapped",
			respStatus: 200,
			respBody:   `not json`,
			gamePk:     1,
			wantNil:    true,
			wantErr:    "boxscore",
		},
		{
			name:       "network failure is wrapped",
			respStatus: 0,
			gamePk:     1,
			wantNil:    true,
			wantErr:    "boxscore",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(c.respStatus)
				_, _ = w.Write([]byte(c.respBody))
			}))
			url := srv.URL
			if c.respStatus == 0 {
				srv.Close()
			} else {
				defer srv.Close()
			}

			client, err := New(WithBaseURL(url))
			if err != nil {
				t.Fatalf("New: %v", err)
			}
			box, err := client.Boxscore(context.Background(), c.gamePk)

			if c.wantErr == "" {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if box == nil {
					t.Fatal("expected non-nil boxscore")
				}
				if c.wantDPs > 0 {
					if got := box.Team(LAD).DoublePlaysTurned(); got != c.wantDPs {
						t.Errorf("DoublePlaysTurned() = %d, want %d", got, c.wantDPs)
					}
				}
				return
			}
			if err == nil {
				t.Fatalf("expected error containing %q, got nil", c.wantErr)
			}
			if !strings.Contains(err.Error(), c.wantErr) {
				t.Errorf("err = %v, want substring %q", err, c.wantErr)
			}
			if c.wantIs != nil && !errors.Is(err, c.wantIs) {
				t.Errorf("errors.Is(err, %v) = false, want true", c.wantIs)
			}
			if c.wantNil && box != nil {
				t.Errorf("box = %v, want nil on error", box)
			}
		})
	}
}

func strPtr(s string) *string { return &s }
func intPtr(n int) *int       { return &n }
