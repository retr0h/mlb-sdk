// Copyright (c) 2026 John Dewey
//
// SPDX-License-Identifier: MIT

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

const seasonsHappyBody = `{
  "copyright": "Copyright 2026 MLB Advanced Media, L.P.",
  "seasons": [{
    "seasonId": "2024",
    "hasWildcard": true,
    "preSeasonStartDate": "2024-01-01",
    "preSeasonEndDate": "2024-02-21",
    "seasonStartDate": "2024-02-22",
    "springStartDate": "2024-02-22",
    "springEndDate": "2024-03-26",
    "regularSeasonStartDate": "2024-03-20",
    "lastDate1stHalf": "2024-07-14",
    "allStarDate": "2024-07-16",
    "firstDate2ndHalf": "2024-07-19",
    "regularSeasonEndDate": "2024-09-30",
    "postSeasonStartDate": "2024-10-01",
    "postSeasonEndDate": "2024-10-30",
    "seasonEndDate": "2024-10-30",
    "offseasonStartDate": "2024-10-31",
    "offSeasonEndDate": "2024-12-31",
    "seasonLevelGamedayType": "P",
    "gameLevelGamedayType": "P",
    "qualifierPlateAppearances": 3.1,
    "qualifierOutsPitched": 3.0
  }]
}`

const seasonsAllHappyBody = `{
  "seasons": [
    {"seasonId": "1876", "hasWildcard": false, "regularSeasonStartDate": "1876-04-22"},
    {"seasonId": "2024", "hasWildcard": true,  "regularSeasonStartDate": "2024-03-20"}
  ]
}`

func TestSeasons_Season(t *testing.T) {
	empty := seasonsFromGen(nil)
	full := &Seasons{Seasons: []Season{
		{SeasonID: "2023"},
		{SeasonID: "2024"},
	}}
	cases := []struct {
		name   string
		s      *Seasons
		lookup string
		wantOk bool
	}{
		{"nil receiver", nil, "2024", false},
		{"empty (no seasons)", empty, "2024", false},
		{"miss", full, "1999", false},
		{"hit", full, "2024", true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := c.s.Season(c.lookup)
			if (got != nil) != c.wantOk {
				t.Fatalf("Season(%q) ok=%v, want %v", c.lookup, got != nil, c.wantOk)
			}
			if got != nil && got.SeasonID != c.lookup {
				t.Errorf("Season(%q).SeasonID = %q", c.lookup, got.SeasonID)
			}
		})
	}
}

func TestParseSeasonDate(t *testing.T) {
	s := "2024-03-20"
	empty := ""
	cases := []struct {
		name   string
		input  *string
		isZero bool
		want   time.Time
	}{
		{"nil pointer yields zero", nil, true, time.Time{}},
		{"empty string yields zero", &empty, true, time.Time{}},
		{
			"bad format yields zero",
			func() *string { v := "not-a-date"; return &v }(),
			true,
			time.Time{},
		},
		{"valid YYYY-MM-DD parses", &s, false, time.Date(2024, 3, 20, 0, 0, 0, 0, time.UTC)},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := parseSeasonDate(c.input)
			if got.IsZero() != c.isZero {
				t.Fatalf("IsZero = %v, want %v (got %v)", got.IsZero(), c.isZero, got)
			}
			if !c.isZero && !got.Equal(c.want) {
				t.Errorf("got %v, want %v", got, c.want)
			}
		})
	}
}

func TestClient_Seasons(t *testing.T) {
	cases := []struct {
		name        string
		query       SeasonsQuery
		respStatus  int
		respBody    string
		wantPath    string
		wantQuery   url.Values
		wantErr     string
		wantIs      error
		wantLen     int
		wantSeason  string // first season's SeasonID (checked when non-empty)
		wantWildC   *bool
		wantRSStart string // first season's RegularSeasonStartDate as YYYY-MM-DD
	}{
		{
			name: "happy path: filtered seasons for 2024",
			query: SeasonsQuery{
				SportID: 1,
				Season:  2024,
				Fields:  "seasons,seasonId,regularSeasonStartDate",
			},
			respStatus: 200,
			respBody:   seasonsHappyBody,
			wantPath:   "/api/v1/seasons",
			wantQuery: url.Values{
				"sportId": {"1"},
				"season":  {"2024"},
				"fields":  {"seasons,seasonId,regularSeasonStartDate"},
			},
			wantLen:     1,
			wantSeason:  "2024",
			wantWildC:   boolPtr(true),
			wantRSStart: "2024-03-20",
		},
		{
			name: "filtered with DivisionID instead of SportID",
			query: SeasonsQuery{
				DivisionID: 200,
				Fields:     "seasons,seasonId",
			},
			respStatus: 200,
			respBody:   `{"seasons":[{"seasonId":"2024"}]}`,
			wantPath:   "/api/v1/seasons",
			wantQuery: url.Values{
				"divisionId": {"200"},
				"fields":     {"seasons,seasonId"},
			},
			wantLen:    1,
			wantSeason: "2024",
		},
		{
			name: "filtered with LeagueID instead of SportID",
			query: SeasonsQuery{
				LeagueID: 104,
			},
			respStatus: 200,
			respBody:   `{"seasons":[{"seasonId":"2024"}]}`,
			wantPath:   "/api/v1/seasons",
			wantQuery: url.Values{
				"leagueId": {"104"},
			},
			wantLen:    1,
			wantSeason: "2024",
		},
		{
			name: "happy path: All=true hits /seasons/all",
			query: SeasonsQuery{
				All:        true,
				DivisionID: 200,
				LeagueID:   103,
				Fields:     "seasons,seasonId",
			},
			respStatus: 200,
			respBody:   seasonsAllHappyBody,
			wantPath:   "/api/v1/seasons/all",
			wantQuery: url.Values{
				"divisionId": {"200"},
				"leagueId":   {"103"},
				"fields":     {"seasons,seasonId"},
			},
			wantLen:     2,
			wantSeason:  "1876",
			wantRSStart: "1876-04-22",
		},
		{
			name:    "missing required combo is rejected before HTTP call",
			query:   SeasonsQuery{Season: 2024},
			wantErr: "one of SportID, DivisionID, or LeagueID",
			wantIs:  ErrInvalidQuery,
		},
		{
			name:       "200 with no seasons yields empty slice",
			query:      SeasonsQuery{SportID: 1},
			respStatus: 200,
			respBody:   `{}`,
			wantLen:    0,
		},
		{
			name:       "200 with explicit empty seasons array",
			query:      SeasonsQuery{SportID: 1},
			respStatus: 200,
			respBody:   `{"seasons": []}`,
			wantLen:    0,
		},
		{
			name:       "404 returns ErrNotFound (filtered path)",
			query:      SeasonsQuery{SportID: 1},
			respStatus: 404,
			respBody:   `{}`,
			wantIs:     ErrNotFound,
			wantErr:    "not found",
		},
		{
			name:       "404 returns ErrNotFound (all path)",
			query:      SeasonsQuery{SportID: 1, All: true},
			respStatus: 404,
			respBody:   `{}`,
			wantIs:     ErrNotFound,
			wantErr:    "not found",
		},
		{
			name:       "5xx is wrapped (filtered)",
			query:      SeasonsQuery{SportID: 1},
			respStatus: 500,
			respBody:   `oops`,
			wantErr:    "unexpected status 500",
		},
		{
			name:       "5xx is wrapped (all)",
			query:      SeasonsQuery{SportID: 1, All: true},
			respStatus: 500,
			respBody:   `oops`,
			wantErr:    "unexpected status 500",
		},
		{
			name:       "malformed JSON is wrapped (filtered)",
			query:      SeasonsQuery{SportID: 1},
			respStatus: 200,
			respBody:   `not json`,
			wantErr:    "seasons",
		},
		{
			name:       "malformed JSON is wrapped (all)",
			query:      SeasonsQuery{SportID: 1, All: true},
			respStatus: 200,
			respBody:   `not json`,
			wantErr:    "seasons",
		},
		{
			name:       "network failure is wrapped (filtered)",
			query:      SeasonsQuery{SportID: 1},
			respStatus: 0,
			wantErr:    "seasons",
		},
		{
			name:       "network failure is wrapped (all)",
			query:      SeasonsQuery{SportID: 1, All: true},
			respStatus: 0,
			wantErr:    "seasons",
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
			s, err := client.Seasons(context.Background(), c.query)

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
			if s == nil {
				t.Fatal("expected non-nil Seasons")
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
			if got := len(s.Seasons); got != c.wantLen {
				t.Errorf("len(Seasons) = %d, want %d", got, c.wantLen)
			}
			if c.wantSeason != "" {
				first := s.Seasons[0]
				if first.SeasonID != c.wantSeason {
					t.Errorf("first SeasonID = %q, want %q", first.SeasonID, c.wantSeason)
				}
				if c.wantWildC != nil && first.HasWildcard != *c.wantWildC {
					t.Errorf("HasWildcard = %v, want %v", first.HasWildcard, *c.wantWildC)
				}
				if c.wantRSStart != "" {
					if got := first.RegularSeasonStartDate.Format(seasonDateFmt); got !=
						c.wantRSStart {
						t.Errorf(
							"RegularSeasonStartDate = %q, want %q",
							got,
							c.wantRSStart,
						)
					}
				}
				// On the fully-hydrated 2024 row, sanity-check the wider field
				// set so the field-promotion audit is exercised end-to-end.
				// Gated to the explicit "happy path: filtered" row that has
				// HasWildcard set so empty-body rows don't trip the asserts.
				if first.SeasonID == "2024" && c.wantWildC != nil && *c.wantWildC {
					if first.AllStarDate.Format(seasonDateFmt) != "2024-07-16" {
						t.Errorf("AllStarDate = %v", first.AllStarDate)
					}
					if first.SpringStartDate.Format(seasonDateFmt) != "2024-02-22" {
						t.Errorf("SpringStartDate = %v", first.SpringStartDate)
					}
					if first.SpringEndDate.Format(seasonDateFmt) != "2024-03-26" {
						t.Errorf("SpringEndDate = %v", first.SpringEndDate)
					}
					if first.PreSeasonStartDate.Format(seasonDateFmt) != "2024-01-01" {
						t.Errorf("PreSeasonStartDate = %v", first.PreSeasonStartDate)
					}
					if first.PreSeasonEndDate.Format(seasonDateFmt) != "2024-02-21" {
						t.Errorf("PreSeasonEndDate = %v", first.PreSeasonEndDate)
					}
					if first.SeasonStartDate.Format(seasonDateFmt) != "2024-02-22" {
						t.Errorf("SeasonStartDate = %v", first.SeasonStartDate)
					}
					if first.LastDate1stHalf.Format(seasonDateFmt) != "2024-07-14" {
						t.Errorf("LastDate1stHalf = %v", first.LastDate1stHalf)
					}
					if first.FirstDate2ndHalf.Format(seasonDateFmt) != "2024-07-19" {
						t.Errorf("FirstDate2ndHalf = %v", first.FirstDate2ndHalf)
					}
					if first.RegularSeasonEndDate.Format(seasonDateFmt) != "2024-09-30" {
						t.Errorf(
							"RegularSeasonEndDate = %v",
							first.RegularSeasonEndDate,
						)
					}
					if first.PostSeasonStartDate.Format(seasonDateFmt) != "2024-10-01" {
						t.Errorf(
							"PostSeasonStartDate = %v",
							first.PostSeasonStartDate,
						)
					}
					if first.PostSeasonEndDate.Format(seasonDateFmt) != "2024-10-30" {
						t.Errorf("PostSeasonEndDate = %v", first.PostSeasonEndDate)
					}
					if first.SeasonEndDate.Format(seasonDateFmt) != "2024-10-30" {
						t.Errorf("SeasonEndDate = %v", first.SeasonEndDate)
					}
					if first.OffseasonStartDate.Format(seasonDateFmt) != "2024-10-31" {
						t.Errorf(
							"OffseasonStartDate = %v",
							first.OffseasonStartDate,
						)
					}
					if first.OffSeasonEndDate.Format(seasonDateFmt) != "2024-12-31" {
						t.Errorf("OffSeasonEndDate = %v", first.OffSeasonEndDate)
					}
					if first.SeasonLevelGamedayType != "P" {
						t.Errorf("SeasonLevelGamedayType = %q",
							first.SeasonLevelGamedayType)
					}
					if first.GameLevelGamedayType != "P" {
						t.Errorf("GameLevelGamedayType = %q",
							first.GameLevelGamedayType)
					}
					if first.QualifierPlateAppearances != 3.1 {
						t.Errorf("QualifierPlateAppearances = %v",
							first.QualifierPlateAppearances)
					}
					if first.QualifierOutsPitched != 3.0 {
						t.Errorf("QualifierOutsPitched = %v",
							first.QualifierOutsPitched)
					}
				}
			}
		})
	}
}

func boolPtr(b bool) *bool { return &b }
