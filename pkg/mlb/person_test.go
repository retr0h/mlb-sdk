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
)
//
const personHappyBody = `{
  "copyright": "Copyright 2026 MLB Advanced Media, L.P.",
  "people": [{
    "id": 660271, "fullName": "Shohei Ohtani",
    "firstName": "Shohei", "lastName": "Ohtani",
    "link": "/api/v1/people/660271",
    "primaryNumber": "17",
    "birthDate": "1994-07-05", "currentAge": 31,
    "birthCity": "Oshu", "birthCountry": "Japan",
    "height": "6' 4\"", "weight": 210,
    "active": true,
    "primaryPosition": {"code": "Y", "name": "Two-Way Player",
                        "type": "Two-Way Player", "abbreviation": "TWP"},
    "useName": "Shohei", "useLastName": "Ohtani",
    "boxscoreName": "Ohtani", "nickName": "Showtime",
    "gender": "M", "isPlayer": true, "isVerified": false,
    "pronunciation": "show-HEY oh-TAWN-ee",
    "mlbDebutDate": "2018-03-29",
    "batSide":   {"code": "L", "description": "Left"},
    "pitchHand": {"code": "R", "description": "Right"},
    "nameFirstLast": "Shohei Ohtani",
    "nameSlug": "shohei-ohtani-660271",
    "firstLastName": "Shohei Ohtani",
    "lastFirstName": "Ohtani, Shohei",
    "lastInitName": "Ohtani, S",
    "initLastName": "S Ohtani",
    "fullFMLName": "Shohei Ohtani",
    "fullLFMName": "Ohtani, Shohei",
    "strikeZoneTop": 3.369,
    "strikeZoneBottom": 1.7
  }]
}`
//
func TestClient_Person(t *testing.T) {
	cases := []struct {
		name         string
		personID     int
		query        PersonQuery
		respStatus   int
		respBody     string
		wantPath     string
		wantQuery    url.Values
		wantErr      string
		wantIs       error
		wantHydrated bool
	}{
		{
			name:       "happy path: Ohtani",
			personID:   660271,
			query:      PersonQuery{Hydrate: "stats", Fields: "people,id,fullName"},
			respStatus: 200,
			respBody:   personHappyBody,
			wantPath:   "/api/v1/people/660271",
			wantQuery: url.Values{
				"hydrate": {"stats"},
				"fields":  {"people,id,fullName"},
			},
			wantHydrated: true,
		},
		{
			name:       "200 with minimal person",
			personID:   1,
			respStatus: 200,
			respBody:   `{"people":[{"id":1,"fullName":"Test"}]}`,
		},
		{
			name:       "200 with empty people maps to ErrNotFound",
			personID:   9999,
			respStatus: 200,
			respBody:   `{"people":[]}`,
			wantIs:     ErrNotFound,
			wantErr:    "not found",
		},
		{
			name:       "200 with missing people key maps to ErrNotFound",
			personID:   9999,
			respStatus: 200,
			respBody:   `{}`,
			wantIs:     ErrNotFound,
			wantErr:    "not found",
		},
		{
			name:       "404 returns ErrNotFound",
			personID:   9999,
			respStatus: 404,
			respBody:   `{}`,
			wantIs:     ErrNotFound,
			wantErr:    "not found",
		},
		{
			name:       "5xx is wrapped",
			personID:   660271,
			respStatus: 500,
			respBody:   `oops`,
			wantErr:    "unexpected status 500",
		},
		{
			name:       "malformed JSON is wrapped",
			personID:   660271,
			respStatus: 200,
			respBody:   `not json`,
			wantErr:    "person",
		},
		{
			name:       "network failure is wrapped",
			personID:   660271,
			respStatus: 0,
			wantErr:    "person",
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
			p, err := client.Person(context.Background(), c.personID, c.query)
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
			if p == nil {
				t.Fatal("expected non-nil PersonDetail")
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
			if p.ID != 660271 || p.FullName != "Shohei Ohtani" ||
				p.FirstName != "Shohei" || p.LastName != "Ohtani" ||
				p.Link != "/api/v1/people/660271" ||
				p.PrimaryNumber != "17" || p.BirthDate != "1994-07-05" ||
				p.CurrentAge != 31 || p.BirthCity != "Oshu" ||
				p.BirthCountry != "Japan" || p.Height != "6' 4\"" ||
				p.Weight != 210 || !p.Active {
				t.Errorf("basic fields = %+v", p)
			}
			if p.PrimaryPosition.Code != "Y" ||
				p.PrimaryPosition.Abbreviation != "TWP" {
				t.Errorf("PrimaryPosition = %+v", p.PrimaryPosition)
			}
			if p.UseName != "Shohei" || p.UseLastName != "Ohtani" ||
				p.BoxscoreName != "Ohtani" || p.NickName != "Showtime" ||
				p.Gender != "M" || !p.IsPlayer || p.IsVerified ||
				p.Pronunciation != "show-HEY oh-TAWN-ee" ||
				p.MlbDebutDate != "2018-03-29" {
				t.Errorf("identity fields = %+v", p)
			}
			if p.BatSide.Code != "L" || p.BatSide.Description != "Left" ||
				p.PitchHand.Code != "R" || p.PitchHand.Description != "Right" {
				t.Errorf("hands = bat=%+v pitch=%+v", p.BatSide, p.PitchHand)
			}
			if p.NameFirstLast != "Shohei Ohtani" ||
				p.NameSlug != "shohei-ohtani-660271" ||
				p.FirstLastName != "Shohei Ohtani" ||
				p.LastFirstName != "Ohtani, Shohei" ||
				p.LastInitName != "Ohtani, S" ||
				p.InitLastName != "S Ohtani" ||
				p.FullFMLName != "Shohei Ohtani" ||
				p.FullLFMName != "Ohtani, Shohei" {
				t.Errorf("name variants = %+v", p)
			}
			if p.StrikeZoneTop != 3.369 || p.StrikeZoneBottom != 1.7 {
				t.Errorf("zone = top=%v bottom=%v", p.StrikeZoneTop, p.StrikeZoneBottom)
			}
		})
	}
}
//
func TestClient_People(t *testing.T) {
	cases := []struct {
		name       string
		query      PeopleQuery
		respStatus int
		respBody   string
		wantPath   string
		wantQuery  url.Values
		wantErr    string
		wantIs     error
		wantLen    int
	}{
		{
			name: "happy path: two people",
			query: PeopleQuery{
				PersonIDs: "660271,545361",
				Hydrate:   "stats",
				Fields:    "people,id",
			},
			respStatus: 200,
			respBody:   `{"people":[{"id":660271,"fullName":"Ohtani"},{"id":545361,"fullName":"Trout"}]}`,
			wantPath:   "/api/v1/people",
			wantQuery: url.Values{
				"personIds": {"660271,545361"},
				"hydrate":   {"stats"},
				"fields":    {"people,id"},
			},
			wantLen: 2,
		},
		{
			name:    "missing PersonIDs rejected",
			query:   PeopleQuery{},
			wantErr: "PersonIDs is required",
			wantIs:  ErrInvalidQuery,
		},
		{
			name:       "200 with empty people",
			query:      PeopleQuery{PersonIDs: "9999"},
			respStatus: 200,
			respBody:   `{}`,
			wantLen:    0,
		},
		{
			name:       "404 returns ErrNotFound",
			query:      PeopleQuery{PersonIDs: "9999"},
			respStatus: 404,
			respBody:   `{}`,
			wantIs:     ErrNotFound,
			wantErr:    "not found",
		},
		{
			name:       "5xx is wrapped",
			query:      PeopleQuery{PersonIDs: "660271"},
			respStatus: 500,
			respBody:   `oops`,
			wantErr:    "unexpected status 500",
		},
		{
			name:       "malformed JSON is wrapped",
			query:      PeopleQuery{PersonIDs: "660271"},
			respStatus: 200,
			respBody:   `not json`,
			wantErr:    "people",
		},
		{
			name:       "network failure is wrapped",
			query:      PeopleQuery{PersonIDs: "660271"},
			respStatus: 0,
			wantErr:    "people",
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
			pp, err := client.People(context.Background(), c.query)
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
			if got := len(pp); got != c.wantLen {
				t.Errorf("len(People) = %d, want %d", got, c.wantLen)
			}
		})
	}
}
