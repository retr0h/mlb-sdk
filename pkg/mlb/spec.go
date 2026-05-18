// Copyright (c) 2026 John Dewey
//
// SPDX-License-Identifier: MIT
//
package mlb
//
import "github.com/retr0h/mlb-sdk/pkg/api"
//
// OpenAPISpec returns the raw OpenAPI 3.0 YAML specification that this SDK
// is generated from. Consumers can use it for code generation, MCP tool
// registration, or documentation purposes. This is a convenience re-export
// of [api.OpenAPISpec].
func OpenAPISpec() []byte {
	return api.OpenAPISpec()
}
