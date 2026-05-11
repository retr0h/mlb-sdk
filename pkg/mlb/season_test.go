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
)

const seasonHappyBody = `{
  "copyright": "Copyright 2026 MLB Advanced Media, L.P.",
  "seasons": [{
    "seasonId": "2024",
    "hasWildcard": true,
    "regularSeasonStartDate": "2024-03-20",
    "regularSeasonEndDate":   "2024-09-30",
    "allStarDate":            "2024-07-16"
  }]
}`

func TestClient_Season(t *testing.T) {
	cases := []struct {
		name        string
		seasonID    string
		query       SeasonQuery
		respStatus  int
		respBody    string
		wantPath    string
		wantQuery   url.Values
		wantErr     string
		wantIs      error
		wantSeason  string
		wantWildC   bool
		wantRSStart string
	}{
		{
			name:     "happy path: 2024 with SportID=1",
			seasonID: "2024",
			query: SeasonQuery{
				SportID: 1,
				Fields:  "seasons,seasonId,regularSeasonStartDate",
			},
			respStatus: 200,
			respBody:   seasonHappyBody,
			wantPath:   "/api/v1/seasons/2024",
			wantQuery: url.Values{
				"sportId": {"1"},
				"fields":  {"seasons,seasonId,regularSeasonStartDate"},
			},
			wantSeason:  "2024",
			wantWildC:   true,
			wantRSStart: "2024-03-20",
		},
		{
			name:     "missing SportID rejected before HTTP",
			seasonID: "2024",
			query:    SeasonQuery{},
			wantErr:  "SportID is required",
			wantIs:   ErrInvalidQuery,
		},
		{
			name:       "200 with empty seasons array maps to ErrNotFound",
			seasonID:   "9999",
			query:      SeasonQuery{SportID: 1},
			respStatus: 200,
			respBody:   `{"seasons":[]}`,
			wantIs:     ErrNotFound,
			wantErr:    "not found",
		},
		{
			name:       "200 with missing seasons key maps to ErrNotFound",
			seasonID:   "9999",
			query:      SeasonQuery{SportID: 1},
			respStatus: 200,
			respBody:   `{}`,
			wantIs:     ErrNotFound,
			wantErr:    "not found",
		},
		{
			name:       "404 returns ErrNotFound",
			seasonID:   "9999",
			query:      SeasonQuery{SportID: 1},
			respStatus: 404,
			respBody:   `{}`,
			wantIs:     ErrNotFound,
			wantErr:    "not found",
		},
		{
			name:       "5xx is wrapped",
			seasonID:   "2024",
			query:      SeasonQuery{SportID: 1},
			respStatus: 500,
			respBody:   `oops`,
			wantErr:    "unexpected status 500",
		},
		{
			name:       "malformed JSON is wrapped",
			seasonID:   "2024",
			query:      SeasonQuery{SportID: 1},
			respStatus: 200,
			respBody:   `not json`,
			wantErr:    "season",
		},
		{
			name:       "network failure is wrapped",
			seasonID:   "2024",
			query:      SeasonQuery{SportID: 1},
			respStatus: 0,
			wantErr:    "season",
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
			s, err := client.Season(context.Background(), c.seasonID, c.query)

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
				t.Fatal("expected non-nil Season")
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
			if s.SeasonID != c.wantSeason {
				t.Errorf("SeasonID = %q, want %q", s.SeasonID, c.wantSeason)
			}
			if s.HasWildcard != c.wantWildC {
				t.Errorf("HasWildcard = %v, want %v", s.HasWildcard, c.wantWildC)
			}
			if c.wantRSStart != "" {
				if got := s.RegularSeasonStartDate.Format(seasonDateFmt); got != c.wantRSStart {
					t.Errorf("RegularSeasonStartDate = %q, want %q", got, c.wantRSStart)
				}
			}
		})
	}
}
