package rag

// Filter defines criteria for filtering search results
type Filter struct {
	// Metadata contains key-value pairs that documents must match
	Metadata map[string]interface{}
}

// Matches checks if a document matches the filter criteria
func (f *Filter) Matches(doc Document) bool {
	if f == nil || f.Metadata == nil {
		return true
	}

	for key, expectedValue := range f.Metadata {
		actualValue, exists := doc.Metadata[key]
		if !exists {
			return false
		}
		if actualValue != expectedValue {
			return false
		}
	}

	return true
}
