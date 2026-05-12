// Copyright (c) 2026 John Dewey
//
// SPDX-License-Identifier: MIT

package api

import (
	"strings"
	"testing"
)

func TestOpenAPISpec(t *testing.T) {
	cases := []struct {
		name     string
		wantLen  bool
		wantYAML string
	}{
		{"non-empty", true, "openapi: 3.0"},
		{"has paths", true, "paths:"},
		{"has components", true, "components:"},
	}
	spec := OpenAPISpec()
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.wantLen && len(spec) == 0 {
				t.Fatal("OpenAPISpec() returned empty")
			}
			if c.wantYAML != "" && !strings.Contains(string(spec), c.wantYAML) {
				t.Errorf("spec missing %q", c.wantYAML)
			}
		})
	}
}
