package rag

import "testing"

func TestFilter_Matches(t *testing.T) {
	tests := []struct {
		name     string
		filter   *Filter
		doc      Document
		expected bool
	}{
		{
			name:   "nil filter matches all",
			filter: nil,
			doc: Document{
				Metadata: map[string]interface{}{"key": "value"},
			},
			expected: true,
		},
		{
			name: "exact match",
			filter: &Filter{
				Metadata: map[string]interface{}{"language": "go"},
			},
			doc: Document{
				Metadata: map[string]interface{}{"language": "go"},
			},
			expected: true,
		},
		{
			name: "no match",
			filter: &Filter{
				Metadata: map[string]interface{}{"language": "go"},
			},
			doc: Document{
				Metadata: map[string]interface{}{"language": "python"},
			},
			expected: false,
		},
		{
			name: "multiple fields match",
			filter: &Filter{
				Metadata: map[string]interface{}{
					"language": "go",
					"level":    "beginner",
				},
			},
			doc: Document{
				Metadata: map[string]interface{}{
					"language": "go",
					"level":    "beginner",
					"extra":    "field",
				},
			},
			expected: true,
		},
		{
			name: "one field doesn't match",
			filter: &Filter{
				Metadata: map[string]interface{}{
					"language": "go",
					"level":    "advanced",
				},
			},
			doc: Document{
				Metadata: map[string]interface{}{
					"language": "go",
					"level":    "beginner",
				},
			},
			expected: false,
		},
		{
			name: "missing field in document",
			filter: &Filter{
				Metadata: map[string]interface{}{"nonexistent": "value"},
			},
			doc: Document{
				Metadata: map[string]interface{}{"language": "go"},
			},
			expected: false,
		},
		{
			name: "empty filter matches all",
			filter: &Filter{
				Metadata: map[string]interface{}{},
			},
			doc: Document{
				Metadata: map[string]interface{}{"language": "go"},
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.filter.Matches(tt.doc)
			if result != tt.expected {
				t.Errorf("Matches() = %v, expected %v", result, tt.expected)
			}
		})
	}
}
