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

// Package mlb is the public, idiomatic Go SDK for the MLB Stats API. It is the
// only package consumers of this module should import; the generated client
// under internal/gen is an implementation detail.
package mlb

import (
	"net/http"

	"github.com/retr0h/mlb-sdk/internal/gen"
)

// DefaultBaseURL is the public MLB Stats API host. Override via WithBaseURL
// for testing against an httptest server or a private mirror.
const DefaultBaseURL = "https://statsapi.mlb.com"

// Client is the public MLB Stats API client. Construct with New.
type Client struct {
	raw *gen.ClientWithResponses
}

// Option configures a Client.
type Option func(*config)

type config struct {
	baseURL    string
	httpClient *http.Client
}

// WithBaseURL overrides the API host. Useful for tests.
func WithBaseURL(u string) Option {
	return func(c *config) { c.baseURL = u }
}

// WithHTTPClient overrides the underlying http.Client. Use to inject timeouts,
// transport-level instrumentation, or recorded fixtures.
func WithHTTPClient(h *http.Client) Option {
	return func(c *config) { c.httpClient = h }
}

// New returns a Client configured for the public MLB Stats API. Options may
// override the base URL and http.Client. Construction never fails — the
// underlying generated client only validates the server URL when it actually
// dispatches a request, surfacing errors there instead.
func New(opts ...Option) *Client {
	cfg := config{baseURL: DefaultBaseURL, httpClient: http.DefaultClient}
	for _, opt := range opts {
		opt(&cfg)
	}

	// gen.NewClientWithResponses returns an error only when one of its
	// ClientOptions fails. We pass only WithHTTPClient, which never errors
	// (it just stores the client). If a future oapi-codegen release adds a
	// fallible option, switch to (*Client, error) here.
	raw, _ := gen.NewClientWithResponses( //nolint:errcheck
		cfg.baseURL,
		gen.WithHTTPClient(cfg.httpClient),
	)
	return &Client{raw: raw}
}
