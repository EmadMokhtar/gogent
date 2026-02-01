package builtin

import (
	"context"
	"fmt"

	"github.com/EmadMokhtar/gogent/tools"
)

// WebSearch is a tool for searching the web
type WebSearch struct{}

// NewWebSearch creates a new web search tool
func NewWebSearch() *WebSearch {
	return &WebSearch{}
}

// Name returns the tool name
func (w *WebSearch) Name() string {
	return "web_search"
}

// Description returns the tool description
func (w *WebSearch) Description() string {
	return "Searches the web for information based on a query"
}

// Parameters returns the tool parameters
func (w *WebSearch) Parameters() []tools.Parameter {
	return []tools.Parameter{
		{
			Name:        "query",
			Type:        tools.TypeString,
			Description: "The search query",
			Required:    true,
		},
		{
			Name:        "max_results",
			Type:        tools.TypeNumber,
			Description: "Maximum number of results to return",
			Required:    false,
			Default:     5,
		},
	}
}

// Execute runs the web search tool
func (w *WebSearch) Execute(ctx context.Context, input map[string]interface{}) (interface{}, error) {
	query, ok := input["query"].(string)
	if !ok {
		return nil, fmt.Errorf("query must be a string")
	}
	
	maxResults := 5
	if mr, ok := input["max_results"]; ok {
		switch v := mr.(type) {
		case int:
			maxResults = v
		case float64:
			maxResults = int(v)
		}
	}
	
	// This is a mock implementation
	// In a real implementation, this would call a search API
	results := []map[string]interface{}{
		{
			"title":   fmt.Sprintf("Mock result 1 for: %s", query),
			"url":     "https://example.com/1",
			"snippet": fmt.Sprintf("This is a mock search result snippet for %s", query),
		},
	}
	
	if maxResults > 1 {
		for i := 2; i <= maxResults && i <= 5; i++ {
			results = append(results, map[string]interface{}{
				"title":   fmt.Sprintf("Mock result %d for: %s", i, query),
				"url":     fmt.Sprintf("https://example.com/%d", i),
				"snippet": fmt.Sprintf("This is mock search result snippet %d for %s", i, query),
			})
		}
	}
	
	return map[string]interface{}{
		"query":   query,
		"results": results,
		"total":   len(results),
	}, nil
}
