// Copyright 2025 Redpanda Data, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package generator

import (
	"testing"
)

func TestHasMCPTag(t *testing.T) {
	tests := []struct {
		name     string
		comment  string
		expected bool
	}{
		{
			name:     "simple @mcp tag",
			comment:  "@mcp",
			expected: true,
		},
		{
			name:     "@mcp with description",
			comment:  "@mcp Enable this endpoint for MCP",
			expected: true,
		},
		{
			name: "multiline comment with @mcp",
			comment: `This is a method description
@mcp
Additional details here`,
			expected: true,
		},
		{
			name: "multiline comment with @mcp and description",
			comment: `This is a method description
@mcp Enable this endpoint for MCP
Additional details here`,
			expected: true,
		},
		{
			name:     "no @mcp tag",
			comment:  "This is a regular method description",
			expected: false,
		},
		{
			name:     "empty comment",
			comment:  "",
			expected: false,
		},
		{
			name:     "@mcp in middle of line (should not match)",
			comment:  "This has @mcp in the middle",
			expected: false,
		},
		{
			name: "multiple lines without @mcp",
			comment: `This is a method description
with multiple lines
but no @mcp tag`,
			expected: false,
		},
		{
			name: "@mcp with indentation",
			comment: `  @mcp
  Enable this endpoint`,
			expected: true,
		},
		{
			name: "mixed tags with @mcp",
			comment: `buf:lint:ignore
@mcp Enable this endpoint
@ignore-comment some other tag`,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hasMCPTag(tt.comment)
			if result != tt.expected {
				t.Errorf("hasMCPTag(%q) = %v, expected %v", tt.comment, result, tt.expected)
			}
		})
	}
}

func TestCleanCommentRemovesMCPTag(t *testing.T) {
	tests := []struct {
		name     string
		comment  string
		expected string
	}{
		{
			name:     "simple @mcp tag removed",
			comment:  "@mcp",
			expected: "",
		},
		{
			name:     "@mcp with description removed",
			comment:  "@mcp Enable this endpoint for MCP",
			expected: "",
		},
		{
			name: "multiline comment with @mcp removed",
			comment: `This is a method description
@mcp
Additional details here`,
			expected: `This is a method description
Additional details here`,
		},
		{
			name: "preserves other content",
			comment: `CreateItem creates a new item
@mcp
Returns the created item with ID`,
			expected: `CreateItem creates a new item
Returns the created item with ID`,
		},
		{
			name: "removes multiple stripped prefixes",
			comment: `buf:lint:ignore
@mcp Enable this endpoint
@ignore-comment some other tag
Actual description here`,
			expected: `Actual description here`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cleanComment(tt.comment)
			if result != tt.expected {
				t.Errorf("cleanComment(%q) = %q, expected %q", tt.comment, result, tt.expected)
			}
		})
	}
}
