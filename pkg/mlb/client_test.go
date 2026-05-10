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
)

func TestNew(t *testing.T) {
	// markerTransport flips a flag the moment it sees a request, so the test
	// can confirm WithHTTPClient actually rewires the transport rather than
	// being silently ignored.
	type markerTransport struct {
		called bool
		inner  http.RoundTripper
	}
	roundTrip := func(m *markerTransport) func(*http.Request) (*http.Response, error) {
		return func(r *http.Request) (*http.Response, error) {
			m.called = true
			return m.inner.RoundTrip(r)
		}
	}

	cases := []struct {
		name      string
		opts      func(srvURL string) (opts []Option, marker *markerTransport)
		wantErr   string
		invokeErr bool // whether to drive Schedule afterwards (to exercise transport)
	}{
		{
			name:    "no options uses defaults",
			opts:    func(_ string) ([]Option, *markerTransport) { return nil, nil },
			wantErr: "",
		},
		{
			name: "WithBaseURL is honored",
			opts: func(u string) ([]Option, *markerTransport) {
				return []Option{WithBaseURL(u)}, nil
			},
			invokeErr: true,
		},
		{
			name: "WithHTTPClient rewires the transport",
			opts: func(u string) ([]Option, *markerTransport) {
				m := &markerTransport{inner: http.DefaultTransport}
				return []Option{
					WithBaseURL(u),
					WithHTTPClient(&http.Client{Transport: roundTripperFunc(roundTrip(m))}),
				}, m
			},
			invokeErr: true,
		},
		{
			name: "options stack — both can be applied",
			opts: func(u string) ([]Option, *markerTransport) {
				return []Option{WithBaseURL(u), WithHTTPClient(http.DefaultClient)}, nil
			},
			invokeErr: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			srv := httptest.NewServer(
				http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					_, _ = w.Write([]byte(`{}`))
				}),
			)
			defer srv.Close()

			opts, marker := c.opts(srv.URL)
			client, err := New(opts...)

			if c.wantErr != "" {
				if err == nil {
					t.Fatalf("expected error containing %q, got nil", c.wantErr)
				}
				if !strings.Contains(err.Error(), c.wantErr) {
					t.Errorf("err = %v, want substring %q", err, c.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if client == nil {
				t.Fatal("New returned nil client without error")
			}

			// Drive the client through one call so we can verify transport
			// rewiring actually took effect when the test asked for it.
			if c.invokeErr {
				if _, err := client.Schedule(context.Background(), ScheduleQuery{}); err != nil {
					t.Fatalf("Schedule via httptest server: %v", err)
				}
			}
			if marker != nil && !marker.called {
				t.Errorf("custom transport was not invoked — WithHTTPClient is not wiring through")
			}
		})
	}
}

// roundTripperFunc is the http.RoundTripper analogue of http.HandlerFunc — lets
// us inject ad-hoc round-trip logic from a test without declaring a new type
// per case.
type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }
