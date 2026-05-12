// Copyright (c) 2026 John Dewey
//
// SPDX-License-Identifier: MIT

package mlb

import "testing"

func TestLeagueID_Int(t *testing.T) {
	cases := []struct {
		name string
		l    LeagueID
		want int
	}{
		{"zero value", LeagueID(0), 0},
		{"AL", AL, 103},
		{"NL", NL, 104},
		{"unknown numeric preserved", LeagueID(999), 999},
		{"negative value preserved", LeagueID(-1), -1},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := c.l.Int(); got != c.want {
				t.Errorf("LeagueID(%d).Int() = %d, want %d", c.l, got, c.want)
			}
		})
	}
}

func TestLeagueID_String(t *testing.T) {
	cases := []struct {
		name string
		l    LeagueID
		want string
	}{
		{"AL", AL, "103"},
		{"NL", NL, "104"},
		{"zero", LeagueID(0), "0"},
		{"unknown numeric", LeagueID(999), "999"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := c.l.String(); got != c.want {
				t.Errorf("LeagueID(%d).String() = %q, want %q", c.l, got, c.want)
			}
		})
	}
}
