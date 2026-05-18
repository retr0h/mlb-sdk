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
	"time"
)
//
const transactionsHappyBody = `{
  "copyright": "Copyright 2026 MLB Advanced Media, L.P.",
  "transactions": [
    {
      "id": 791179,
      "person": {"id": 801903, "fullName": "Brock Tibbitts",
                 "link": "/api/v1/people/801903"},
      "toTeam": {"id": 424, "name": "Dunedin Blue Jays"},
      "date": "2024-07-30", "effectiveDate": "2024-07-30",
      "resolutionDate": "2024-07-30",
      "typeCode": "ASG", "typeDesc": "Assigned",
      "description": "C Brock Tibbitts assigned to Dunedin Blue Jays."
    },
    {
      "id": 791029,
      "person": {"id": 691002, "fullName": "Yohandy Morales",
                 "link": "/api/v1/people/691002"},
      "fromTeam": {"id": 547, "name": "Harrisburg Senators"},
      "toTeam":   {"id": 436, "name": "Fredericksburg Nationals"},
      "date": "2024-07-30", "effectiveDate": "2024-07-30",
      "typeCode": "ASG", "typeDesc": "Assigned"
    }
  ]
}`
//
func TestClient_Transactions(t *testing.T) {
	start := time.Date(2024, 7, 30, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 7, 31, 0, 0, 0, 0, time.UTC)
	on := time.Date(2024, 7, 30, 0, 0, 0, 0, time.UTC)
//
	cases := []struct {
		name         string
		query        TransactionsQuery
		respStatus   int
		respBody     string
		wantPath     string
		wantQuery    url.Values
		wantErr      string
		wantIs       error
		wantLen      int
		wantHydrated bool
	}{
		{
			name: "happy path: date range + sportId + fields",
			query: TransactionsQuery{
				StartDate: start,
				EndDate:   end,
				SportID:   1,
				Fields:    "transactions,id,person",
			},
			respStatus: 200,
			respBody:   transactionsHappyBody,
			wantPath:   "/api/v1/transactions",
			wantQuery: url.Values{
				"startDate": {"2024-07-30"},
				"endDate":   {"2024-07-31"},
				"sportId":   {"1"},
				"fields":    {"transactions,id,person"},
			},
			wantLen:      2,
			wantHydrated: true,
		},
		{
			name:       "filter with TeamID alone",
			query:      TransactionsQuery{TeamID: 119},
			respStatus: 200,
			respBody:   `{"transactions":[]}`,
			wantPath:   "/api/v1/transactions",
			wantQuery:  url.Values{"teamId": {"119"}},
			wantLen:    0,
		},
		{
			name:       "filter with PlayerID alone",
			query:      TransactionsQuery{PlayerID: 691002},
			respStatus: 200,
			respBody:   `{"transactions":[]}`,
			wantPath:   "/api/v1/transactions",
			wantQuery:  url.Values{"playerId": {"691002"}},
			wantLen:    0,
		},
		{
			name:       "filter with On alone",
			query:      TransactionsQuery{On: on},
			respStatus: 200,
			respBody:   `{"transactions":[]}`,
			wantPath:   "/api/v1/transactions",
			wantQuery:  url.Values{"date": {"2024-07-30"}},
			wantLen:    0,
		},
		{
			name:    "missing required combo rejected before HTTP",
			query:   TransactionsQuery{SportID: 1},
			wantErr: "one of TeamID, PlayerID, On, or StartDate+EndDate",
			wantIs:  ErrInvalidQuery,
		},
		{
			name:    "StartDate without EndDate rejected",
			query:   TransactionsQuery{StartDate: start},
			wantErr: "StartDate and EndDate must be set together",
			wantIs:  ErrInvalidQuery,
		},
		{
			name:    "EndDate without StartDate rejected",
			query:   TransactionsQuery{EndDate: end},
			wantErr: "StartDate and EndDate must be set together",
			wantIs:  ErrInvalidQuery,
		},
		{
			name:       "200 with no transactions yields empty slice",
			query:      TransactionsQuery{TeamID: 119},
			respStatus: 200,
			respBody:   `{}`,
			wantLen:    0,
		},
		{
			name:       "404 returns ErrNotFound",
			query:      TransactionsQuery{TeamID: 119},
			respStatus: 404,
			respBody:   `{}`,
			wantIs:     ErrNotFound,
			wantErr:    "not found",
		},
		{
			name:       "5xx is wrapped",
			query:      TransactionsQuery{TeamID: 119},
			respStatus: 500,
			respBody:   `oops`,
			wantErr:    "unexpected status 500",
		},
		{
			name:       "malformed JSON is wrapped",
			query:      TransactionsQuery{TeamID: 119},
			respStatus: 200,
			respBody:   `not json`,
			wantErr:    "transactions",
		},
		{
			name:       "network failure is wrapped",
			query:      TransactionsQuery{TeamID: 119},
			respStatus: 0,
			wantErr:    "transactions",
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
			tx, err := client.Transactions(context.Background(), c.query)
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
			if tx == nil {
				t.Fatal("expected non-nil Transactions")
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
			if got := len(tx.Transactions); got != c.wantLen {
				t.Errorf("len(Transactions) = %d, want %d", got, c.wantLen)
			}
			if !c.wantHydrated {
				return
			}
			first := tx.Transactions[0]
			if first.ID != 791179 || first.Date != "2024-07-30" ||
				first.EffectiveDate != "2024-07-30" ||
				first.ResolutionDate != "2024-07-30" ||
				first.TypeCode != "ASG" || first.TypeDesc != "Assigned" ||
				first.Description != "C Brock Tibbitts assigned to Dunedin Blue Jays." {
				t.Errorf("transactions[0] = %+v", first)
			}
			if first.Person.ID != 801903 || first.Person.FullName != "Brock Tibbitts" ||
				first.Person.Link != "/api/v1/people/801903" {
				t.Errorf("transactions[0].Person = %+v", first.Person)
			}
			if first.ToTeam.ID != TeamID(424) || first.ToTeam.Name != "Dunedin Blue Jays" {
				t.Errorf("transactions[0].ToTeam = %+v", first.ToTeam)
			}
			// FromTeam zero on the first row (initial assignment).
			if first.FromTeam.ID != 0 || first.FromTeam.Name != "" {
				t.Errorf("transactions[0].FromTeam = %+v (want zero)", first.FromTeam)
			}
			// Second row has FromTeam set.
			second := tx.Transactions[1]
			if second.FromTeam.ID != TeamID(547) ||
				second.FromTeam.Name != "Harrisburg Senators" {
				t.Errorf("transactions[1].FromTeam = %+v", second.FromTeam)
			}
		})
	}
}
