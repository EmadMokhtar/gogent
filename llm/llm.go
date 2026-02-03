package llm

import "context"

// Message represents a message in a conversation
type Message struct {
	Role    string
	Content string
}

// Response represents an LLM response
type Response struct {
	Content string
	Model   string
	Usage   *Usage
}

// Usage represents token usage information
type Usage struct {
	PromptTokens     int
	CompletionTokens int
	TotalTokens      int
}

// LLM defines the interface for language model providers
type LLM interface {
	// Generate generates a response for the given messages
	Generate(ctx context.Context, messages []Message) (*Response, error)

	// Model returns the model name
	Model() string
}
