// Copyright (c) 2026 John Dewey
//
// SPDX-License-Identifier: MIT

// Package api embeds the hand-authored OpenAPI 3.0 specification for the
// MLB Stats API. Import this package to access the raw spec bytes for code
// generation, MCP tool registration, or documentation.
package api

import _ "embed"

//go:embed openapi.yaml
var spec []byte

// OpenAPISpec returns the raw OpenAPI 3.0 YAML specification.
func OpenAPISpec() []byte {
	return spec
}
