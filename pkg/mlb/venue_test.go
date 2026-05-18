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
const venueHappyBody = `{
  "copyright": "Copyright 2026 MLB Advanced Media, L.P.",
  "venues": [{
    "id": 22,
    "name": "UNIQLO Field at Dodger Stadium",
    "link": "/api/v1/venues/22",
    "active": true,
    "season": "2026",
    "location": {
      "address1": "1000 Vin Scully Avenue",
      "address2": "Suite 100",
      "city": "Los Angeles",
      "state": "California",
      "stateAbbrev": "CA",
      "postalCode": "90012-1199",
      "defaultCoordinates": {"latitude": 34.07368, "longitude": -118.24053},
      "azimuthAngle": 26.0,
      "elevation": 515,
      "country": "USA",
      "phone": "(323) 224-1500"
    },
    "timeZone": {"tz": "PDT", "id": "America/Los_Angeles", "offset": -7, "offsetAtGameTime": -7},
    "fieldInfo": {
      "capacity": 56000, "turfType": "Grass", "roofType": "Open",
      "leftLine": 330, "leftCenter": 385, "center": 395, "rightCenter": 385, "rightLine": 330
    }
  }]
}`
//
func TestClient_Venue(t *testing.T) {
	cases := []struct {
		name       string
		venueID    int
		query      VenueQuery
		respStatus int
		respBody   string
		wantPath   string
		wantQuery  url.Values
		wantErr    string
		wantIs     error
		wantName   string
		wantCity   string
		wantCap    int
		wantTZ     string
	}{
		{
			name:    "happy path: hydrated venue",
			venueID: 22,
			query: VenueQuery{
				Season:  2026,
				Hydrate: "location,fieldInfo,timezone",
				Fields:  "venues,id,name",
			},
			respStatus: 200,
			respBody:   venueHappyBody,
			wantPath:   "/api/v1/venues/22",
			wantQuery: url.Values{
				"season":  {"2026"},
				"hydrate": {"location,fieldInfo,timezone"},
				"fields":  {"venues,id,name"},
			},
			wantName: "UNIQLO Field at Dodger Stadium",
			wantCity: "Los Angeles",
			wantCap:  56000,
			wantTZ:   "PDT",
		},
		{
			name:       "200 with minimal venue parses cleanly",
			venueID:    1,
			respStatus: 200,
			respBody:   `{"venues":[{"id":1,"name":"Angel Stadium"}]}`,
			wantName:   "Angel Stadium",
		},
		{
			name:       "200 with empty venues slice maps to ErrNotFound",
			venueID:    9999,
			respStatus: 200,
			respBody:   `{"venues":[]}`,
			wantIs:     ErrNotFound,
			wantErr:    "not found",
		},
		{
			name:       "200 with missing venues key maps to ErrNotFound",
			venueID:    9999,
			respStatus: 200,
			respBody:   `{}`,
			wantIs:     ErrNotFound,
			wantErr:    "not found",
		},
		{
			name:       "404 returns ErrNotFound",
			venueID:    9999,
			respStatus: 404,
			respBody:   `{}`,
			wantIs:     ErrNotFound,
			wantErr:    "not found",
		},
		{
			name:       "5xx is wrapped",
			venueID:    22,
			respStatus: 500,
			respBody:   `oops`,
			wantErr:    "unexpected status 500",
		},
		{
			name:       "malformed JSON is wrapped",
			venueID:    22,
			respStatus: 200,
			respBody:   `not json`,
			wantErr:    "venue",
		},
		{
			name:       "network failure is wrapped",
			venueID:    22,
			respStatus: 0,
			wantErr:    "venue",
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
			v, err := client.Venue(context.Background(), c.venueID, c.query)
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
			if v == nil {
				t.Fatal("expected non-nil Venue")
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
			if v.Name != c.wantName {
				t.Errorf("Name = %q, want %q", v.Name, c.wantName)
			}
			if c.wantCity != "" {
				if v.Location.City != c.wantCity {
					t.Errorf("Location.City = %q, want %q", v.Location.City, c.wantCity)
				}
				if v.Location.Address1 != "1000 Vin Scully Avenue" {
					t.Errorf("Location.Address1 = %q", v.Location.Address1)
				}
				if v.Location.Address2 != "Suite 100" {
					t.Errorf("Location.Address2 = %q", v.Location.Address2)
				}
				if v.Location.State != "California" || v.Location.StateAbbrev != "CA" {
					t.Errorf("State/Abbrev = %q/%q", v.Location.State, v.Location.StateAbbrev)
				}
				if v.Location.PostalCode != "90012-1199" {
					t.Errorf("PostalCode = %q", v.Location.PostalCode)
				}
				if v.Location.Country != "USA" || v.Location.Phone != "(323) 224-1500" {
					t.Errorf("Country/Phone = %q/%q", v.Location.Country, v.Location.Phone)
				}
				if v.Location.AzimuthAngle != 26.0 || v.Location.Elevation != 515 {
					t.Errorf("Azimuth/Elevation = %v/%v",
						v.Location.AzimuthAngle, v.Location.Elevation)
				}
				if v.Location.DefaultCoordinates.Latitude != 34.07368 ||
					v.Location.DefaultCoordinates.Longitude != -118.24053 {
					t.Errorf("DefaultCoordinates = %+v", v.Location.DefaultCoordinates)
				}
			}
			if c.wantCap > 0 {
				if v.FieldInfo.Capacity != c.wantCap {
					t.Errorf("FieldInfo.Capacity = %d, want %d", v.FieldInfo.Capacity, c.wantCap)
				}
				if v.FieldInfo.TurfType != "Grass" || v.FieldInfo.RoofType != "Open" {
					t.Errorf("Turf/Roof = %q/%q", v.FieldInfo.TurfType, v.FieldInfo.RoofType)
				}
				if v.FieldInfo.LeftLine != 330 || v.FieldInfo.LeftCenter != 385 ||
					v.FieldInfo.Center != 395 || v.FieldInfo.RightCenter != 385 ||
					v.FieldInfo.RightLine != 330 {
					t.Errorf("dimensions = %+v", v.FieldInfo)
				}
			}
			if c.wantTZ != "" {
				if v.TimeZone.TZ != c.wantTZ {
					t.Errorf("TimeZone.TZ = %q, want %q", v.TimeZone.TZ, c.wantTZ)
				}
				if v.TimeZone.ID != "America/Los_Angeles" ||
					v.TimeZone.Offset != -7 || v.TimeZone.OffsetAtGameTime != -7 {
					t.Errorf("TimeZone = %+v", v.TimeZone)
				}
			}
			if c.wantName == "UNIQLO Field at Dodger Stadium" {
				if v.ID != 22 || v.Link != "/api/v1/venues/22" ||
					!v.Active || v.Season != "2026" {
					t.Errorf(
						"core fields = id=%d link=%q active=%v season=%q",
						v.ID, v.Link, v.Active, v.Season,
					)
				}
			}
		})
	}
}
