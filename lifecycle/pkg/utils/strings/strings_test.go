// Copyright © 2025 sealos.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package strings_test

import (
	"testing"

	"github.com/labring/sealos/pkg/utils/strings"
)

func TestNotInIPList(t *testing.T) {
	tests := []struct {
		name     string
		slice    []string
		key      string
		expected bool
	}{
		{
			name:     "empty slice",
			slice:    []string{},
			key:      "192.168.1.1",
			expected: true,
		},
		{
			name:     "key not in slice",
			slice:    []string{"192.168.1.2:80", "192.168.1.3:80"},
			key:      "192.168.1.1",
			expected: true,
		},
		{
			name:     "key in slice",
			slice:    []string{"192.168.1.1:80", "192.168.1.2:80"},
			key:      "192.168.1.1",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := strings.NotInIPList(tt.slice, tt.key)
			if result != tt.expected {
				t.Errorf("NotInIPList() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIsEmptyLine(t *testing.T) {
	tests := []struct {
		name     string
		str      string
		expected bool
	}{
		{
			name:     "empty string",
			str:      "",
			expected: true,
		},
		{
			name:     "whitespace only",
			str:      "   ",
			expected: true,
		},
		{
			name:     "tabs and newlines",
			str:      "\t\n",
			expected: true,
		},
		{
			name:     "non-empty string",
			str:      "hello",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := strings.IsEmptyLine(tt.str)
			if result != tt.expected {
				t.Errorf("IsEmptyLine() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestTrimWS(t *testing.T) {
	tests := []struct {
		name     string
		str      string
		expected string
	}{
		{
			name:     "empty string",
			str:      "",
			expected: "",
		},
		{
			name:     "string with whitespace",
			str:      "\nhello\t",
			expected: "hello",
		},
		{
			name:     "string with multiple whitespace",
			str:      "\n\thello\n\t",
			expected: "hello",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := strings.TrimWS(tt.str)
			if result != tt.expected {
				t.Errorf("TrimWS() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestTrimSpaceWS(t *testing.T) {
	tests := []struct {
		name     string
		str      string
		expected string
	}{
		{
			name:     "empty string",
			str:      "",
			expected: "",
		},
		{
			name:     "string with right whitespace",
			str:      "hello \t\n",
			expected: "hello",
		},
		{
			name:     "string with left whitespace",
			str:      " \t\nhello",
			expected: " \t\nhello",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := strings.TrimSpaceWS(tt.str)
			if result != tt.expected {
				t.Errorf("TrimSpaceWS() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestFilterNonEmptyFromSlice(t *testing.T) {
	tests := []struct {
		name     string
		list     []string
		expected []string
	}{
		{
			name:     "empty slice",
			list:     []string{},
			expected: nil,
		},
		{
			name:     "slice with empty strings",
			list:     []string{"", "  ", "\t"},
			expected: nil,
		},
		{
			name:     "slice with non-empty strings",
			list:     []string{"hello", "", "world", "  "},
			expected: []string{"hello", "world"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := strings.FilterNonEmptyFromSlice(tt.list)
			if len(result) != len(tt.expected) {
				t.Errorf("FilterNonEmptyFromSlice() length = %v, want %v", len(result), len(tt.expected))
				return
			}
			for i := range result {
				if result[i] != tt.expected[i] {
					t.Errorf("FilterNonEmptyFromSlice()[%d] = %v, want %v", i, result[i], tt.expected[i])
				}
			}
		})
	}
}

func TestFilterNonEmptyFromString(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		sep      string
		expected []string
	}{
		{
			name:     "empty string",
			s:        "",
			sep:      ",",
			expected: nil,
		},
		{
			name:     "string with empty parts",
			s:        ",,  ,\t,",
			sep:      ",",
			expected: nil,
		},
		{
			name:     "string with non-empty parts",
			s:        "hello,,world,  ,test",
			sep:      ",",
			expected: []string{"hello", "world", "test"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := strings.FilterNonEmptyFromString(tt.s, tt.sep)
			if len(result) != len(tt.expected) {
				t.Errorf("FilterNonEmptyFromString() length = %v, want %v", len(result), len(tt.expected))
				return
			}
			for i := range result {
				if result[i] != tt.expected[i] {
					t.Errorf("FilterNonEmptyFromString()[%d] = %v, want %v", i, result[i], tt.expected[i])
				}
			}
		})
	}
}

func TestRemoveDuplicate(t *testing.T) {
	tests := []struct {
		name     string
		list     []string
		expected []string
	}{
		{
			name:     "empty slice",
			list:     []string{},
			expected: nil,
		},
		{
			name:     "slice with duplicates",
			list:     []string{"a", "b", "a", "c", "b"},
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "slice without duplicates",
			list:     []string{"a", "b", "c"},
			expected: []string{"a", "b", "c"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := strings.RemoveDuplicate(tt.list)
			if len(result) != len(tt.expected) {
				t.Errorf("RemoveDuplicate() length = %v, want %v", len(result), len(tt.expected))
				return
			}
			for i := range result {
				if result[i] != tt.expected[i] {
					t.Errorf("RemoveDuplicate()[%d] = %v, want %v", i, result[i], tt.expected[i])
				}
			}
		})
	}
}

func TestRemoveSubSlice(t *testing.T) {
	tests := []struct {
		name     string
		src      []string
		dst      []string
		expected []string
	}{
		{
			name:     "empty source",
			src:      []string{},
			dst:      []string{"a", "b"},
			expected: nil,
		},
		{
			name:     "empty destination",
			src:      []string{"a", "b"},
			dst:      []string{},
			expected: []string{"a", "b"},
		},
		{
			name:     "remove elements",
			src:      []string{"a", "b", "c", "d"},
			dst:      []string{"b", "d"},
			expected: []string{"a", "c"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := strings.RemoveSubSlice(tt.src, tt.dst)
			if len(result) != len(tt.expected) {
				t.Errorf("RemoveSubSlice() length = %v, want %v", len(result), len(tt.expected))
				return
			}
			for i := range result {
				if result[i] != tt.expected[i] {
					t.Errorf("RemoveSubSlice()[%d] = %v, want %v", i, result[i], tt.expected[i])
				}
			}
		})
	}
}

func TestRemoveFromSlice(t *testing.T) {
	tests := []struct {
		name     string
		ss       []string
		s        string
		expected []string
	}{
		{
			name:     "empty slice",
			ss:       []string{},
			s:        "a",
			expected: nil,
		},
		{
			name:     "remove existing element",
			ss:       []string{"a", "b", "c"},
			s:        "b",
			expected: []string{"a", "c"},
		},
		{
			name:     "remove non-existing element",
			ss:       []string{"a", "b", "c"},
			s:        "d",
			expected: []string{"a", "b", "c"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := strings.RemoveFromSlice(tt.ss, tt.s)
			if len(result) != len(tt.expected) {
				t.Errorf("RemoveFromSlice() length = %v, want %v", len(result), len(tt.expected))
				return
			}
			for i := range result {
				if result[i] != tt.expected[i] {
					t.Errorf("RemoveFromSlice()[%d] = %v, want %v", i, result[i], tt.expected[i])
				}
			}
		})
	}
}

func TestMerge(t *testing.T) {
	tests := []struct {
		name     string
		ss       []string
		s        string
		expected []string
	}{
		{
			name:     "empty slice",
			ss:       []string{},
			s:        "a",
			expected: []string{"a"},
		},
		{
			name:     "merge new element",
			ss:       []string{"a", "b"},
			s:        "c",
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "merge existing element",
			ss:       []string{"a", "b", "c"},
			s:        "b",
			expected: []string{"a", "c", "b"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := strings.Merge(tt.ss, tt.s)
			if len(result) != len(tt.expected) {
				t.Errorf("Merge() length = %v, want %v", len(result), len(tt.expected))
				return
			}
			for i := range result {
				if result[i] != tt.expected[i] {
					t.Errorf("Merge()[%d] = %v, want %v", i, result[i], tt.expected[i])
				}
			}
		})
	}
}

func TestFormatSize(t *testing.T) {
	tests := []struct {
		name     string
		size     int64
		expected string
	}{
		{
			name:     "bytes",
			size:     500,
			expected: "500.00B",
		},
		{
			name:     "kilobytes",
			size:     1500,
			expected: "1.46KB",
		},
		{
			name:     "megabytes",
			size:     1500000,
			expected: "1.43MB",
		},
		{
			name:     "gigabytes",
			size:     1500000000,
			expected: "1.40GB",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := strings.FormatSize(tt.size)
			if result != tt.expected {
				t.Errorf("FormatSize() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIsLetterOrNumber(t *testing.T) {
	tests := []struct {
		name     string
		k        string
		expected bool
	}{
		{
			name:     "empty string",
			k:        "",
			expected: true,
		},
		{
			name:     "letters only",
			k:        "abcDEF",
			expected: true,
		},
		{
			name:     "numbers only",
			k:        "123456",
			expected: true,
		},
		{
			name:     "letters and numbers",
			k:        "abc123DEF",
			expected: true,
		},
		{
			name:     "with underscore",
			k:        "abc_123",
			expected: true,
		},
		{
			name:     "with special characters",
			k:        "abc@123",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := strings.IsLetterOrNumber(tt.k)
			if result != tt.expected {
				t.Errorf("IsLetterOrNumber() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestRenderShellWithEnv(t *testing.T) {
	tests := []struct {
		name     string
		shell    string
		envs     map[string]string
		expected string
	}{
		{
			name:     "no envs",
			shell:    "echo hello",
			envs:     map[string]string{},
			expected: "echo hello",
		},
		{
			name:  "single env",
			shell: "echo $MESSAGE",
			envs: map[string]string{
				"MESSAGE": "hello",
			},
			expected: `export MESSAGE="hello" ; echo $MESSAGE`,
		},
		{
			name:  "multiple envs",
			shell: "echo $GREETING $NAME",
			envs: map[string]string{
				"NAME":     "world",
				"GREETING": "hello",
			},
			expected: `export NAME="world" GREETING="hello" ; echo $GREETING $NAME`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := strings.RenderShellWithEnv(tt.shell, tt.envs)
			if result != tt.expected {
				t.Errorf("RenderShellWithEnv() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestRenderTextWithEnv(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		envs     map[string]string
		expected string
	}{
		{
			name: "no placeholders",
			text: "Hello world",
			envs: map[string]string{
				"NAME": "John",
			},
			expected: "Hello world",
		},
		{
			name: "with $(VAR) format",
			text: "Hello $(NAME)",
			envs: map[string]string{
				"NAME": "John",
			},
			expected: "Hello John",
		},
		{
			name: "with ${VAR} format",
			text: "Hello ${NAME}",
			envs: map[string]string{
				"NAME": "John",
			},
			expected: "Hello John",
		},
		{
			name: "with $VAR format",
			text: "Hello $NAME",
			envs: map[string]string{
				"NAME": "John",
			},
			expected: "Hello John",
		},
		{
			name: "multiple variables",
			text: "$(GREETING) ${NAME}! $TIME",
			envs: map[string]string{
				"GREETING": "Hello",
				"NAME":     "John",
				"TIME":     "morning",
			},
			expected: "Hello John! morning",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := strings.RenderTextWithEnv(tt.text, tt.envs)
			if result != tt.expected {
				t.Errorf("RenderTextWithEnv() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestTrimQuotes(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		expected string
	}{
		{
			name:     "no quotes",
			s:        "hello",
			expected: "hello",
		},
		{
			name:     "double quotes",
			s:        `"hello"`,
			expected: "hello",
		},
		{
			name:     "single quotes",
			s:        "'hello'",
			expected: "hello",
		},
		{
			name:     "mixed quotes",
			s:        `"hello'`,
			expected: `"hello'`,
		},
		{
			name:     "empty string",
			s:        "",
			expected: "",
		},
		{
			name:     "single character",
			s:        "a",
			expected: "a",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := strings.TrimQuotes(tt.s)
			if result != tt.expected {
				t.Errorf("TrimQuotes() = %v, want %v", result, tt.expected)
			}
		})
	}
}
