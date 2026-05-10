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

func TestBoxscoreTeamLookupAndDoublePlays(t *testing.T) {
	// Compose a minimal Boxscore via the same conversion path the live
	// client uses, so the test exercises both gen→public mapping and the
	// DP parser.
	dodgerInfo := gen.InfoSection{
		Title: strPtr("FIELDING"),
		FieldList: &[]gen.InfoItem{
			{Label: strPtr("DP"), Value: strPtr("2 (Freeland, A; Betts-Freeman).")},
		},
	}
	bravesInfo := gen.InfoSection{
		Title: strPtr("FIELDING"),
		FieldList: &[]gen.InfoItem{
			{Label: strPtr("E"), Value: strPtr("Jarvis (1, throw).")},
		},
	}

	resp := &gen.BoxscoreResponse{
		Teams: &gen.BoxscoreTeams{
			Home: &gen.BoxscoreSide{
				Team: &gen.Team{Id: intPtr(int(LAD)), Name: strPtr("Los Angeles Dodgers")},
				Info: &[]gen.InfoSection{dodgerInfo},
			},
			Away: &gen.BoxscoreSide{
				Team: &gen.Team{Id: intPtr(int(ATL)), Name: strPtr("Atlanta Braves")},
				Info: &[]gen.InfoSection{bravesInfo},
			},
		},
	}
	box := boxscoreFromGen(resp)

	if got := box.Team(LAD).DoublePlaysTurned(); got != 2 {
		t.Errorf("Dodgers DoublePlaysTurned = %d, want 2", got)
	}
	if got := box.Team(ATL).DoublePlaysTurned(); got != 0 {
		t.Errorf("Braves DoublePlaysTurned = %d, want 0 (no DP entry)", got)
	}
	if got := box.Team(NYY); got != nil {
		t.Errorf("Team(NYY) = %v, want nil for team not in boxscore", got)
	}
}

func TestBoxscoreOverHTTP(t *testing.T) {
	body := `{
		"teams": {
			"home": {
				"team": {"id": 119, "name": "Los Angeles Dodgers"},
				"info": [{"title":"FIELDING","fieldList":[{"label":"DP","value":"3 (a; b; c)."}]}]
			},
			"away": {"team": {"id": 144, "name": "Atlanta Braves"}}
		}
	}`
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/api/v1/game/") {
			t.Errorf("unexpected path %q", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(body))
	}))
	defer srv.Close()

	c, err := New(WithBaseURL(srv.URL))
	if err != nil {
		t.Fatal(err)
	}
	box, err := c.Boxscore(context.Background(), 823957)
	if err != nil {
		t.Fatal(err)
	}
	if got := box.Team(LAD).DoublePlaysTurned(); got != 3 {
		t.Errorf("DoublePlaysTurned over HTTP = %d, want 3", got)
	}
}

func strPtr(s string) *string { return &s }
func intPtr(n int) *int       { return &n }
